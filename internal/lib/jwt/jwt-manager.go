package lib

import (
	"errors"
	"time"

	go_jwt "github.com/golang-jwt/jwt/v5"

	"main/internal/models"
	coreconfig "main/tools/pkg/core_config"
	tvoerrors "main/tools/pkg/tvo_errors"
)

var ErrCastClaims = errors.New("failed to cast token claims to UserClaims")

// JWTManager is a struct responsible for managing JWT tokens.
type JWTManager struct {
	cfg *coreconfig.JWT
}

// UserClaims represents the claims stored in JWT tokens for users.
type UserClaims struct {
	ID   int64 `json:"uid"`
	Role int64 `json:"role_id"`
	go_jwt.RegisteredClaims
}

// NewJWTManager initializes a new JWTManager with the provided configuration.
func NewJWTManager(cfg *coreconfig.JWT) *JWTManager {
	return &JWTManager{
		cfg: cfg,
	}
}

// GetTokenTTL returns the time-to-live (TTL) duration for authentication tokens.
func (manager *JWTManager) GetTokenTTL() time.Duration {
	return manager.cfg.AuthExpired
}

// GetRefreshTTL returns the time-to-live (TTL) duration for refresh tokens.
func (manager *JWTManager) GetRefreshTTL() time.Duration {
	return manager.cfg.RefreshExpired
}

// Generate creates new JWTManager token for given user and app.
func (manager *JWTManager) Generate(user *models.User) (string, error) {
	now := time.Now()

	token := go_jwt.NewWithClaims(go_jwt.SigningMethodHS512, UserClaims{
		ID:   user.ID,
		Role: int64(user.RoleID),
		RegisteredClaims: go_jwt.RegisteredClaims{
			ExpiresAt: go_jwt.NewNumericDate(now.Add(manager.cfg.AuthExpired)),
			IssuedAt:  go_jwt.NewNumericDate(now),
		},
	})

	tokenString, err := token.SignedString([]byte(manager.cfg.Secret))
	if err != nil {
		return "", tvoerrors.Wrap("error signing token", err)
	}

	return tokenString, nil
}

// Verify parses the provided JWTManager token using the given secret and returns the corresponding UserClaims struct.
func (manager *JWTManager) Verify(jwtToken string) (*UserClaims, error) {
	token, err := go_jwt.ParseWithClaims(
		jwtToken,
		&UserClaims{},
		func(token *go_jwt.Token) (interface{}, error) {
			return []byte(manager.cfg.Secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, tvoerrors.Wrap("invalid token", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, tvoerrors.Wrap("invalid token", ErrCastClaims)
	}

	return claims, nil
}
