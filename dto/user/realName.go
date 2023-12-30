package dto

type RealName struct {
	Id     int64  `json:"id" ,binding:"required"`
	Name   string `json:"real_name" ,binding:"required"`
	RealId string `json:"real_id" ,binding:"required"`
}
