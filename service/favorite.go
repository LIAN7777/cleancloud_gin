package service

import (
	"GinProject/model"
	"GinProject/query"
	"context"
)

func AddFavorite(favor *model.Favorite) bool {
	ctx := context.Background()
	F := query.Favorite
	err := F.WithContext(ctx).Create(favor)
	if err != nil {
		return false
	}
	return true
}

func DeleteFavorite(favor *model.Favorite) bool {
	ctx := context.Background()
	F := query.Favorite
	_, err := F.WithContext(ctx).Delete(favor)
	if err != nil {
		return false
	}
	return true
}

func JudgeFavorite(favor *model.Favorite) bool {
	ctx := context.Background()
	F := query.Favorite
	count, err := F.WithContext(ctx).Where(F.UserID.Eq(favor.UserID), F.BlogID.Eq(favor.BlogID)).Count()
	if err != nil {
		return false
	}
	return count > 0
}
