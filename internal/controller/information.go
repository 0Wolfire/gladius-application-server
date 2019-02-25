package controller

import (
	"github.com/gladiusio/gladius-application-server/internal/models"
	"github.com/jinzhu/gorm"
)

func PoolInformation(db *gorm.DB) (models.PoolInformation, error) {
	var profile models.PoolInformation

	err := db.First(&profile).Error

	return profile, err
}

func PoolWalletOwner(db *gorm.DB, walletAddress string) (bool, error) {
	var profile models.PoolInformation
	var count = 0

	err := db.Model(&profile).
		Where("lower(wallet) like lower(?)", walletAddress).
		First(&profile).
		Count(&count).Error

	return count > 0, err
}
