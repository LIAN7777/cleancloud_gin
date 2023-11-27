package dao

import (
	"GinProject/model"
	"GinProject/utils"
)

func GetCityByName(cityName string) (*model.City, error) {
	city := &model.City{}
	err := utils.DBlink.Where("Name=?", cityName).First(city).Error
	if err != nil {
		return nil, err
	} else {
		return city, nil
	}
}
