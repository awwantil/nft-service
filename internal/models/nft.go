package models

import "time"

type NftDataModel struct {
	ID            int64     `json:"id"`
	TokenId       int64     `json:"token_id" example:"1"`
	Description   string    `json:"description" example:"About this token"`
	CidV0         string    `json:"cid_v0" example:"dss"`
	CidV1         string    `json:"cid_v1" example:"dss"`
	FileName      string    `json:"file_name" example:"pic12.png"`
	FileSize      string    `json:"file_size" example:"12kb"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	DeletedAt     time.Time `json:"-"`
	LastVisitedAt time.Time `json:"-"`
}
