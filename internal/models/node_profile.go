package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type NodeProfile struct {
	// Default Gorm Model Properties
	// Stating here to remove from JSON responses with `json:"-"`
	ID             uint       `json:"-" gorm:"primary_key"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `json:"-"`
	Name           string     `json:"name" gorm:"not null"`
	Email          string     `json:"email" gorm:"not null"`
	Bio            string     `json:"bio" gorm:"not null"`
	Location       string     `json:"location" gorm:"not null"`
	IPAddress      string     `json:"-" gorm:"not null"`
	EstimatedSpeed int        `json:"estimatedSpeed" gorm:"not null"`
	PoolAccepted   bool       `json:"-" gorm:"default:false"`
	NodeAccepted   bool       `json:"-" gorm:"default:false"`
	Pending        bool       `json:"pending" gorm:"default:true"`
	Approved       bool       `json:"approved" gorm:"default:false"`
	Wallet         string     `json:"wallet" gorm:"not null; unique"`
}

type NodeRequestPayload struct {
	EstimatedSpeed int    `json:"estimatedSpeed"`
	Wallet         string `json:"wallet"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Bio            string `json:"bio"`
	Location       string `json:"location"`
	IPAddress      string `json:"ipAddress"`
}

func CreateApplication(payload *NodeRequestPayload) NodeProfile {
	profile := NodeProfile{
		IPAddress:      payload.IPAddress,
		EstimatedSpeed: payload.EstimatedSpeed,
		Wallet:         strings.ToLower(payload.Wallet),
		Name:           payload.Name,
		Email:          payload.Email,
		Bio:            payload.Bio,
		Location:       payload.Location,
	}

	return profile
}

func (profile *NodeProfile) AfterUpdate(db *gorm.DB) (err error) {
	if (profile.Approved != (profile.PoolAccepted && profile.NodeAccepted)) {
		db.First(&profile)
		profile.Approved = (profile.PoolAccepted && profile.NodeAccepted)
		profile.Pending = !(profile.PoolAccepted && profile.NodeAccepted)

		db.Save(&profile)
		log.Debug().Msg("Wallet: " + profile.Wallet + " application status updated")
	}

	if profile.Wallet != strings.ToLower(profile.Wallet) {
		db.Model(&NodeProfile{}).Where("id like ?", profile.Wallet).Update("wallet", strings.ToLower(profile.Wallet))
	}

	return
}
