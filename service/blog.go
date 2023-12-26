package service

import (
	dto "GinProject/dto/blog"
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"golang.org/x/sync/singleflight"
	"log"
	"strconv"
	"time"
)

// 缓存互斥锁，用于解决缓存击穿问题
// var cacheMutex = &sync.Mutex{}
var g = &singleflight.Group{}

func GetBlogById(id string) interface{} {
	//该问题也可以使用基于redis的分布式锁来实现
	//单飞模式解决互相等待问题
	blog, err, _ := g.Do(id, func() (interface{}, error) {
		//读取redis
		blog, err := utils.RedisGetModel("cache:blog:"+id, model.Blog{})
		//err为nil: 1.得到空值 2.得到数据;均直接返回blog，上层controller进行判断
		if err == nil {
			return blog, nil
		} else {
			//不存在则读mysql
			dBlog := query.Blog
			ctx := context.Background()
			blogId, _ := strconv.Atoi(id)
			blog, err := dBlog.WithContext(ctx).Where(dBlog.BlogID.Eq(int64(blogId))).First()
			//不存在返回
			if err != nil {
				//向redis存入空值防止缓存穿透
				utils.Client.Set("cache:blog:"+id, "", time.Minute*30)
				return nil, err
			} else {
				//存在写redis，返回
				utils.RedisSetModel("cache:blog:"+id, blog)
				return blog, nil
			}
		}
	})

	if err != nil {
		return nil
	} else {
		return blog
	}
	//一般的互斥锁方案，存在问题：当goroutine的数量大大增加时，会有很多goroutine在等待，大大降低了服务器性能，
	//按理来说当A释放锁以后，说明数据已经被缓存，其他等待的goroutine不需要再等待依次获取锁了，而是直接去访问缓存即可
	////读取redis
	//blog, err := utils.RedisGetModel("cache:blog:"+id, model.Blog{})
	////err为nil: 1.得到空值 2.得到数据;均直接返回blog，上层controller进行判断
	//if err == nil {
	//	return blog
	//}
	////缓存不存在，加锁，高并发下只有第一个线程拿到锁
	//cacheMutex.Lock()
	//defer cacheMutex.Unlock()
	//
	////锁结束后再次尝试获取redis，因为可能redis已经被其他线程更新，就无需再读数据库
	//blog, err = utils.RedisGetModel("cache:blog:"+id, model.Blog{})
	//if err == nil {
	//	return blog
	//} else {
	//	//不存在则读mysql
	//	dBlog := query.Blog
	//	ctx := context.Background()
	//	blogId, _ := strconv.Atoi(id)
	//	blog, err := dBlog.WithContext(ctx).Where(dBlog.BlogID.Eq(int64(blogId))).First()
	//	//不存在返回
	//	if err != nil {
	//		//向redis存入空值防止缓存穿透
	//		utils.Client.Set("cache:blog:"+id, "", time.Minute*30)
	//		return nil
	//	} else {
	//		//存在写redis，返回
	//		utils.RedisSetModel("cache:blog:"+id, blog)
	//		return blog
	//	}
	//}
}

func UpdateBlog(blog *model.Blog) bool {
	id := blog.BlogID
	//开启事务，保证数据库操作和redis的一致性
	tx := utils.DBlink.Begin()
	//更新数据库
	ctx := context.Background()
	err := tx.WithContext(ctx).Updates(blog).Error
	if err != nil {
		tx.Rollback()
		return false
	}
	//删除redis缓存
	err = utils.Client.Del("cache:blog:" + strconv.Itoa(int(id))).Err()
	if err != nil {
		tx.Rollback()
		return false
	} else {
		tx.Commit()
		return true
	}
}

func GetBlogThumb(id string) int64 {
	//先查询Redis
	res, err := utils.Client.Get("cache:blog:thumb:" + id).Result()
	if err == nil {
		count, _ := strconv.Atoi(res)
		return int64(count)
	}
	//Redis无数据,查MySQL
	ctx := context.Background()
	dThumb := query.Thumb
	blogId, _ := strconv.Atoi(id)
	count, err := dThumb.WithContext(ctx).Where(dThumb.BlogID.Eq(int64(blogId))).Count()
	if err != nil {
		return 0
	}
	//写入Redis缓存
	utils.Client.Set("cache:blog:thumb:"+id, count, time.Minute*30)
	//返回值
	return count
}

