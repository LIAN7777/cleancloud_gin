package service

import (
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
)

func AddReport(report *model.Report) bool {
	ctx := context.Background()
	R := query.Report
	B := query.Blog
	userId := report.UserID
	blogId := report.BlogID
	count, err := R.WithContext(ctx).Where(R.UserID.Eq(userId), R.BlogID.Eq(blogId)).Count()
	//举报已存在
	if err != nil || count > 0 {
		return false
	}
	user := GetUserById(userId)
	blog, err := B.WithContext(ctx).Where(B.BlogID.Eq(blogId)).First()
	//用户或博客不存在
	if err != nil || user == nil || blog == nil {
		return false
	}
	//保存举报
	//开启事务
	tx := utils.DBlink.Begin()
	err = tx.Create(report).Error
	if err != nil {
		tx.Rollback()
		return false
	}
	//修改博客举报数
	_, err = B.WithContext(ctx).Where(B.BlogID.Eq(blogId)).Update(B.Reports, *blog.Reports+1)
	if err != nil {
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func JudgeReport(report *model.Report) bool {
	ctx := context.Background()
	R := query.Report
	uid := report.UserID
	bid := report.BlogID
	count, err := R.WithContext(ctx).Where(R.UserID.Eq(uid), R.BlogID.Eq(bid)).Count()
	if err != nil || count == 0 {
		return false
	}
	return true
}
