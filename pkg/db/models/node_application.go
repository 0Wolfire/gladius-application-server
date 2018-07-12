package models

import (
	"github.com/jinzhu/gorm"
)

type NodeApplication struct {
	gorm.Model

	IPAddress      string `json:"-" gorm:"not null"`
	EstimatedSpeed int    `json:"estimatedSpeed" gorm:"not null"`
	PoolAccepted   bool   `json:"-" gorm:"not null;default:false"`
	NodeAccepted   bool   `json:"-" gorm:"not null;default:false"`
	Accepted       bool   `json:"-" gorm:"not null;default:false"`
	Wallet         string `json:"wallet" gorm:"not null; unique"`
	Profile        NodeProfile
	ProfileID      uint `json:"-"`
	// Schedule
}

type NodeRequestPayload struct {
	EstimatedSpeed int    `json:"estimatedSpeed"`
	Wallet         string `json:"wallet"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Bio            string `json:"bio"`
	Location       string `json:"location"`
}

func CreateApplication(payload *NodeRequestPayload) NodeApplication {
	// placeholder until REST request can pull ip in
	ipAddress := "0.0.0.0"

	application := NodeApplication{
		IPAddress:      ipAddress,
		EstimatedSpeed: payload.EstimatedSpeed,
		Wallet:         payload.Wallet,
		Profile: NodeProfile{
			Name:     payload.Name,
			Email:    payload.Email,
			Bio:      payload.Bio,
			Location: payload.Location,
		},
	}

	return application
}
