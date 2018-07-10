package models

import (
	"github.com/jinzhu/gorm"
)

type PoolInformation struct {
	gorm.Model
	OnlyRow   int    `json:"-" gorm:"unique;not null;default:1"`
	Address   string `json:"address" gorm:"not null"`
	Name      string `json:"name" gorm:"not null"`
	Bio       string `json:"bio"`
	Location  string `json:"location" gorm:"not null"`
	Rating    int    `json:"rating"`
	NodeCount int    `json:"nodeCount" gorm:"not null;default:0"`
	Wallet    string `json:"wallet" gorm:"not null"`
	Public    bool   `json:"public" grom:"not null;default:false"`
}
