package utils

import (
	"GinProject/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBlink *gorm.DB

func SetupDBLink() error {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/cleancloud?charset=utf8&parseTime=True&loc=Local"
	DBlink, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err == nil {
		query.SetDefault(DBlink)
		return nil
	} else {
		return err
	}
}
