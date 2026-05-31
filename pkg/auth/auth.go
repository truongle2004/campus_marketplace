package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/truongle2004/campus_marketplace/pkg/env"
	"github.com/truongle2004/campus_marketplace/pkg/respond"
)

var (
	ErrNotConfigured     = errors.New("clerk auth is not configured")
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
	ErrInvalidToken      = errors.New("invalid token")
)

var clerkAlgorithms = []string{"RS256"}

type ctxKey string

const userKey ctxKey = "authUser"

// User holds identity extracted from a verified Clerk session JWT.
type User struct {
	ClerkUserID string
	Email       string
}

// Middleware verifies Clerk session tokens via JWKS.
type Middleware struct {
	jwks              keyfunc.Keyfunc
	issuer            string
	authorizedParties map[string]struct{}
}

// NewMiddleware loads Clerk JWKS from CLERK_JWKS_URL. Returns a no-op verifier
// when the URL is unset (protected routes respond with 503).
func NewMiddleware() (*Middleware, error) {
	jwksURL := env.GetEnv("CLERK_JWKS_URL", "")
	if jwksURL == "" {
		return &Middleware{}, nil
	}

	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("auth.NewMiddleware: %w", err)
	}

	return newMiddleware(
		jwks,
		env.GetEnv("CLERK_ISSUER", ""),
		parseCSV(env.GetEnv("CLERK_AUTHORIZED_PARTIES", "")),
	), nil
}

func newMiddleware(jwks keyfunc.Keyfunc, issuer string, authorizedParties []string) *Middleware {
	parties := make(map[string]struct{}, len(authorizedParties))
	for _, p := range authorizedParties {
		parties[p] = struct{}{}
	}
	return &Middleware{
		jwks:              jwks,
		issuer:            issuer,
		authorizedParties: parties,
	}
}

// RequireAuth rejects requests without a valid Clerk Bearer token.
func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.jwks == nil {
			respond.Error(c, 503, "authentication is not configured")
			c.Abort()
			return
		}

		tokenStr, err := bearerToken(c.GetHeader("Authorization"))
		if err != nil {
			respond.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		user, err := m.parseToken(tokenStr)
		if err != nil {
			respond.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), userKey, user))
		c.Next()
	}
}

func (m *Middleware) parseToken(tokenStr string) (*User, error) {
	token, err := jwt.Parse(tokenStr, m.jwks.Keyfunc, jwt.WithValidMethods(clerkAlgorithms))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if m.issuer != "" {
		iss, _ := claims["iss"].(string)
		if iss != m.issuer {
			return nil, fmt.Errorf("%w: invalid issuer", ErrInvalidToken)
		}
	}

	if len(m.authorizedParties) > 0 {
		azp, _ := claims["azp"].(string)
		if _, ok := m.authorizedParties[azp]; !ok {
			return nil, fmt.Errorf("%w: invalid authorized party", ErrInvalidToken)
		}
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, fmt.Errorf("%w: missing subject", ErrInvalidToken)
	}

	email, _ := claims["email"].(string)
	if email == "" {
		email, _ = primaryEmail(claims)
	}

	return &User{ClerkUserID: sub, Email: email}, nil
}

func primaryEmail(claims jwt.MapClaims) (string, bool) {
	emails, ok := claims["email_addresses"].([]any)
	if !ok || len(emails) == 0 {
		return "", false
	}
	first, ok := emails[0].(map[string]any)
	if !ok {
		return "", false
	}
	email, _ := first["email_address"].(string)
	return email, email != ""
}

func bearerToken(header string) (string, error) {
	if header == "" {
		return "", ErrMissingAuthHeader
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", ErrInvalidAuthHeader
	}
	return parts[1], nil
}

func parseCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// UserFromContext returns the authenticated Clerk user set by RequireAuth.
func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey).(*User)
	return user, ok
}
