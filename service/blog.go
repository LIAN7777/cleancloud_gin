package service

import (
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
	"golang.org/x/sync/singleflight"
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
