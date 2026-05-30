package campus

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCampusRepo struct {
	listActiveFn func(ctx context.Context) ([]Campus, error)
	getByIDFn    func(ctx context.Context, id uuid.UUID) (*Campus, error)
}

func (m *mockCampusRepo) Create(_ context.Context, _ *Campus) error {
	return nil
}

func (m *mockCampusRepo) ListActive(ctx context.Context) ([]Campus, error) {
	if m.listActiveFn != nil {
		return m.listActiveFn(ctx)
	}
	return nil, nil
}

func (m *mockCampusRepo) GetByID(ctx context.Context, id uuid.UUID) (*Campus, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, ErrNotFound
}

func TestService_ListActive(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	domain := "hcmut.edu.vn"
	id := uuid.MustParse("00000000-0000-4000-8000-000000000001")

	tests := []struct {
		name    string
		repo    *mockCampusRepo
		want    []Response
		wantErr string
	}{
		{
			name: "returns active campuses mapped to responses",
			repo: &mockCampusRepo{
				listActiveFn: func(_ context.Context) ([]Campus, error) {
					return []Campus{
						{
							ID:        id,
							Name:      "HCMUT",
							Slug:      "hcmut",
							Domain:    &domain,
							Country:   "Vietnam",
							City:      "Ho Chi Minh City",
							IsActive:  true,
							CreatedAt: createdAt,
						},
					}, nil
				},
			},
			want: []Response{
				{
					ID:        id,
					Name:      "HCMUT",
					Slug:      "hcmut",
					Domain:    &domain,
					Country:   "Vietnam",
					City:      "Ho Chi Minh City",
					CreatedAt: createdAt,
				},
			},
		},
		{
			name: "returns empty slice when no campuses",
			repo: &mockCampusRepo{
				listActiveFn: func(_ context.Context) ([]Campus, error) {
					return []Campus{}, nil
				},
			},
			want: []Response{},
		},
		{
			name: "wraps repository errors",
			repo: &mockCampusRepo{
				listActiveFn: func(_ context.Context) ([]Campus, error) {
					return nil, errors.New("db unavailable")
				},
			},
			wantErr: "CampusService.ListActive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewService(tt.repo)
			got, err := svc.ListActive(context.Background())

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	domain := "hcmut.edu.vn"
	id := uuid.MustParse("00000000-0000-4000-8000-000000000001")

	tests := []struct {
		name    string
		rawID   string
		repo    *mockCampusRepo
		want    *Response
		wantErr error
		errMsg  string
	}{
		{
			name:  "returns campus by id",
			rawID: id.String(),
			repo: &mockCampusRepo{
				getByIDFn: func(_ context.Context, campusID uuid.UUID) (*Campus, error) {
					assert.Equal(t, id, campusID)
					return &Campus{
						ID:        id,
						Name:      "HCMUT",
						Slug:      "hcmut",
						Domain:    &domain,
						Country:   "Vietnam",
						City:      "Ho Chi Minh City",
						IsActive:  true,
						CreatedAt: createdAt,
					}, nil
				},
			},
			want: &Response{
				ID:        id,
				Name:      "HCMUT",
				Slug:      "hcmut",
				Domain:    &domain,
				Country:   "Vietnam",
				City:      "Ho Chi Minh City",
				CreatedAt: createdAt,
			},
		},
		{
			name:    "rejects invalid uuid",
			rawID:   "not-a-uuid",
			repo:    &mockCampusRepo{},
			wantErr: nil,
			errMsg:  "invalid id",
		},
		{
			name:  "returns not found from repository",
			rawID: id.String(),
			repo: &mockCampusRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*Campus, error) {
					return nil, ErrNotFound
				},
			},
			wantErr: ErrNotFound,
		},
		{
			name:  "passes through repository errors",
			rawID: id.String(),
			repo: &mockCampusRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*Campus, error) {
					return nil, errors.New("db unavailable")
				},
			},
			wantErr: nil,
			errMsg:  "db unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewService(tt.repo)
			got, err := svc.GetByID(context.Background(), tt.rawID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}

			if tt.errMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
