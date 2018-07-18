package models

import (
	"github.com/jinzhu/gorm"
	"database/sql"
)

type NodeProfile struct {
	gorm.Model

	Name           string       `json:"name" gorm:"not null"`
	Email          string       `json:"email" gorm:"not null"`
	Bio            string       `json:"bio" gorm:"not null"`
	Location       string       `json:"location" gorm:"not null"`
	IPAddress      string       `json:"-" gorm:"not null"`
	EstimatedSpeed int          `json:"estimatedSpeed" gorm:"not null"`
	PoolAccepted   sql.NullBool `json:"-" gorm:"default:null"`
	NodeAccepted   sql.NullBool `json:"-" gorm:"default:null"`
	Accepted       sql.NullBool `json:"-" gorm:"default:null"`
	Wallet         string       `json:"wallet" gorm:"not null; unique"`
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
		Wallet:         payload.Wallet,
		Name:           payload.Name,
		Email:          payload.Email,
		Bio:            payload.Bio,
		Location:       payload.Location,
	}

	return profile
}

func (profile *NodeProfile) AfterUpdate(tx *gorm.DB) (err error) {
	if profile.Accepted.Bool != (profile.PoolAccepted.Bool && profile.NodeAccepted.Bool) {
		tx.Model(&NodeProfile{}).Where("id = ?", profile.ID).Update("accepted", profile.PoolAccepted.Bool && profile.NodeAccepted.Bool)
	}

	return
}
