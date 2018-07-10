package controller

import (
	"github.com/gladiusio/gladius-application-server/pkg/db/models"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func NodeApplyToPool(payload models.NodeRequestPayload) {
	application := models.CreateApplication(&payload)

	DB.Create(&application)
	//DB.Save(application)

	//DB.Create(&application.Profile)
}
