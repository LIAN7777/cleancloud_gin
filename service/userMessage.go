package service

import (
	"GinProject/model"
	"GinProject/query"
	"context"
	"encoding/json"
	"log"
	"time"
)

func ReceiveComment(msg []byte) {
	comment := &model.Comment{}
	err := json.Unmarshal(msg, comment)
	if err != nil {
		log.Print("comment format error")
		return
	}
	blogId := comment.BlogID
	dBlog := query.Blog
	user, err := dBlog.WithContext(context.Background()).Select(dBlog.UserID).Where(dBlog.BlogID.Eq(blogId)).First()
	message := model.Usermessage{}
	t := time.Now().Format("2006-01-02 15:04:05")
	message.Time = &t
	message.UserID = &user.UserID
	content := "您收到一条新的评论，请及时查看：\n" + *comment.Content
	message.Content = &content
	source := "系统"
	message.Source = &source
	dUMessage := query.Usermessage
	err = dUMessage.WithContext(context.Background()).Create(&message)
	if err != nil {
		log.Print("receive user message fail")
	}
}

func GetMessageByUser(userId int64) []*model.Usermessage {
	ctx := context.Background()
	M := query.Usermessage
	messages, err := M.WithContext(ctx).Where(M.UserID.Eq(userId)).Find()
	if err != nil {
		return nil
	}
	return messages
}

func AddUserMessage(message *model.Usermessage) bool {
	ctx := context.Background()
	M := query.Usermessage
	now := time.Now().Format("2006-01-02 15:04:05")
	message.Time = &now
	err := M.WithContext(ctx).Create(message)
	if err != nil {
		return false
	}
	return true
}

func DeleteUserMessage(id int64) bool {
	ctx := context.Background()
	M := query.Usermessage
	_, err := M.WithContext(ctx).Where(M.UserMessageID.Eq(id)).Delete()
	if err != nil {
		return false
	}
	return true
}
