package db

import (
	"github.com/gladiusio/gladius-application-server/pkg/controller"
	"github.com/gladiusio/gladius-application-server/pkg/db/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite for now, SQL later on
)

func Initialize() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&models.PoolInformation{})
	db.AutoMigrate(&models.NodeApplication{})
	db.AutoMigrate(&models.NodeProfile{})

	controller.DB = db

	// Temp for testing
	request := models.NodeRequestPayload{
		EstimatedSpeed: 100,
		Wallet:         "0x975432957943875235",
		Name:           "Name",
		Email:          "email@fds.com",
		Bio:            "bio",
		Location:       "location",
	}

	controller.NodeApplyToPool(request)
}
