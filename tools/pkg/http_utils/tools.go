package httputils

import (
	"context"
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/metadata"

	"main/tools/pkg/constants"
	"main/tools/pkg/helpers"
	"main/tools/pkg/logger"
	tvoerrors "main/tools/pkg/tvo_errors"
	tvomodels "main/tools/pkg/tvo_models"
)

// errorResponse represents a JSON error response.
type errorResponse struct {
	Message string `json:"message"`
}

// HandleError handles the error response for various status codes.
func HandleError(c *fiber.Ctx, statusCode int, e error) error {
	if err := c.Status(statusCode).JSON(errorResponse{
		Message: e.Error(),
	}); err != nil {
		return tvoerrors.Wrap("sending response", err)
	}
	return nil
}

// ParseRequestBody утилита для парсинга входящего запроса
func ParseRequestBody(c *fiber.Ctx, req interface{}, method string, logger *logger.Logger) error {
	if err := c.BodyParser(req); err != nil {
		// TODO: add sentry send
		logger.Error("cant parse body", "error", err, "method", method)
		return err
	}
	return nil
}

// ParseIntQueryParam парсинг параметра QueryString
func ParseIntQueryParam(c *fiber.Ctx, key string, logger *logger.Logger) (int, error) {
	valueString := c.Params(key)
	value, err := strconv.Atoi(valueString)
	if err != nil || value == 0 {
		logger.Error("invalid data", "error", err, "value", valueString)
		return 0, err
	}
	return value, nil
}

// CtxWithAuthToken утилита для подмешивания токена в контекст gRPC
func CtxWithAuthToken(c *fiber.Ctx) context.Context {
	var rawToken string
	tokenData, ok := c.Locals(constants.TOKEN_DATA_KEY).(tvomodels.TokenData)
	if ok {
		rawToken = tokenData.RawToken
	}

	return metadata.NewOutgoingContext(c.Context(), metadata.Pairs("authorization", rawToken))
}

// UserIDFromToken extracts the user ID from the JWT token stored in the Fiber context
func UserIDFromToken(c *fiber.Ctx, method string, logger *logger.Logger) (int64, error) {
	tokenData, ok := c.Locals(constants.TOKEN_DATA_KEY).(tvomodels.TokenData)
	if !ok {
		logger.Error("can't cast user data", "error", "can't extract user data from context", "method", method)
		return 0, tvoerrors.Wrap("can't extract user data from context", tvoerrors.ErrInvalidRequestData)
	}
	return tokenData.UserID, nil
}

// RoleIDFromToken extracts the role ID from the JWT token stored in the Fiber context
func RoleIDFromToken(c *fiber.Ctx, method string, logger *logger.Logger) (int64, error) {
	tokenData, ok := c.Locals(constants.TOKEN_DATA_KEY).(tvomodels.TokenData)
	if !ok {
		logger.Error("can't cast user data", "error", "can't extract user data from context", "method", method)
		return 0, tvoerrors.Wrap("can't extract user data from context", tvoerrors.ErrInvalidRequestData)
	}
	return int64(tokenData.UserRoleID), nil
}

// IsAuthorized checks if the request is authorized
func IsAuthorized(c *fiber.Ctx) bool {
	return c.Locals(constants.TOKEN_DATA_KEY) != nil
}

// GetTokenDataFromCtx retrieves the authenticated user's token data from the context.
func GetTokenDataFromCtx(ctx context.Context) (*tvomodels.TokenData, error) {
	tokenData, ok := ctx.Value(constants.TOKEN_DATA_KEY).(*tvomodels.TokenData)
	if !ok {
		return nil, tvoerrors.Wrap("invalid token", tvoerrors.ErrCastClaims)
	}

	return tokenData, nil
}

func GetTelegramSignatureFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", tvoerrors.Wrap("invalid token", tvoerrors.ErrCastClaims)
	}

	signature := ""
	if values := md.Get("x-telegram-signature"); len(values) > 0 {
		signature = values[0]
	}

	return signature, nil
}

// GetBytesFromMultipartFile retrieves the bytes data from a multipart file in the request.
func GetBytesFromMultipartFile(c *fiber.Ctx, logger *logger.Logger, key string, required bool) ([]byte, error) {
	fileHeader, err := c.FormFile(key)

	if err != nil && !required {
		return nil, nil
	}

	if err != nil {
		logger.Error("Failed get file from request", "key", key, "error", err)
		return nil, tvoerrors.Wrap("Failed get file from request", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Error("Failed to open file", "key", key, "error", err)
		return nil, tvoerrors.Wrap("Failed to open file", err)
	}
	defer helpers.DeferClose(file, logger)

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read file", "key", key, "error", err)
		return nil, tvoerrors.Wrap("Failed to read file", err)
	}
	return bytes, nil
}
