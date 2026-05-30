package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepo interface {
	ListTopLevel(ctx context.Context) ([]Category, error)
	ListByParentID(ctx context.Context, parentID uuid.UUID) ([]Category, error)
	ParentExists(ctx context.Context, parentID uuid.UUID) (bool, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) ListTopLevel(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND parent_id IS NULL", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("categoryRepo.ListTopLevel: %w", err)
	}
	return categories, nil
}

func (r *categoryRepo) ListByParentID(ctx context.Context, parentID uuid.UUID) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND parent_id = ?", true, parentID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("categoryRepo.ListByParentID: %w", err)
	}
	return categories, nil
}

func (r *categoryRepo) ParentExists(ctx context.Context, parentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Category{}).
		Where("id = ? AND is_active = ?", parentID, true).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("categoryRepo.ParentExists: %w", err)
	}
	return count > 0, nil
}

func parseUUID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid id")
	}
	return id, nil
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
