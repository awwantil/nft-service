package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"main/internal/models"
	tvoerrors "main/tools/pkg/tvo_errors"
)

// RoleRepository handles database operations related to roles.
type RoleRepository struct {
	db *pgxpool.Pool
}

// NewRoleRepository creates a new RoleRepository instance.
func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

// RoleByName retrieves a role from the database by name.
func (rr *RoleRepository) RoleByName(ctx context.Context, roleName string) (*models.Role, error) {
	const op = "postgresql.RoleRepository.RoleByName"
	var role models.Role
	query := "SELECT id, name FROM roles WHERE name = $1;"

	if err := rr.db.QueryRow(ctx, query, roleName).Scan(&role.ID, &role.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}
		return nil, tvoerrors.Wrap(op, err)
	}

	return &role, nil
}

// RoleById retrieves a role from the database by ID.
func (rr *RoleRepository) RoleById(ctx context.Context, roleId int64) (*models.Role, error) {
	const op = "postgresql.RoleRepository.RoleById"
	var role models.Role

	query := "SELECT id, name FROM roles WHERE id = $1;"
	if err := rr.db.QueryRow(ctx, query, roleId).Scan(&role.ID, &role.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tvoerrors.Wrap(op, tvoerrors.ErrNotFound)
		}
		return nil, tvoerrors.Wrap(op, err)
	}

	return &role, nil
}
