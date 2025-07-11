package repository

import (
	"context"
	"time"

	"main/internal/models"
)

// RoleRepository provides methods for managing roles.
type RoleRepository interface {
	RoleByName(ctx context.Context, roleName string) (*models.Role, error)
	RoleById(ctx context.Context, roleId int64) (*models.Role, error)
}

// UserTokenRepository provides methods for managing user tokens.
type UserTokenRepository interface {
	Create(ctx context.Context, id int64, accessToken, refreshToken string, tokenTTL, refreshTokenTTL time.Duration) (*models.UserToken, error)
	GetRefreshToken(ctx context.Context, refresh string) (*models.UserToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	TokenReset(ctx context.Context, userID int64) error
	ActiveTokens(ctx context.Context, userID int64) ([]models.UserToken, error)
	GetUserIdByToken(ctx context.Context, token string) (*models.UserToken, error)
}

// UserRepository provides methods for managing user-related operations.
type UserRepository interface {
	PhoneExists(ctx context.Context, phoneNumber string) (bool, error)
	UserByPhone(ctx context.Context, phone string) (*models.User, error)
	UserById(ctx context.Context, id int64) (*models.User, error)
	UpdateTelegramId(ctx context.Context, id, telegramId int64) error
	CreateUser(ctx context.Context, phone, password string) (*models.User, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
	UpdateLastVisit(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	DigUpUser(ctx context.Context, id int64) (*models.User, error)
	UpdatePhone(ctx context.Context, phone string, id int64) error
	ChangeRole(ctx context.Context, id, roleId int64) error
}
