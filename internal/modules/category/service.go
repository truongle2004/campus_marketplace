package category

import (
	"context"
	"fmt"
)

type CategoryService interface {
	ListTopLevel(ctx context.Context) ([]Response, error)
	ListByParentID(ctx context.Context, parentID string) ([]Response, error)
}

type categoryService struct {
	repo CategoryRepo
}

func NewService(repo CategoryRepo) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) ListTopLevel(ctx context.Context) ([]Response, error) {
	categories, err := s.repo.ListTopLevel(ctx)
	if err != nil {
		return nil, fmt.Errorf("CategoryService.ListTopLevel: %w", err)
	}
	return toResponses(categories), nil
}

func (s *categoryService) ListByParentID(ctx context.Context, parentID string) ([]Response, error) {
	id, err := parseUUID(parentID)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ParentExists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CategoryService.ListByParentID: %w", err)
	}
	if !exists {
		return nil, ErrNotFound
	}

	categories, err := s.repo.ListByParentID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CategoryService.ListByParentID: %w", err)
	}
	return toResponses(categories), nil
}
