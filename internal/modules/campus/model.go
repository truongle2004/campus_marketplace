package campus

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Campus struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Slug      string    `gorm:"size:100;not null;uniqueIndex"`
	Domain    *string   `gorm:"size:100;uniqueIndex"`
	Country   string    `gorm:"size:100;not null"`
	City      string    `gorm:"size:100;not null"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (Campus) TableName() string {
	return "campuses"
}

func (c *Campus) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
