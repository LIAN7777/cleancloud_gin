package dto

type GeoRange struct {
	UserId int64   `json:"user_id"`
	Radius float64 `json:"radius"`
	Count  int     `json:"count"`
}
