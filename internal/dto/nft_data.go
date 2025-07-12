package dto

import "mime/multipart"

type CreateNftDataRequest struct {
	NftId       string                `json:"nft_id" example:"1"`
	Description string                `json:"description" example:"About this token"`
	ImageFile   *multipart.FileHeader `json:"image_file" form:"image_file" example:"pic12.png"`
}

type NftData struct {
	TokenId     int64  `json:"token_id" example:"1"`
	Description string `json:"description" example:"About this token"`
	CidV0       string `json:"cid_v0" example:"dss"`
	CidV1       string `json:"cid_v1" example:"dss"`
	FileName    string `json:"file_name" example:"pic12.png"`
	FileSize    string `json:"file_size" example:"12kb"`
}

type CreateNftDataResponse struct {
	Message string `json:"message"`
}
