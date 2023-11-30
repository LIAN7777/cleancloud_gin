package utils

import (
	"GinProject/model"
	"encoding/json"
	"fmt"
	"testing"
)

func TestInitClient(t *testing.T) {
	err := InitClient()
	if err != nil {
		fmt.Print("wrong")
	} else {
		name := Client.Get("name")
		fmt.Print(name)
	}
}

func TestRedisOp(t *testing.T) {
	err := InitClient()
	defer Client.Close()
	if err != nil {
		fmt.Print("wrong")
	} else {
		name := "lian"
		psw := "777"
		user := &model.User{
			UserID:   666,
			UserName: &name,
			Password: &psw,
		}
		jsonData, err := json.Marshal(user)
		if err != nil {
			fmt.Print("encode error")
		}
		err1 := Client.Set("test:user:2", jsonData, -1).Err()
		if err1 != nil {
			fmt.Print("add user error")
		}
	}
}

func TestRedisGet(t *testing.T) {
	err := InitClient()
	if err != nil {
		fmt.Print("wrong")
	} else {
		jsonData, err := Client.Get("test:user:1").Result()
		if err != nil {
			panic(err)
		}

		// 将 JSON 字符串转换为模型对象
		var user *model.User // 替换成模型类型
		err = json.Unmarshal([]byte(jsonData), &user)
		if err != nil {
			panic(err)
		}
	}
}
