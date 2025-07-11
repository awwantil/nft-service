package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"main/internal/auth/tools"
	"main/internal/models"
	tvoerrors "main/tools/pkg/tvo_errors"
	tvomodels "main/tools/pkg/tvo_models"
)

// UserRepository handles user-related operations in PostgreSQL.
type UserRepository struct {
	db     *pgxpool.Pool
	secret string
}

// NewUserRepository creates a new instance of UserRepository with the given PostgreSQL connection pool.
func NewUserRepository(db *pgxpool.Pool, secret string) *UserRepository {
	return &UserRepository{
		db:     db,
		secret: secret,
	}
}

// PhoneExists checks if a user with the given phone number exists in the database.
func (ur *UserRepository) PhoneExists(ctx context.Context, phoneNumber string) (bool, error) {
	const op = "postgresql.UserRepository.PhoneExists"
	var id int64

	query := "SELECT id FROM users WHERE phone = $1 AND deleted_at IS NULL;"

	if err := ur.db.QueryRow(ctx, query, phoneNumber).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}
		return false, tvoerrors.Wrap(op, err)
	}
	return true, nil
}

// UserByPhone retrieves a user from the database by their phone number.
func (ur *UserRepository) UserByPhone(ctx context.Context, phone string) (*models.User, error) {
	const op = "postgresql.UserRepository.UserByPhone"
	var user models.User

	query := `SELECT id, phone, password, salt, role_id FROM users WHERE phone = $1 AND deleted_at IS NULL;`

	if err := ur.db.QueryRow(ctx, query, phone).Scan(&user.ID, &user.Phone,
		&user.Password, &user.Salt, &user.RoleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}
		return nil, tvoerrors.Wrap(op, err)
	}

	return &user, nil
}

// UserById retrieves a user from the database by their ID.
func (ur *UserRepository) UserById(ctx context.Context, id int64) (*models.User, error) {
	const op = "postgresql.UserRepository.UserById"
	var user models.User

	query := `SELECT id, phone, password, salt, role_id FROM users WHERE id = $1 AND deleted_at IS NULL;`

	if err := ur.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Phone,
		&user.Password, &user.Salt, &user.RoleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}
		return nil, tvoerrors.Wrap(op, err)
	}

	return &user, nil
}

// CreateUser saves a new user to the database with the provided phone number, password, and salt.
func (ur *UserRepository) CreateUser(ctx context.Context, phone, password string) (*models.User, error) {
	const op = "postgresql.UserRepository.CreateUser"
	var user models.User

	salt := tools.GenerateRandomSalt(tools.SaltSize)
	hashPassword := tools.HashPassword(password, ur.secret, salt)

	query := "INSERT INTO users (phone, password, salt, role_id, telegram_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, role_id, telegram_id"
	if err := ur.db.QueryRow(ctx, query, phone, hashPassword, salt, tvomodels.USER).
		Scan(&user.ID, &user.RoleID); err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}

	return &user, nil
}

// UpdatePassword updates the password and salt of a user in the database.
func (ur *UserRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	const op = "postgresql.UserRepository.UpdatePassword"

	now := time.Now().UTC()
	salt := tools.GenerateRandomSalt(tools.SaltSize)
	hashPassword := tools.HashPassword(password, ur.secret, salt)

	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return tvoerrors.Wrap(op, err)
	}

	query := "SELECT id FROM users WHERE id = $1 AND deleted_at IS NULL FOR UPDATE;"

	if err = tx.QueryRow(ctx, query, id).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return tvoerrors.Wrap(op, rollbackErr)
			}
			return tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}

		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	query = "UPDATE users SET password = $1, salt = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL;"

	result, err := tx.Exec(ctx, query, hashPassword, salt, now, id)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	if result.RowsAffected() != 1 {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	query = "UPDATE user_tokens SET refresh_token = NULL WHERE user_id = $1;"

	// ignore res, the user may have never logged in.
	if _, err = tx.Exec(ctx, query, id); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return tvoerrors.Wrap(op, err)
	}

	return nil
}

