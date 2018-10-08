package controller

import (
	"github.com/gladiusio/gladius-common/pkg/db/models"
	"github.com/jinzhu/gorm"
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
