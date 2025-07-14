package database

import (
	"fmt"
	"log"
	"subscription-service/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to establish connection with database: %s", err.Error())
	}

	return db

}

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(&Subscription{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %s", err.Error())
	}
}
