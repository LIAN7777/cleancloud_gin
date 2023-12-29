package dto

type UserUpdateForm struct {
	Id         int64  `json:"id" ,binding:"required"`
	UserName   string `json:"user_name" ,binding:"required"`
	Address    string `json:"address"`
	Introduce  string `json:"introduce"`
	Profession string `json:"profession"`
	Age        int64  `json:"age"`
}
