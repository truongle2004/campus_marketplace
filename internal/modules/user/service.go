package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/truongle2004/campus_marketplace/internal/modules/campus"
	"github.com/truongle2004/campus_marketplace/pkg/auth"
)

type UserService interface {
	Bootstrap(ctx context.Context, authUser *auth.User, req BootstrapRequest) (*ProfileResponse, error)
	GetProfile(ctx context.Context, clerkUserID string) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, clerkUserID string, req UpdateProfileRequest) (*ProfileResponse, error)
	Deactivate(ctx context.Context, clerkUserID string) error
}

type userService struct {
	repo   UserRepo
	campus campus.CampusService
}

func NewService(repo UserRepo, campus campus.CampusService) UserService {
	return &userService{repo: repo, campus: campus}
}

func (s *userService) Bootstrap(ctx context.Context, authUser *auth.User, req BootstrapRequest) (*ProfileResponse, error) {
	if authUser.Email == "" {
		return nil, ErrEmailRequired
	}

	if _, err := s.repo.GetByClerkUserID(ctx, authUser.ClerkUserID); err == nil {
		return nil, ErrAlreadyExists
	} else if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("UserService.Bootstrap: %w", err)
	}

	campusResp, err := s.campus.GetByID(ctx, req.CampusID.String())
	if err != nil {
		if errors.Is(err, campus.ErrNotFound) {
			return nil, ErrInvalidCampus
		}
		return nil, fmt.Errorf("UserService.Bootstrap: %w", err)
	}

	user := &User{
		CampusID:    req.CampusID,
		ClerkUserID: authUser.ClerkUserID,
		Email:       authUser.Email,
		FullName:    req.FullName,
		Role:        RoleStudent,
		IsVerified:  emailMatchesCampusDomain(authUser.Email, campusResp.Domain),
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("UserService.Bootstrap: %w", err)
	}

	resp := toProfileResponse(*user)
	return &resp, nil
}

func (s *userService) GetProfile(ctx context.Context, clerkUserID string) (*ProfileResponse, error) {
	user, err := s.repo.GetByClerkUserID(ctx, clerkUserID)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetProfile: %w", err)
	}
	resp := toProfileResponse(*user)
	return &resp, nil
}

func (s *userService) UpdateProfile(ctx context.Context, clerkUserID string, req UpdateProfileRequest) (*ProfileResponse, error) {
	user, err := s.repo.GetByClerkUserID(ctx, clerkUserID)
	if err != nil {
		return nil, fmt.Errorf("UserService.UpdateProfile: %w", err)
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("UserService.UpdateProfile: %w", err)
	}

	resp := toProfileResponse(*user)
	return &resp, nil
}

func (s *userService) Deactivate(ctx context.Context, clerkUserID string) error {
	user, err := s.repo.GetByClerkUserID(ctx, clerkUserID)
	if err != nil {
		return fmt.Errorf("UserService.Deactivate: %w", err)
	}

	user.IsActive = false
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("UserService.Deactivate: %w", err)
	}
	return nil
}

func emailMatchesCampusDomain(email string, domain *string) bool {
	if domain == nil || *domain == "" {
		return false
	}
	suffix := "@" + strings.ToLower(strings.TrimSpace(*domain))
	return strings.HasSuffix(strings.ToLower(strings.TrimSpace(email)), suffix)
}
