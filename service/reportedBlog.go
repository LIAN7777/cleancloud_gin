package service

import (
	dto "GinProject/dto/report"
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
	"errors"
	"strconv"
	"time"
)

func AddReportedBlog(blogId int64) (bool, error) {
	ctx := context.Background()
	B := query.Blog
	//获取原博客
	blog, err := B.WithContext(ctx).Where(B.BlogID.Eq(blogId)).First()
	if err != nil || blog == nil {
		return false, errors.New("blog not exist")
	}
	//删除原博客以及redis缓存
	//开启事务
	tx := utils.DBlink.Begin()
	err = tx.WithContext(ctx).Delete(blog).Error
	if err != nil {
		tx.Rollback()
		return false, errors.New("delete blog fail")
	}
	_ = utils.Client.Del("cache:blog:" + strconv.FormatInt(blogId, 10))
	//添加被举报博客
	err = tx.WithContext(ctx).Create(Blog2Report(blog)).Error
	if err != nil {
		tx.Rollback()
		return false, errors.New("add reported blog fail")
	}
	tx.Commit()
	return true, nil
}

func Blog2Report(blog *model.Blog) *model.Reportedblog {
	//将博客转化为被举报
	reported := &model.Reportedblog{}
	reported.UserID = blog.UserID
	reported.UserName = blog.UserName
	reported.Introduce = blog.Introduce
	reported.Content = blog.Content
	reported.BlogClass = blog.BlogClass
	reported.File = blog.File
	reported.Image = blog.Image
	reported.Tag = blog.Tag
	reported.Title = blog.Title
	status := "0"
	reported.Status = &status
	now := time.Now().Format("2006-01-02 15:04:05")
	reported.Time = &now
	return reported
}

func DeleteReportedBlog(id int64) bool {
	ctx := context.Background()
	R := query.Reportedblog
	_, err := R.WithContext(ctx).Where(R.RBlogID.Eq(id)).Delete()
	return err == nil
}

func GetReportedBlogById(id int64) *model.Reportedblog {
	ctx := context.Background()
	R := query.Reportedblog
	blog, err := R.WithContext(ctx).Where(R.RBlogID.Eq(id)).First()
	if err != nil {
		return nil
	}
	return blog
}

func AddAssistantComment(comment *dto.AssistantComment) (bool, error) {
	ctx := context.Background()
	R := query.Reportedblog
	blog, err := R.WithContext(ctx).Where(R.RBlogID.Eq(comment.BlogId)).First()
	if err != nil {
		return false, errors.New("blog not exist")
	}
	if blog.AssistantComment == nil {
		_, err = R.WithContext(ctx).Where(R.RBlogID.Eq(comment.BlogId)).Update(R.AssistantComment, comment.Comment)
		if err != nil {
			return false, errors.New("add comment fail")
		}
		return true, nil
	}
	newComment := *blog.AssistantComment + "\n" + comment.Comment
	_, err = R.WithContext(ctx).Where(R.RBlogID.Eq(comment.BlogId)).Update(R.AssistantComment, newComment)
	if err != nil {
		return false, errors.New("add comment fail")
	}
	return true, nil
}
