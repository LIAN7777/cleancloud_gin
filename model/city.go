package model

type City struct {
	ID          int    `gorm:"column:ID" ,json:"id"`
	Name        string `gorm:"column:Name" ,json:"name"`
	CountryCode string `gorm:"column:CountryCode" ,json:"code"`
	District    string `gorm:"column:District" ,json:"district"`
	Population  int    `gorm:"column:Population" ,json:"population"`
}

func (City) TableName() string {
	return "City"
}
