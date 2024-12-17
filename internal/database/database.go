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

type FanStatus struct {
	gorm.Model
	Status string // "off", "clockwise", "counterclockwise"
}

func New() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&DetectEnvironment{}, &FanStatus{})
	return db
}

// GetLatestFanStatus 取得最新的風扇狀態
func GetLatestFanStatus(db *gorm.DB) (*FanStatus, error) {
	var status FanStatus
	result := db.Last(&status)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果找不到記錄，返回預設狀態
			return &FanStatus{
				Status: "off", // 預設為關閉
			}, nil
		}
		return nil, result.Error
	}
	return &status, nil
}

// InsertFanStatus 插入新的風扇狀態
func InsertFanStatus(db *gorm.DB, status string) error {
	return db.Create(&FanStatus{
		Status: status,
	}).Error
}
