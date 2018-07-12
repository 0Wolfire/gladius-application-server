package controller

import (
	"errors"
	"github.com/gladiusio/gladius-application-server/pkg/db/models"
	"github.com/jinzhu/gorm"
	"log"
)

func db() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

// temp
func TempDBCalls() {
	//Temp for testing
	request := models.NodeRequestPayload{
		EstimatedSpeed: 100,
		Wallet:         "0x975432957943875235",
		Name:           "Name",
		Email:          "email@fds.com",
		Bio:            "bio",
		Location:       "location",
	}

	NodeApplyToPool(request)

	requestUpdate := models.NodeRequestPayload{
		Wallet:   "0x975432957943875235",
		Name:     "Name Updated",
		Email:    "email@fds.com Updated",
		Bio:      "bio Updated",
		Location: "location Updated",
	}

	_, err := NodeUpdateProfile(requestUpdate)
	if err != nil {
		log.Fatal(err)
	}

	poolInfo := models.PoolInformation{
		Name:     "Gladius Pool",
		Address:  "124.232.83.8",
		Bio:      "Gladius Testing Pool",
		Location: "Washington D.C.",
		Rating:   5,
		Wallet:   "0x96585865865",
		Public:   true,
	}

	PoolCreateUpdateData(poolInfo)
}

func PoolCreateUpdateData(poolInfo models.PoolInformation) {
	var pool models.PoolInformation

	db().Model(&pool).FirstOrCreate(&pool)
	db().Model(&pool).Updates(&poolInfo)
}

func NodeApplyToPool(payload models.NodeRequestPayload) {
	application := models.CreateApplication(&payload)

	db().Model(&application).Where("wallet = ?", payload.Wallet).FirstOrCreate(&application)
	db().Save(&application)
}

func NodeUpdateProfile(payload models.NodeRequestPayload) (models.NodeProfile, error) {
	profile, err := NodeProfile(payload.Wallet)
	if err != nil {
		return profile, err
	}

	profile.Name = payload.Name
	profile.Bio = payload.Bio
	profile.Email = payload.Email
	profile.Location = payload.Location

	db().Save(&profile)
	return profile, nil
}

func NodeProfile(wallet string) (models.NodeProfile, error) {
	var application models.NodeApplication
	var profile models.NodeProfile

	if err := db().Model(&application).Where("wallet = ?", wallet).First(&application).Error; err != nil {
		return models.NodeProfile{}, errors.New("NodeProfile() profile not found for given wallet address")
	}
	db().Model(&application).Association("Profile").Find(&profile)
	return profile, nil
}
