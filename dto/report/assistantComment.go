package dto

type AssistantComment struct {
	BlogId  int64  `json:"blog_id" ,binding:"required"`
	Comment string `json:"comment" ,binding:"required"`
}
