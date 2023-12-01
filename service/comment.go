package service

import (
	dto "GinProject/dto/user"
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
	"encoding/json"
	"log"
	"time"
)

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
	err = dComment.WithContext(ctx).Create(comment)
	if err != nil {
		log.Print("add comment fail:caused by\n", err)
		return
	}
	err = utils.Publish("amq.direct", "comment_message", comment)
	if err != nil {
		log.Print("send comment message fail")
	}
}
