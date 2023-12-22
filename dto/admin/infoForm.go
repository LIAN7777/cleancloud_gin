package dto

type AdminUpdateForm struct {
	Id        int64  `json:"id" ,binding:"required"`
	AdminName string `json:"name" ,binding:"required"`
	Contact   string `json:"contact" ,binding:"required"`
	Address   string `json:"address" ,binding:"required"`
}
