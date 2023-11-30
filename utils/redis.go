package utils

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

var Client *redis.Client

func InitClient() error {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     "59.110.54.159:6379",
		Password: "020109", // no password set
		DB:       0,        // use default DB
		PoolSize: 10,       // connection pool size
	})

	_, err := Rdb.Ping().Result()
	if err != nil {
		Client = nil
		return err
	}
	Client = Rdb
	return nil
}

func RedisGetModel(key string, model interface{}) (interface{}, error) {
	jsonData, err := Client.Get(key).Result()
	//不存在缓存
	if err != nil {
		return nil, err
	}
	//得到空值，返回空对象，service得到空对象后直接返回而不查询数据库
	if jsonData == "" {
		return nil, nil
	}
	// 将 JSON 字符串转换为模型对象
	err = json.Unmarshal([]byte(jsonData), &model)
	if err != nil {
		return nil, err
	}
	//若命中，增加持续时间
	Client.Expire(key, time.Minute*60)
	return model, nil
}

func RedisSetModel(key string, model interface{}) bool {
	jsonData, err := json.Marshal(model)
	if err != nil {
		return false
	}
	err1 := Client.Set(key, jsonData, time.Minute*60).Err()
	if err1 != nil {
		return false
	}
	return true
}
