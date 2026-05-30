package database

import (
	"fmt"

	"github.com/truongle2004/campus_marketplace/internal/modules/campus"
	"github.com/truongle2004/campus_marketplace/internal/modules/category"
	"github.com/truongle2004/campus_marketplace/internal/modules/user"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&campus.Campus{},
		&user.User{},
		&category.Category{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
