package controller

import (
	"github.com/gladiusio/gladius-application-server/pkg/db/models"
	)

func Nodes() ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db().Find(&profiles).Error

	return profiles, err
}

func NodesPendingPoolConfirmation() ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db().Where("pool_accepted is ?", nil).Find(&profiles).Error

	return profiles, err
}

func NodesPendingNodeConfirmation() ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db().Where("node_accepted is ?", nil).Find(&profiles).Error

	return profiles, err
}

func NodesAccepted() ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db().Where("accepted = ?", true).Find(&profiles).Error

	return profiles, err
}

func NodesRejected() ([]models.NodeProfile, error) {
	var profiles []models.NodeProfile

	err := db().Where("accepted = ?", false).Find(&profiles).Error

	return profiles, err
}