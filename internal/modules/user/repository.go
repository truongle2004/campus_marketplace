package user

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepo interface {
	GetByClerkUserID(ctx context.Context, clerkUserID string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetByClerkUserID(ctx context.Context, clerkUserID string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).
		Where("clerk_user_id = ? AND is_active = ?", clerkUserID, true).
		First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("userRepo.GetByClerkUserID: %w", err)
	}
	return &u, nil
}

func (r *userRepo) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
