package user

import (
	"time"

	"github.com/google/uuid"
)

type ProfileResponse struct {
	ID         uuid.UUID `json:"id"`
	CampusID   uuid.UUID `json:"campus_id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
	Role       Role      `json:"role"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BootstrapRequest struct {
	CampusID uuid.UUID `json:"campus_id" binding:"required"`
	FullName string    `json:"full_name" binding:"required,min=1,max=255"`
}

type UpdateProfileRequest struct {
	FullName  *string `json:"full_name" binding:"omitempty,min=1,max=255"`
	AvatarURL *string `json:"avatar_url" binding:"omitempty,max=2048"`
}

func toProfileResponse(u User) ProfileResponse {
	return ProfileResponse{
		ID:         u.ID,
		CampusID:   u.CampusID,
		Email:      u.Email,
		FullName:   u.FullName,
		AvatarURL:  u.AvatarURL,
		Role:       u.Role,
		IsVerified: u.IsVerified,
		IsActive:   u.IsActive,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
