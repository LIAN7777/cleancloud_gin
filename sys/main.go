package main

import (
	"GinProject/router"
	"GinProject/service"
	"GinProject/utils"
	"log"
)

func init() {
	err := utils.SetupDBLink()
	if err != nil {
		log.Fatal("db link error")
	}
	err = utils.InitClient()
	if err != nil {
		log.Print("redis link error")
		//log.Fatal("redis link error")
	}
	err = utils.InitRabbit()
	if err != nil {
		log.Print("rabbitmq link error")
	}
	// 启动goroutine监听队列
	go service.StartConsumer("comment_queue", service.PrintComment)
}

func main() {
	//启动gin实例
	r := router.Router()
	err := r.Run(":8888")
	if err != nil {
		log.Printf("server run error")
		return
	}
}
