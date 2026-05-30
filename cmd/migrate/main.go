package main

import (
	"log"

	"github.com/truongle2004/campus_marketplace/internal/database"
)

func main() {
	db := database.NewDatabase()
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migration completed")
}
