package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	coreconfig "main/tools/pkg/core_config"
	tvoerrors "main/tools/pkg/tvo_errors"
)

// NewClient initializes the database connection pool.
func NewClient(ctx context.Context, cfg *coreconfig.Database) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN to config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, tvoerrors.Wrap("unable to create connection pool: %w", err)
	}

	return pool, nil
}
