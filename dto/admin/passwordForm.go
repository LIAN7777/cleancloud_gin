package dto

type AdminPswForm struct {
	Id          int64  `json:"id" ,binding:"required"`
	OldPassword string `json:"old_password" ,binding:"required"`
	NewPassword string `json:"new_password" ,binding:"required"`
}
