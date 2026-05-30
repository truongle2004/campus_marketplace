package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/truongle2004/campus_marketplace/internal/modules/campus"
	"gorm.io/gorm"
)

type Role string

const (
	RoleStudent Role = "student"
	RoleStaff   Role = "staff"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CampusID    uuid.UUID     `gorm:"type:uuid;not null;index"`
	Campus      campus.Campus `gorm:"foreignKey:CampusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ClerkUserID string    `gorm:"size:128;not null;uniqueIndex"`
	Email       string    `gorm:"size:255;not null;uniqueIndex"`
	FullName    string    `gorm:"size:255;not null"`
	AvatarURL   *string   `gorm:"type:text"`
	Role        Role      `gorm:"size:20;not null;default:student"`
	IsVerified  bool      `gorm:"not null;default:false"`
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
