package campus

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CampusRepo interface {
	Create(ctx context.Context, campus *Campus) error
	GetByID(ctx context.Context, id uuid.UUID) (*Campus, error)
	ListActive(ctx context.Context) ([]Campus, error)
}

type campusRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) CampusRepo {
	return &campusRepo{db: db}
}

func (r *campusRepo) Create(ctx context.Context, campus *Campus) error {
	return r.db.WithContext(ctx).Create(campus).Error
}

func (r *campusRepo) GetByID(ctx context.Context, id uuid.UUID) (*Campus, error) {
	var c Campus
	err := r.db.WithContext(ctx).First(&c, "id = ? AND is_active = ?", id, true).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("campusRepo.GetByID: %w", err)
	}
	return &c, nil
}

func (r *campusRepo) ListActive(ctx context.Context) ([]Campus, error) {
	var campuses []Campus
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&campuses).Error
	if err != nil {
		return nil, fmt.Errorf("campusRepo.ListActive: %w", err)
	}
	return campuses, nil
}
