package campus

import (
	"context"
	"fmt"
)

type CampusService interface {
	ListActive(ctx context.Context) ([]Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
}

type campusService struct {
	repo CampusRepo
}

func NewService(repo CampusRepo) CampusService {
	return &campusService{repo: repo}
}

func (s *campusService) ListActive(ctx context.Context) ([]Response, error) {
	campuses, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("CampusService.ListActive: %w", err)
	}
	return toResponses(campuses), nil
}

func (s *campusService) GetByID(ctx context.Context, id string) (*Response, error) {
	campusID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	campus, err := s.repo.GetByID(ctx, campusID)
	if err != nil {
		return nil, err
	}

	resp := toResponse(*campus)
	return &resp, nil
}
