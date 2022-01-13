package service

import (
	"github.com/leffss/go-admin/models"
)

func GetUserByUsername(username string) *models.User {
	var user *models.User
	models.DB.Where(&models.User{Username: username}).First(&user)
	return user
}

func GetUserByID(id uint) *models.User {
	var user *models.User
	models.DB.Where(&models.User{BaseModel: models.BaseModel{Id: id}}).First(&user)
	return user
}
