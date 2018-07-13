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

	// Pool Accepts Application
	PoolApplicationStatus("0x975432957943875235", true)
	// Node Denies Application
	NodeApplicationStatus("0x975432957943875235", false)

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
	profile := models.CreateApplication(&payload)
	db().Model(&profile).Where("wallet = ?", payload.Wallet).FirstOrCreate(&profile)
}

func NodeUpdateProfile(payload models.NodeRequestPayload) (models.NodeProfile, error) {
	profile, err := NodeProfile(payload.Wallet)
	if err != nil {
		return profile, err
	}

	db().Model(&profile).Updates(
		models.NodeProfile{
			Name:     payload.Name,
			Bio:      payload.Bio,
			Email:    payload.Email,
			Location: payload.Location,
		},
	)

	return profile, nil
}

func NodeProfile(wallet string) (models.NodeProfile, error) {
	var profile models.NodeProfile

	if err := db().Model(&profile).Where("wallet = ?", wallet).First(&profile).Error; err != nil {
		return models.NodeProfile{}, errors.New("NodeProfile() profile not found for given wallet address")
	}

	return profile, nil
}

func PoolApplicationStatus(wallet string, accepted bool) {
	profile, _ := NodeProfile(wallet)
	profile.PoolAccepted = accepted
	db().Save(&profile)
}

func NodeApplicationStatus(wallet string, accepted bool) {
	profile, _ := NodeProfile(wallet)
	profile.NodeAccepted = accepted
	db().Save(&profile)
}
