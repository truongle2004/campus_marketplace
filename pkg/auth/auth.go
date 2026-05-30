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

var ErrNotConfigured = errors.New("clerk auth is not configured")

type ctxKey string

const userKey ctxKey = "authUser"

type User struct {
	ClerkUserID string
	Email       string
}

type Middleware struct {
	jwks   keyfunc.Keyfunc
	issuer string
}

func NewMiddleware() (*Middleware, error) {
	jwksURL := env.GetEnv("CLERK_JWKS_URL", "")
	if jwksURL == "" {
		return &Middleware{}, nil
	}

	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("auth.NewMiddleware: %w", err)
	}

	return &Middleware{
		jwks:   jwks,
		issuer: env.GetEnv("CLERK_ISSUER", ""),
	}, nil
}

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
	token, err := jwt.Parse(tokenStr, m.jwks.Keyfunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if m.issuer != "" {
		iss, _ := claims["iss"].(string)
		if iss != m.issuer {
			return nil, errors.New("invalid issuer")
		}
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, errors.New("missing subject")
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
		return "", errors.New("missing authorization header")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}

func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey).(*User)
	return user, ok
}
