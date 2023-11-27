package main

import (
	"GinProject/router"
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
}

func main() {
	r := router.Router()
	err := r.Run(":8888")
	if err != nil {
		log.Printf("server run error")
		return
	}
}
