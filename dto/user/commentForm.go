package dto

type CommentForm struct {
	UserId   int64  `json:"user_id"`
	BlogId   int64  `json:"blog_id"`
	Content  string `json:"content"`
	UserName string `json:"user_name"`
}
