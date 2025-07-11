package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"main/internal/models"
	tvoerrors "main/tools/pkg/tvo_errors"
)

// UserTokenRepository handles user token-related operations in PostgreSQL.
type UserTokenRepository struct {
	db    *pgxpool.Pool
	debug bool
}

// NewUserTokenRepository creates a new instance of UserTokenRepository with the given PostgreSQL connection pool.
func NewUserTokenRepository(db *pgxpool.Pool, debug bool) *UserTokenRepository {
	return &UserTokenRepository{
		db:    db,
		debug: debug,
	}
}

// Create saves a new user token to the database.
func (utr *UserTokenRepository) Create(ctx context.Context, id int64, accessToken, refreshToken string, tokenTTL,
	refreshTokenTTL time.Duration) (*models.UserToken, error) {
	const op = "postgresql.UserTokenRepository.Create"

	expiredAt := time.Now().UTC().Add(tokenTTL)
	refreshExpiredAt := time.Now().UTC().Add(refreshTokenTTL)

	query := `INSERT INTO user_tokens (user_id, token, refresh_token, expired_at, refresh_expired_at) 
		VALUES ($1, $2, $3, $4, $5);`

	res, err := utr.db.Exec(ctx, query, id, accessToken, refreshToken, expiredAt, refreshExpiredAt)
	if err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}
	if res.RowsAffected() != 1 {
		return nil, tvoerrors.Wrap(op, tvoerrors.ErrInsertFailed)
	}

	userToken := &models.UserToken{
		UserID:           id,
		Token:            accessToken,
		ExpiredAt:        expiredAt,
		RefreshToken:     refreshToken,
		RefreshExpiredAt: refreshExpiredAt,
	}

	return userToken, nil
}

// GetRefreshToken retrieves a user token from the database based on the access token and refresh token.
func (utr *UserTokenRepository) GetRefreshToken(ctx context.Context, refresh string) (*models.UserToken, error) {
	const op = "postgresql.UserTokenRepository.GetRefreshToken"
	var userToken models.UserToken
	now := time.Now().UTC()

	query := "SELECT user_id FROM user_tokens WHERE refresh_token = $1 AND refresh_expired_at >= $2;"
	if err := utr.db.QueryRow(ctx, query, refresh, now).Scan(&userToken.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}
		return nil, tvoerrors.Wrap(op, err)
	}

	return &userToken, nil
}

// DeleteRefreshToken resets the token associated with a user in the database.
func (utr *UserTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	const op = "postgresql.UserTokenRepository.DeleteRefreshToken"

	query := "UPDATE user_tokens SET refresh_token = NULL WHERE token = $1;"
	res, err := utr.db.Exec(ctx, query, token)
	if err != nil {
		return tvoerrors.Wrap(op, err)
	}

	if res.RowsAffected() != 1 {
		return tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	return nil
}

// TokenReset resets the refresh token for the user identified by the provided userID.
func (utr *UserTokenRepository) TokenReset(ctx context.Context, userID int64) error {
	const op = "postgresql.UserTokenRepository.TokenReset"

	query := "UPDATE user_tokens SET refresh_token = NULL WHERE user_id = $1;"
	// ignore res, the user may have never logged in.
	if _, err := utr.db.Exec(ctx, query, userID); err != nil {
		return tvoerrors.Wrap(op, err)
	}

	return nil
}

// ActiveTokens retrieves the active user tokens associated with the provided userID.
func (utr *UserTokenRepository) ActiveTokens(ctx context.Context, userID int64) ([]models.UserToken, error) {
	const op = "postgresql.UserTokenRepository.ActiveTokens"

	now := time.Now().UTC()

	query := "SELECT token, expired_at FROM user_tokens WHERE user_id = $1 AND expired_at > $2 AND refresh_token IS NOT NULL;"
	rows, err := utr.db.Query(ctx, query, userID, now)
	if err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}
	defer rows.Close()

	var tokens []models.UserToken
	for rows.Next() {
		token := models.UserToken{}
		if err = rows.Scan(&token.Token, &token.ExpiredAt); err != nil {
			return nil, tvoerrors.Wrap(op, err)
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// GetUserIdByToken получает информацию о пользователе по токену
func (utr *UserTokenRepository) GetUserIdByToken(ctx context.Context, token string) (*models.UserToken, error) {
	const op = "postgresql.UserTokenRepository.GetUserIdByToken"

	now := time.Now().UTC()

	var t models.UserToken
	query := `SELECT user_id, expired_at FROM user_tokens WHERE token = $1 AND expired_at > $2  AND refresh_token IS NOT NULL`
	if err := utr.db.QueryRow(ctx, query, token, now).Scan(&t.UserID, &t.ExpiredAt); err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}

	return &t, nil
}
