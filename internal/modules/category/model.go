package category

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index"`
	Parent    *Category  `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name      string     `gorm:"size:100;not null"`
	Slug      string     `gorm:"size:100;not null;uniqueIndex"`
	IconURL   *string    `gorm:"type:text"`
	SortOrder int        `gorm:"not null;default:0"`
	IsActive  bool       `gorm:"not null;default:true"`
}

func (Category) TableName() string {
	return "categories"
}

func (c *Category) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
