package auth

// import (
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestBearerToken(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name    string
// 		header  string
// 		want    string
// 		wantErr error
// 	}{
// 		{name: "valid bearer", header: "Bearer token123", want: "token123"},
// 		{name: "case insensitive scheme", header: "bearer token123", want: "token123"},
// 		{name: "missing header", header: "", wantErr: ErrMissingAuthHeader},
// 		{name: "invalid scheme", header: "Basic token123", wantErr: ErrInvalidAuthHeader},
// 		{name: "missing token", header: "Bearer ", wantErr: ErrInvalidAuthHeader},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			got, err := bearerToken(tt.header)
// 			if tt.wantErr != nil {
// 				require.ErrorIs(t, err, tt.wantErr)
// 				return
// 			}
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestPrimaryEmail(t *testing.T) {
// 	t.Parallel()

// 	email, ok := primaryEmail(jwt.MapClaims{
// 		"email_addresses": []any{
// 			map[string]any{"email_address": "student@hcmut.edu.vn"},
// 		},
// 	})
// 	assert.True(t, ok)
// 	assert.Equal(t, "student@hcmut.edu.vn", email)
// }

// func TestRequireAuth(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
// 	require.NoError(t, err)

// 	keyfunc := func(token *jwt.Token) (any, error) {
// 		return &privateKey.PublicKey, nil
// 	}

// 	issuer := "https://clerk.example.com"
// 	mw := newMiddleware(keyfunc, issuer, []string{"http://localhost:5173"})

// 	sign := func(claims jwt.MapClaims) string {
// 		signed, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
// 		require.NoError(t, err)
// 		return signed
// 	}

// 	now := time.Now()
// 	validClaims := jwt.MapClaims{
// 		"sub":   "user_abc123",
// 		"iss":   issuer,
// 		"azp":   "http://localhost:5173",
// 		"exp":   now.Add(time.Hour).Unix(),
// 		"iat":   now.Unix(),
// 		"email": "student@hcmut.edu.vn",
// 	}

// 	tests := []struct {
// 		name       string
// 		middleware *Middleware
// 		authHeader string
// 		wantStatus int
// 		wantUser   bool
// 	}{
// 		{
// 			name:       "accepts valid token",
// 			middleware: mw,
// 			authHeader: "Bearer " + sign(validClaims),
// 			wantStatus: http.StatusOK,
// 			wantUser:   true,
// 		},
// 		{
// 			name:       "rejects missing header",
// 			middleware: mw,
// 			wantStatus: http.StatusUnauthorized,
// 		},
// 		{
// 			name:       "rejects invalid token",
// 			middleware: mw,
// 			authHeader: "Bearer not-a-jwt",
// 			wantStatus: http.StatusUnauthorized,
// 		},
// 		{
// 			name:       "rejects wrong issuer",
// 			middleware: mw,
// 			authHeader: "Bearer " + sign(jwt.MapClaims{
// 				"sub": "user_abc123",
// 				"iss": "https://wrong.example.com",
// 				"azp": "http://localhost:5173",
// 				"exp": now.Add(time.Hour).Unix(),
// 			}),
// 			wantStatus: http.StatusUnauthorized,
// 		},
// 		{
// 			name:       "rejects wrong authorized party",
// 			middleware: mw,
// 			authHeader: "Bearer " + sign(jwt.MapClaims{
// 				"sub": "user_abc123",
// 				"iss": issuer,
// 				"azp": "http://evil.example.com",
// 				"exp": now.Add(time.Hour).Unix(),
// 			}),
// 			wantStatus: http.StatusUnauthorized,
// 		},
// 		{
// 			name:       "returns 503 when clerk is not configured",
// 			middleware: &Middleware{},
// 			authHeader: "Bearer " + sign(validClaims),
// 			wantStatus: http.StatusServiceUnavailable,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := gin.New()
// 			r.GET("/protected", tt.middleware.RequireAuth(), func(c *gin.Context) {
// 				user, ok := UserFromContext(c.Request.Context())
// 				if tt.wantUser {
// 					require.True(t, ok)
// 					assert.Equal(t, "user_abc123", user.ClerkUserID)
// 					assert.Equal(t, "student@hcmut.edu.vn", user.Email)
// 				}
// 				c.Status(http.StatusOK)
// 			})

// 			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
// 			if tt.authHeader != "" {
// 				req.Header.Set("Authorization", tt.authHeader)
// 			}
// 			rec := httptest.NewRecorder()
// 			r.ServeHTTP(rec, req)

// 			assert.Equal(t, tt.wantStatus, rec.Code)
// 		})
// 	}
// }
