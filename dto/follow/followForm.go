package dto

type FollowForm struct {
	UserId    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	Image     string `json:"image"`
	Introduce string `json:"introduce"`
	Field     string `json:"field"`
	BlogNum   int64  `json:"blog_num"`
}
