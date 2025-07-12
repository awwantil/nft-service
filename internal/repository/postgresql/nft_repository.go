package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"main/internal/dto"
	"main/internal/models"
	tvoerrors "main/tools/pkg/tvo_errors"
)

// NftDataRepository handles nft-related operations in PostgreSQL.
type NftDataRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new instance of NftDataRepository with the given PostgreSQL connection pool.
func NewNftDataRepository(db *pgxpool.Pool) *NftDataRepository {
	return &NftDataRepository{
		db: db,
	}
}

// CreateNftData saves a new nft data
func (ur *NftDataRepository) CreateNftData(ctx context.Context, data *dto.NftData) error {
	const op = "postgresql.NftDataRepository.CreateNftData"
	var nft models.NftDataModel

	query := "INSERT INTO nft_data (token_id, content, cidv0, cidv1, file_size, file_name) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	if err := ur.db.QueryRow(ctx, query, data.TokenId, data.Description, data.CidV0, data.CidV1, data.FileSize, data.FileName).
		Scan(&nft.ID); err != nil {
		return tvoerrors.Wrap(op, err)
	}
	return nil
}
