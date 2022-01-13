package service

import "github.com/leffss/go-admin/models"

func UpdatesUser(user *models.User, values map[string]interface{})  {
	models.DB.Model(user).Updates(values)
}

