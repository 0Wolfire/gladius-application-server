package models

import (
	"github.com/jinzhu/gorm"
)

type NodeProfile struct {
	gorm.Model

	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Bio      string `json:"bio" gorm:"not null"`
	Location string `json:"location" gorm:"not null"`
	Wallet   string `json:"wallet" gorm:"not null"`
}
