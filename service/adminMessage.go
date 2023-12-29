package service

import (
	"GinProject/model"
	"GinProject/query"
	"context"
	"time"
)

func GetMessageByAdmin(adminId int64) []*model.Adminmessage {
	ctx := context.Background()
	M := query.Adminmessage
	messages, err := M.WithContext(ctx).Where(M.AdminID.Eq(adminId)).Find()
	if err != nil {
		return nil
	}
	return messages
}

func AddAdminMessage(message *model.Adminmessage) bool {
	ctx := context.Background()
	M := query.Adminmessage
	now := time.Now().Format("2006-01-02 15:04:05")
	message.Time = &now
	err := M.WithContext(ctx).Create(message)
	if err != nil {
		return false
	}
	return true
}

func DeleteAdminMessage(id int64) bool {
	ctx := context.Background()
	M := query.Adminmessage
	_, err := M.WithContext(ctx).Where(M.MessageID.Eq(id)).Delete()
	if err != nil {
		return false
	}
	return true
}
