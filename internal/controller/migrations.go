package controller

import (
	"github.com/gladiusio/gladius-application-server/internal/models"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func Initialize(db *gorm.DB) (*gorm.DB, error) {
	var err error

	if db == nil {
		db, err = gorm.Open("sqlite3", "test.db")
	}

	if err != nil {
		return db, err
	}

	// Migrate the schemas
	db.AutoMigrate(&models.PoolInformation{})
	db.AutoMigrate(&models.NodeProfile{})

	return db, err
}

func InitializePoolManager(db *gorm.DB) {
	poolInfo := models.PoolInformation{
		Address:"0x0000000000000000000000000000000000000000",
		Bio:"Pool Initialized with Default Values",
		Email:"pool@test-values.com",
		Location:"Placeholder, PH",
		Name: "Place Holder",
		Rating: 0,
		Public:false,
		Url:"localhost:" + viper.GetString("API.Port"),
		Wallet: "0x0000000000000000000000000000000000000000",
		NodeCount:0,
	}
	db.NewRecord(poolInfo)
	db.Create(&poolInfo)
}