// UpdateLastVisit updates the last visit timestamp of a user in the database.
func (ur *UserRepository) UpdateLastVisit(ctx context.Context, id int64) error {
	const op = "postgresql.UserRepository.UpdateLastVisit"

	now := time.Now().UTC()
	query := "UPDATE users SET last_visited_at = $1 WHERE id = $2;"

	result, err := ur.db.Exec(ctx, query, now, id)
	if err != nil {
		return tvoerrors.Wrap(op, err)
	}

	if result.RowsAffected() != 1 {
		return tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	return nil
}

func (ur *UserRepository) UpdateTelegramId(ctx context.Context, id, telegramId int64) error {
	const op = "postgresql.UserRepository.UpdateTelegramId"

	now := time.Now().UTC()
	query := "UPDATE users SET updated_at = $1, telegram_id = $2 WHERE id = $3;"

	result, err := ur.db.Exec(ctx, query, now, telegramId, id)
	if err != nil {
		return tvoerrors.Wrap(op, err)
	}

	if result.RowsAffected() != 1 {
		return tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	return nil
}

// DeleteUser deletes a user record from the database with the given ID.
func (ur *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	const op = "postgresql.UserRepository.DeleteUser"
	var userId int64
	now := time.Now().UTC()

	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return tvoerrors.Wrap(op, err)
	}

	query := "SELECT id FROM users WHERE id = $1 AND deleted_at IS NULL FOR UPDATE;"

	if err = tx.QueryRow(ctx, query, id).Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return tvoerrors.Wrap(op, rollbackErr)
			}
			return tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}

		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	query = "UPDATE users SET deleted_at = $1, updated_at = $2 WHERE id = $3;"
	result, err := tx.Exec(ctx, query, now, now, id)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	if result.RowsAffected() != 1 {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	query = "UPDATE user_tokens SET refresh_token = NULL WHERE user_id = $1;"

	// ignore res, the user may have never logged in.
	if _, err = tx.Exec(ctx, query, id); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return tvoerrors.Wrap(op, err)
	}

	return nil
}

// DigUpUser restores a previously deleted user record identified by the given ID.
func (ur *UserRepository) DigUpUser(ctx context.Context, id int64) (*models.User, error) {
	const op = "postgresql.UserRepository.DigUpUser"

	var exists bool
	var user models.User
	now := time.Now().UTC()

	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	query := "SELECT id, phone, role_id FROM users WHERE id = $1 AND deleted_at IS NOT NULL FOR UPDATE;"

	if err = tx.QueryRow(ctx, query, id).Scan(&user.ID, &user.Phone, &user.RoleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}
		return nil, tvoerrors.Wrap(op, err)
	}

	// check phone exists

	query = "SELECT EXISTS (select id FROM users WHERE (phone = $1) AND deleted_at IS NULL);"
	if err = tx.QueryRow(ctx, query, user.Phone).Scan(&exists); err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}

	if exists {
		return nil, tvoerrors.Wrap(op, tvoerrors.Wrap("Phone already exists", tvoerrors.ErrUpdateFailed))
	}

	query = "UPDATE users SET deleted_at = NULL, updated_at = $1 WHERE id = $2;"

	result, err := tx.Exec(ctx, query, now, id)
	if err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}

	if result.RowsAffected() != 1 {
		return nil, tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}

	return &user, nil
}

// UpdatePhone updates the phone number of a user with the given ID.
func (ur *UserRepository) UpdatePhone(ctx context.Context, phone string, id int64) error {
	const op = "postgresql.UserRepository.UpdatePhone"
	now := time.Now().UTC()

	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return tvoerrors.Wrap(op, err)
	}

	query := "SELECT id FROM users WHERE id = $1 AND deleted_at IS NULL FOR UPDATE;"

	if err = tx.QueryRow(ctx, query, id).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return tvoerrors.Wrap(op, rollbackErr)
			}
			return tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}

		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	query = "UPDATE users SET phone = $1, updated_at = $2 WHERE id = $3;"
	result, err := tx.Exec(ctx, query, phone, now, id)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	if result.RowsAffected() != 1 {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	if err = tx.Commit(ctx); err != nil {
		return tvoerrors.Wrap(op, err)
	}

	return nil
}

// ChangeRole updates the role of a user in the database
func (ur *UserRepository) ChangeRole(ctx context.Context, id, roleId int64) error {
	const op = "postgresql.UserRepository.ChangeRole"
	now := time.Now().UTC()

	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return tvoerrors.Wrap(op, err)
	}

	query := "SELECT id FROM users WHERE id = $1 AND deleted_at IS NULL FOR UPDATE;"

	if err = tx.QueryRow(ctx, query, id).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return tvoerrors.Wrap(op, rollbackErr)
			}
			return tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}

		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	query = "UPDATE users SET role_id = $1, updated_at = $2 WHERE id = $3;"

	result, err := tx.Exec(ctx, query, roleId, now, id)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	if result.RowsAffected() != 1 {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, tvoerrors.ErrUpdateFailed)
	}

	query = "UPDATE user_tokens SET refresh_token = NULL WHERE user_id = $1;"

	// ignore res, the user may have never logged in.
	if _, err = tx.Exec(ctx, query, id); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return tvoerrors.Wrap(op, rollbackErr)
		}
		return tvoerrors.Wrap(op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return tvoerrors.Wrap(op, err)
	}

	return nil
}
