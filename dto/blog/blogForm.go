package dto

type BlogForm struct {
	UserID    int64   `json:"user_id" ,binding:"required"`
	Title     *string `json:"title" ,binding:"required"`
	Introduce *string `json:"introduce" ,binding:"required"`
	Content   *string `json:"content" ,binding:"required"`
	Image     *string `json:"image"`
	File      *string `json:"file"`
	Time      *string `json:"time"`
	BlogClass *string `json:"blog_class"`
	Tag       *string `json:"tag"`
	UserName  *string `json:"user_name" ,binding:"required"`
}
