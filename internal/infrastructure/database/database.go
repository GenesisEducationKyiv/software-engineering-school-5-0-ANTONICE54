package database

import (
	"fmt"
	"log"
	"weather-forecast/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(DBHost, DBUser, DBPassword, DBName, DBPort string) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DBHost,
		DBUser,
		DBPassword,
		DBName,
		DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to establish connection with database: %s", err.Error())
	}

	return db

}

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(&models.Subscription{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %s", err.Error())
	}
}
