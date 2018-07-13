package models

import (
	"github.com/jinzhu/gorm"
)

type NodeProfile struct {
	gorm.Model

	Name           string `json:"name" gorm:"not null"`
	Email          string `json:"email" gorm:"not null"`
	Bio            string `json:"bio" gorm:"not null"`
	Location       string `json:"location" gorm:"not null"`
	IPAddress      string `json:"-" gorm:"not null"`
	EstimatedSpeed int    `json:"estimatedSpeed" gorm:"not null"`
	PoolAccepted   bool   `json:"-" gorm:"not null;default:false"`
	NodeAccepted   bool   `json:"-" gorm:"not null;default:false"`
	Accepted       bool   `json:"-" gorm:"not null;default:false"`
	Wallet         string `json:"wallet" gorm:"not null; unique"`
}

type NodeRequestPayload struct {
	EstimatedSpeed int    `json:"estimatedSpeed"`
	Wallet         string `json:"wallet"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Bio            string `json:"bio"`
	Location       string `json:"location"`
}

func CreateApplication(payload *NodeRequestPayload) NodeProfile {
	// placeholder until REST request can pull ip in
	ipAddress := "0.0.0.0"

	profile := NodeProfile{
		IPAddress:      ipAddress,
		EstimatedSpeed: payload.EstimatedSpeed,
		Wallet:         payload.Wallet,
		Name:           payload.Name,
		Email:          payload.Email,
		Bio:            payload.Bio,
		Location:       payload.Location,
	}

	return profile
}

func (profile *NodeProfile) AfterUpdate(tx *gorm.DB) (err error) {
	if profile.Accepted != (profile.PoolAccepted && profile.NodeAccepted) {
		tx.Model(&NodeProfile{}).Where("id = ?", profile.ID).Update("accepted", profile.PoolAccepted && profile.NodeAccepted)
	}
	return
}
