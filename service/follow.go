package service

import (
	dto "GinProject/dto/follow"
	"GinProject/model"
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

func AddFollow(follow *model.Follow) bool {
	ctx := context.Background()
	F := query.Follow
	err := F.WithContext(ctx).Create(follow)
	if err != nil {
		return false
	}
	return true
}

func DeleteFollow(follow *model.Follow) bool {
	ctx := context.Background()
	F := query.Follow
	_, err := F.WithContext(ctx).Delete(follow)
	if err != nil {
		return false
	}
	return true
}

func JudgeFollow(follow *model.Follow) bool {
	ctx := context.Background()
	F := query.Follow
	count, err := F.WithContext(ctx).Where(F.UserID.Eq(follow.UserID), F.FollowID.Eq(follow.FollowID)).Count()
	if err != nil {
		return false
	}
	return count > 0
}
