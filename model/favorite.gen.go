// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameFavorite = "favorite"

// Favorite mapped from table <favorite>
type Favorite struct {
	UserID int64 `gorm:"column:user_id;type:int;primaryKey" json:"user_id"`
	BlogID int64 `gorm:"column:blog_id;type:int;primaryKey" json:"blog_id"`
}

// TableName Favorite's table name
func (*Favorite) TableName() string {
	return TableNameFavorite
}
