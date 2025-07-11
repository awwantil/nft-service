package httpmiddlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"main/tools/pkg/constants"
	httputils "main/tools/pkg/http_utils"
	"main/tools/pkg/logger"
	tvoerrors "main/tools/pkg/tvo_errors"
	tvomodels "main/tools/pkg/tvo_models"
)

type CheckTokenCallback func(ctx context.Context, token string) (*tvomodels.TokenData, error)

// NewAuthMiddleware panic recover middleware
func NewAuthMiddleware(checkFunc CheckTokenCallback, allowUnauth bool, logger *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// extract token from request
		var token string
		authHeader := c.Get(fiber.HeaderAuthorization)
		l := len(constants.AUTH_SCHEMA)
		if len(authHeader) > l+1 && strings.EqualFold(authHeader[:l], constants.AUTH_SCHEMA) {
			token = strings.TrimSpace(authHeader[l:])
		}

		// check token not empty
		if token == "" {
			if allowUnauth {
				// разрешаем доступ пользователям без токена
				c.Locals(constants.TOKEN_DATA_KEY, nil)
				return c.Next()
			} else {
				logger.Error("empty token", "auth_header", authHeader)
				return httputils.HandleError(c, fiber.StatusForbidden, tvoerrors.ErrForbidden)
			}
		}

		var err error
		tokenData := &tvomodels.TokenData{}

		// check token
		if tokenData, err = checkFunc(c.Context(), token); err != nil {
			logger.Error("check token error", "token", token, "error", err)
			return httputils.HandleError(c, fiber.StatusUnauthorized, tvoerrors.ErrInvalidJWT)
		}

		tokenData.RawToken = token

		// save token data to request
		c.Locals(constants.TOKEN_DATA_KEY, *tokenData)

		// next step or error
		return c.Next()
	}
}
