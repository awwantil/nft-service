package handlers

import (
	"fmt"
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

	if err := httputils.ParseRequestBody(c, &request, "CreateNftData", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}
	file, err := c.FormFile("file")
	if err != nil {
		log.Error("Error reading image file", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	ctx := httputils.CtxWithAuthToken(c)
	roleId, err := httputils.RoleIDFromToken(c, "CreateNftData", h.logger)
	if err != nil {
		return nil, tvoerrors.ErrCastClaims
	}

	if roleId != 100 {
		log.Error("Wrong user role")
		return nil, tvoerrors.ErrForbidden
	}

	isExist, err := h.nftDataRepository.TokenIdExists(ctx, request.Id)
	if err != nil {
		log.Error("Error accessing to DB", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}
	if isExist {
		log.Error("Wrong token id", "error", err)
		return nil, status.Error(codes.Internal, "wrong token id (is exist)") //nolint
	}
	addResponse, cidV1, _, err := service.AddFileToIPFS(file)
	if err != nil {
		log.Error("Error creating nft data ", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	nftData := &dto.NftData{
		TokenId:     request.Id,
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

func (h *NftHandlers) ReadNft(c *fiber.Ctx) (interface{}, error) {
	strId := c.Params("id")
	if strId == "" {
		log.Error("Error reading nft id", "error")
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}
	tokenId, err := strconv.ParseInt(strId, 10, 64)

	if err != nil {
		log.Error("Error parsing nft id", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	ctx := c.Context()

	nft, err := h.nftDataRepository.ReadNftData(ctx, tokenId)
	if err != nil {
		log.Error("Error accessing to DB", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}
	return &dto.ReadNftResponse{
		Info: &dto.NftInfo{
			TokenId:     nft.TokenId,
			Description: nft.Description,
			CidV0:       nft.CidV0,
			CidV1:       nft.CidV1,
			Link:        fmt.Sprintf(service.KuboGatewayUrlTemplate, nft.CidV1),
		},
	}, nil
}

func (h *NftHandlers) ReadAllNft(c *fiber.Ctx) (interface{}, error) {
	strLimit := c.Params("limit")
	if strLimit == "" {
		log.Error("Error reading limit", "error")
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}
	limit, err := strconv.ParseInt(strLimit, 10, 64)

	if err != nil {
		log.Error("Error parsing limit", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	ctx := c.Context()

	nfts, err := h.nftDataRepository.ReadAllNftData(ctx, int(limit))
	if err != nil {
		log.Error("Error accessing to DB", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	infos := []dto.NftInfo{}
	if len(nfts) > 0 {
		for _, nft := range nfts {
			infos = append(infos, dto.NftInfo{
				TokenId:     nft.TokenId,
				Description: nft.Description,
				CidV0:       nft.CidV0,
				CidV1:       nft.CidV1,
				Link:        fmt.Sprintf(service.KuboGatewayUrlTemplate, nft.CidV1),
			})
		}
	}

	return &dto.ReadAllNftResponse{
		Infos: &infos,
	}, nil
}
