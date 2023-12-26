package service

import (
	dto "GinProject/dto/follow"
	"GinProject/query"
	"context"
)

func getStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func GetFollowByUser(userId int64) []*dto.FollowForm {
	ctx := context.Background()
	F := query.Follow
	follows, err := F.WithContext(ctx).Where(F.UserID.Eq(userId)).Find()
	if err != nil {
		return nil
	}
	var res []*dto.FollowForm
	for _, follow := range follows {
		id := follow.FollowID
		user, err := query.User.WithContext(ctx).Where(query.User.UserID.Eq(id)).First()
		if err == nil {
			form := dto.FollowForm{
				UserId:    user.UserID,
				UserName:  getStringValue(user.UserName),
				Introduce: getStringValue(user.Introduce),
				Field:     getStringValue(user.Profession),
				BlogNum:   0,
			}
			count, err := query.Blog.WithContext(ctx).Where(query.Blog.UserID.Eq(user.UserID)).Count()
			if err == nil {
				form.BlogNum = count
			}
			res = append(res, &form)
		}
	}
	return res
}
