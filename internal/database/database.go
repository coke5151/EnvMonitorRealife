package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DetectEnvironment struct {
	gorm.Model
	Temperature        float64
	HumidityPercentage float64
}

func New() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&DetectEnvironment{})
	return db
}
