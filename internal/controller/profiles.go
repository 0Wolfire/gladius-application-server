package controller

import (
	"github.com/gladiusio/gladius-application-server/internal/models"
	"github.com/jinzhu/gorm"
)

//func Nodes(db *gorm.DB) ([]models.NodeProfile, error) {
//	var profiles []models.NodeProfile
//
//	err := db.Find(&profiles).Error
//
//	return profiles, err
//}

//func NodesPendingPoolConfirmation(db *gorm.DB) ([]models.NodeProfile, error) {
//	var profiles []models.NodeProfile
//
//	err := db.Where("pool_accepted is ?", nil).Find(&profiles).Error
//
//	return profiles, err
//}
//
//func NodesPendingNodeConfirmation(db *gorm.DB) ([]models.NodeProfile, error) {
//	var profiles []models.NodeProfile
//
//	err := db.Where("pool_accepted = ? AND node_accepted = ?", "true", nil).Find(&profiles).Error
//
//	return profiles, err
//}

func NodesAccepted(db *gorm.DB) ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db.
		Where(&models.NodeProfile{Approved: true, Pending: false}).
		Find(&profiles).Error

	return profiles, err
}

//func NodesRejected(db *gorm.DB) ([]models.NodeProfile, error) {
//	var profiles []models.NodeProfile
//
//	err := db.Where("pool_accepted = ? OR node_accepted = ? OR approved = ?", "false", "false", "false").Find(&profiles).Error
//
//	return profiles, err
//}

func NodeInPool(db *gorm.DB, walletAddress string) (bool, error) {
	var profile models.NodeProfile
	var count int

	err := db.Model(&profile).
		Where("lower(wallet) like lower(?)", walletAddress).
		Where(&models.NodeProfile{Approved: true, Pending: false}).
		//Where("approved = ? AND pending = ?", "1", "0").
		First(&profile).
		Count(&count).Error

	return count > 0, err
}
