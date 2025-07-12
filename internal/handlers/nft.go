package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"main/internal/dto"
	"main/internal/repository"
	"main/internal/service"
	httputils "main/tools/pkg/http_utils"
	"main/tools/pkg/logger"
	tvoerrors "main/tools/pkg/tvo_errors"
	"strconv"
)

// NftHandlers
type NftHandlers struct {
	logger            *logger.Logger
	nftDataRepository repository.NftDataRepository
}

func NewNftHandlers(logger *logger.Logger, nftRepository repository.NftDataRepository) *NftHandlers {
	return &NftHandlers{
		logger:            logger,
		nftDataRepository: nftRepository,
	}
}

func (h *NftHandlers) CreateNftData(c *fiber.Ctx) (interface{}, error) {
	var request dto.CreateNftDataRequest

	ctx := httputils.CtxWithAuthToken(c)

	if err := httputils.ParseRequestBody(c, &request, "CreateNftData", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	file, err := c.FormFile("image_file")
	if err != nil {
		log.Error("Error reading image file", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	addResponse, cidV1, _, err := service.AddFileToIPFS(file)
	if err != nil {
		log.Error("Error creating nft data ", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}
	tokenId, err := strconv.ParseInt(request.NftId, 10, 64)
	if err != nil {
		log.Error("Error parsing token id ", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}
	nftData := &dto.NftData{
		TokenId:     tokenId,
		Description: request.Description,
		CidV0:       addResponse.Hash,
		CidV1:       cidV1,
		FileName:    addResponse.Name,
		FileSize:    addResponse.Size,
	}

	err = h.nftDataRepository.CreateNftData(ctx, nftData)
	if err != nil {
		log.Error("Error creating nft data", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.CreateNftDataResponse{
		Message: "NFT data created successful",
	}, nil
}
