package database

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/truongle2004/campus_marketplace/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		env.GetEnv("DB_HOST", "localhost"),
		env.GetEnv("DB_PORT", "5432"),
		env.GetEnv("DB_USERNAME", "postgres"),
		env.GetEnv("DB_PASSWORD", "postgres"),
		env.GetEnv("DB_DATABASE", "campus_marketplace"),
		env.GetEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