func PublishBlogThumb(thumb *model.Thumb) bool {
	//先查看是否存在Redis，若存在则将Redis值自增，方便查询
	blogId := strconv.Itoa(int(thumb.BlogID))
	err := utils.Client.Get("cache:blog:thumb:" + blogId).Err()
	if err == nil {
		//存在redis，自增
		utils.Client.Incr("cache:blog:thumb:" + blogId)
	}
	//向rabbitmq发送消息，消费者消费后新增数据
	err = utils.Publish("amq.direct", "thumb", thumb)
	if err != nil {
		return false
	}
	return true
}

func AddBlogThumb(msg []byte) {
	thumb := &model.Thumb{}
	err := json.Unmarshal(msg, thumb)
	if err != nil {
		log.Print("thumb format error:", err)
	}
	err = query.Thumb.WithContext(context.Background()).Create(thumb)
	if err != nil {
		log.Print("add thumb error")
		//TODO:重新投递消息
	}
}

func GetBlogByUserFavorite(userId string) []interface{} {
	var blogs []interface{}
	var idSet []string
	//先到Redis中查询用户收藏的博客id
	res, err := utils.Client.SMembers("cache:user:favorite:" + userId).Result()
	if err == nil && cap(res) != 0 {
		idSet = res
	} else {
		//Redis中不存在，则到MySQL查
		dFavor := query.Favorite
		id, _ := strconv.Atoi(userId)
		favors, err := dFavor.WithContext(context.Background()).Where(dFavor.UserID.Eq(int64(id))).Find()
		if err != nil {
			return nil
		}
		//blogId写入缓存
		for _, favor := range favors {
			idSet = append(idSet, strconv.Itoa(int(favor.BlogID)))
		}
		err = utils.Client.SAdd("cache:user:favorite:"+userId, idSet).Err()
		if err != nil {
			log.Print("user favor cache add fail")
		}
	}
	//根据博客id查询博客
	for _, blogId := range idSet {
		//先在Redis中查询指定id的博客
		blog, err := utils.RedisGetModel("cache:blog:"+blogId, model.Blog{})
		if err == nil {
			blogs = append(blogs, blog)
			continue
		}
		//未查询到则到MySQL中查
		id, _ := strconv.Atoi(blogId)
		blog, err = query.Blog.WithContext(context.Background()).Where(query.Blog.BlogID.Eq(int64(id))).First()
		if err == nil {
			blogs = append(blogs, blog)
			//写入缓存
			if ok := utils.RedisSetModel("cache:blog:"+blogId, blog); !ok {
				log.Print("add blog cache fail")
			}
		}
	}
	return blogs
}

func GetBlogByUserId(userId string) []interface{} {
	var idSet []string
	var blogRes []interface{}
	//查redis用户的博客idSet
	res, err := utils.Client.SMembers("cache:user:blog:" + userId).Result()
	if err == nil && cap(res) != 0 {
		idSet = res
	} else {
		//没查到，查MySQL
		q := query.Blog
		ctx := context.Background()
		id, _ := strconv.Atoi(userId)
		blogs, err := q.WithContext(ctx).Where(q.UserID.Eq(int64(id))).Find()
		if err != nil {
			return nil
		}
		for _, blog := range blogs {
			blogId := strconv.Itoa(int(blog.BlogID))
			idSet = append(idSet, blogId)
			//将每个博客写入缓存
			utils.RedisSetModel("cache:blog:"+blogId, blog)
			blogRes = append(blogRes, blog)
		}
		//idSet存入redis
		utils.Client.SAdd("cache:user:blog:"+userId, idSet)
		//直接返回结果
		return blogRes
	}
	//redis中有idSet缓存根据idSet查询博客
	for _, id := range idSet {
		//先查redis
		blog, err := utils.RedisGetModel("cache:blog:"+id, model.Blog{})
		if blog != nil && err == nil {
			blogRes = append(blogRes, blog)
		} else {
			//查MySQL
			q := query.Blog
			ctx := context.Background()
			blogId, _ := strconv.Atoi(id)
			blog, err = q.WithContext(ctx).Where(q.BlogID.Eq(int64(blogId))).First()
			if err != nil {
				blogRes = append(blogRes, blog)
				//把博客存入redis
				utils.RedisSetModel("cache:blog:"+id, blog)
			}
		}
	}
	return blogRes
}

