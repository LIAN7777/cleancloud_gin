// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameAdminmessage = "adminmessage"

// Adminmessage mapped from table <adminmessage>
type Adminmessage struct {
	AdminID   *int64  `gorm:"column:admin_id;type:int" json:"admin_id"`
	Time      *string `gorm:"column:time;type:varchar(255)" json:"time"`
	Source    *string `gorm:"column:source;type:varchar(255)" json:"source"`
	Content   *string `gorm:"column:content;type:varchar(255)" json:"content"`
	MessageID int64   `gorm:"column:message_id;type:int;primaryKey;autoIncrement:true" json:"message_id"`
}

// TableName Adminmessage's table name
func (*Adminmessage) TableName() string {
	return TableNameAdminmessage
}
