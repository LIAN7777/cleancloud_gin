package dto

type CommentForm struct {
	UserId  int64  `json:"userId"`
	BlogId  int64  `json:"blogId"`
	Content string `json:"content"`
}
