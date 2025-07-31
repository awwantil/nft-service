package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
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

// ReadNftData takes one nft data
func (ur *NftDataRepository) ReadNftData(ctx context.Context, tokenId int64) (models.NftDataModel, error) {
	const op = "postgresql.NftDataRepository.ReadNftData"
	var nft models.NftDataModel
	query := "SELECT token_id, content, cidv0, cidv1  FROM nft_data where token_id = $1 LIMIT 1;"

	if err := ur.db.QueryRow(ctx, query, tokenId).Scan(
		&nft.TokenId, &nft.Description, &nft.CidV0, &nft.CidV1); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nft, tvoerrors.Wrap("postgresql.NftDataRepository.ReadNftData", err)
		}
	}

	return nft, nil
}

// ReadAllNftData takes all nft data
func (ur *NftDataRepository) ReadAllNftData(ctx context.Context, limit int) ([]models.NftDataModel, error) {
	const op = "postgresql.NftDataRepository.ReadNftData"
	query := "SELECT token_id, content, cidv0, cidv1  FROM nft_data LIMIT $1;"

	rows, err := ur.db.Query(ctx, query, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.NftDataModel{}, nil
		}
		return nil, tvoerrors.Wrap(op, err)
	}
	defer rows.Close()

	var nfts []models.NftDataModel
	for rows.Next() {
		var nft models.NftDataModel
		if err := rows.Scan(&nft.TokenId, &nft.Description, &nft.CidV0, &nft.CidV1); err != nil {
			return nil, tvoerrors.Wrap(op, err)
		}
		nfts = append(nfts, nft)
	}

	if err = rows.Err(); err != nil {
		return nil, tvoerrors.Wrap(op, err)
	}

	return nfts, nil
}

// TokenIdExists checks if a nft data exists by its token iD.
func (ur *NftDataRepository) TokenIdExists(ctx context.Context, tokenId int64) (bool, error) {
	const op = "postgresql.NftDataRepository.TokenIdExists"

	query := "SELECT EXISTS(SELECT id FROM nft_data WHERE token_id = $1);"
	var exists bool
	err := ur.db.QueryRow(ctx, query, tokenId).Scan(&exists)
	if err != nil {
		return false, tvoerrors.Wrap(op, err)
	}
	return exists, nil
}