func PublishBlog(blog *dto.BlogForm) bool {
	err := utils.Publish("amq.direct", "blog", blog)
	if err != nil {
		return false
	}
	return true
}

func AddUnreviewedBlog(msg []byte) {
	blog := &model.Reportedblog{}
	err := json.Unmarshal(msg, blog)
	if err != nil {
		log.Print("blog format error")
		return
	}
	//TODO:AI审核
	status := "1"
	blog.Status = &status
	ctx := context.Background()
	B := query.Reportedblog
	err = B.WithContext(ctx).Create(blog)
	if err != nil {
		log.Print("add new blog fail:caused by\n", err)
		return
	}
}

func AddBlogHits(id string) bool {
	//增加博客点击量；考虑并发问题，采用分布式锁
	lock := utils.NewRedisLock("lock:blogHit", utils.Client)
	ctx := context.Background()
	_ = lock.Lock(ctx)
	defer func(lock *utils.RedisLock) {
		_ = lock.Unlock()
	}(lock)
	//添加点击量
	_, err := utils.Client.ZIncrBy("hotlist:blog", 1, id).Result()
	if err != nil {
		return false
	}
	_ = lock.Unlock()
	//添加一分钟内的点击量
	//为了和删除点击量互斥，需要上锁
	lock1 := utils.NewRedisLock("lock:blogHit_delete", utils.Client)
	_ = lock1.Lock(ctx)
	_, err = utils.Client.ZIncrBy("hotlist_delete:blog", 1, id).Result()
	_ = lock1.Unlock()
	return true
}

func GetHotBlogs(limit int) []interface{} {
	//获取限定数量的热门博客
	//获取热门博客的id集合
	res, err := utils.Client.ZRevRangeByScore("hotlist:blog", redis.ZRangeBy{
		Min:   "0",
		Max:   "inf",
		Count: int64(limit),
	}).Result()
	if err != nil {
		return nil
	}
	//根据id集合去获取博客信息
	var blogs []interface{}
	for _, blogId := range res {
		blogs = append(blogs, GetBlogById(blogId))
	}
	return blogs
}

func DecrBlogHits() {
	//启用定时任务，每分钟收集一次博客点击量，发送到延时队列中，一小时后删除
	for range time.Tick(time.Minute) {
		//上锁
		ctx := context.Background()
		lock := utils.NewRedisLock("lock:blogHit_delete", utils.Client)
		_ = lock.Lock(ctx)
		//获取一分钟内点击量
		res, _ := utils.Client.ZRangeWithScores("hotlist_delete:blog", 0, -1).Result()
		//删除一分钟内的点击量
		_ = utils.Client.Del("hotlist_delete:blog")
		_ = lock.Unlock()
		if res != nil {
			//res不为空，则将每个键发送到延时队列
			for _, z := range res {
				message := dto.BlogHits{BlogId: z.Member, Hits: int64(z.Score)}
				_ = utils.DelayPublish("delayed_exchange", "blog_hits", message, 3600000)
			}
		}
	}
}

func DeleteBlogHits(msg []byte) {
	hits := &dto.BlogHits{}
	err := json.Unmarshal(msg, hits)
	if err != nil {
		log.Print("blogHits json invalid")
	}
	//在真正的博客点击量中减少点击数
	_, err = utils.Client.ZIncrBy("hotlist:blog", -float64(hits.Hits), hits.BlogId.(string)).Result()
}
