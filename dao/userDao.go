package dao

import (
	"GinProject/model"
	"GinProject/query"
	"context"
)

func GetUserById(userId int64) (*model.User, error) {
	dUser := query.Q.User
	ctx := context.Background()
	user, err := dUser.WithContext(ctx).Where(dUser.UserID.Eq(userId)).First()
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}
