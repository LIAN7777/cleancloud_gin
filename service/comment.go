package service

import (
	dto "GinProject/dto/user"
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

type BoolValue bool

func (b BoolValue) Value() (driver.Value, error) {
	if b {
		return true, nil
	}
	return false, nil
}

func PublishComment(comment *dto.CommentForm) bool {
	err := utils.Publish("amq.direct", "comment", comment)
	if err != nil {
		return false
	}
	return true
}

func AddComment(msg []byte) {
	comment := &model.Comment{}
	err := json.Unmarshal(msg, comment)
	if err != nil {
		log.Print("comment format error")
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	comment.Time = &t
	comment.Status = &[]uint8{1}
	ctx := context.Background()
	dComment := query.Comment
	//评论存入数据库
	err = dComment.WithContext(ctx).Create(comment)
	if err != nil {
		log.Print("add comment fail:caused by\n", err)
		return
	}
	//发送通知给用户有新评论
	err = utils.Publish("amq.direct", "comment_message", comment)
	if err != nil {
		log.Print("send comment message fail")
	}
}

func GetCommentById(id string) interface{} {
	//使用Redis分布式锁
	ctx := context.Background()
	lock := utils.NewRedisLock("lock:comment", utils.Client)
	//上锁
	_ = lock.Lock(ctx)
	//函数返回前释放锁
	defer func(lock *utils.RedisLock) {
		_ = lock.Unlock()
	}(lock)
	//先查Redis
	comment, err := utils.RedisGetModel("cache:comment:"+id, model.Comment{})
	if err == nil {
		return comment
	}
	//查Mysql
	C := query.Comment
	Id, _ := strconv.Atoi(id)
	comment, err = C.WithContext(ctx).Where(C.CommentID.Eq(int64(Id))).First()
	if err != nil {
		//缓存空值
		utils.Client.Set("cache:comment:"+id, "", time.Minute*30)
		return nil
	} else {
		utils.RedisSetModel("cache:comment:"+id, comment)
		return comment
	}
}

func GetCommentByBlog(blogId string) []interface{} {
	var idSet []string    //blog对应的comment idSet
	var res []interface{} //comment最终结果
	//先查redis是否有idSet
	ids, err := utils.Client.SMembers("cache:comment:blog:" + blogId).Result()
	if err == nil && cap(ids) != 0 {
		idSet = ids
	} else {
		//不存在idSet 查MySQL
		C := query.Comment
		ctx := context.Background()
		BlogId, _ := strconv.Atoi(blogId)
		comments, err := C.WithContext(ctx).Where(C.BlogID.Eq(int64(BlogId))).Find()
		if err != nil {
			return nil
		}
		for _, comment := range comments {
			idSet = append(idSet, strconv.Itoa(int(comment.CommentID)))
			res = append(res, comment)
		}
		//idSet存入Redis
		utils.Client.SAdd("cache:comment:blog:"+blogId, idSet)
		//返回结果
		return res
	}
	//idSet存在，根据每个id分别找到comment
	for _, i := range idSet {
		res = append(res, GetCommentById(i))
	}
	return res
}

func GetReportedComment() []*model.Comment {
	ctx := context.Background()
	C := query.Comment
	comments, err := C.WithContext(ctx).Where(C.Status.Eq(BoolValue(false))).Find()
	if err != nil {
		return nil
	} else {
		return comments
	}
}

func DeleteCommentById(id int64) bool {
	ctx := context.Background()
	C := query.Comment
	_, err := C.WithContext(ctx).Where(C.CommentID.Eq(id)).Delete()
	if err != nil {
		return false
	} else {
		return true
	}
}

func ChangeStatus(id int64) bool {
	//改变评论状态
	ctx := context.Background()
	C := query.Comment
	comment, err := C.WithContext(ctx).Where(C.CommentID.Eq(id)).First()
	if err != nil {
		return false
	}
	s := *comment.Status
	if s[0] == 1 {
		*comment.Status = []uint8{0}
	} else {
		*comment.Status = []uint8{1}
	}
	err = C.WithContext(ctx).Where(C.CommentID.Eq(id)).Save(comment)
	if err != nil {
		return false
	}
	return true
}

func AddCommentThumb(commentId string, blogId string) bool {
	//添加评论点赞
	key := "hotlist:comment:blogId:" + blogId
	err := utils.Client.ZIncrBy(key, 1, commentId).Err()
	if err != nil {
		return false
	}
	return true
}

func GetHotComments(blogId string, count int64) []interface{} {
	//获取热门评论
	key := "hotlist:comment:blogId:" + blogId
	ids, err := utils.Client.ZRevRangeByScore(key, redis.ZRangeBy{
		Min:   "0",
		Max:   "inf",
		Count: count,
	}).Result()
	if err != nil {
		return nil
	}
	if len(ids) > 0 {
		var res []interface{}
		for _, id := range ids {
			res = append(res, GetCommentById(id))
		}
		return res
	}
	return nil
}
