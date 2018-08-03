package controller

import (
	"github.com/gladiusio/gladius-application-server/pkg/db/models"
	"github.com/jinzhu/gorm"
)

func Nodes(db *gorm.DB) ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db.Find(&profiles).Error

	return profiles, err
}

func NodesPendingPoolConfirmation(db *gorm.DB) ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db.Where("pool_accepted is ?", nil).Find(&profiles).Error

	return profiles, err
}

func NodesPendingNodeConfirmation(db *gorm.DB) ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db.Where("pool_accepted is ? AND node_accepted is ?", "true", nil).Find(&profiles).Error

	return profiles, err
}

func NodesAccepted(db *gorm.DB) ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db.Where("accepted is ?", "true").Find(&profiles).Error

	return profiles, err
}

func NodesRejected(db *gorm.DB) ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db.Where("pool_accepted is ? OR node_accepted is ? OR accepted is ?", "false", "false", "false").Find(&profiles).Error

	return profiles, err
}
