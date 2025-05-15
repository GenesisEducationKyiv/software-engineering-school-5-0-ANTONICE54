package models

import "time"

type (
	Frequency string

	Subscription struct {
		Email     string `gorm:"primaryKey"`
		City      string
		Token     string `gorm:"unique"`
		Frequency Frequency
		Confirmed bool      `gorm:"default:false"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}
)

const (
	Daily  Frequency = "daily"
	Hourly Frequency = "hourly"
)
