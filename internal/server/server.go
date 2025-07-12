package server

import (
	"context"
	"log/slog"
	"main/internal/dto"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"main/internal/handlers"
	httpmiddlewares "main/tools/pkg/http_middlewares"
	httputils "main/tools/pkg/http_utils"
	"main/tools/pkg/logger"
	tvoerrors "main/tools/pkg/tvo_errors"
	tvomodels "main/tools/pkg/tvo_models"
)

func NewServer() *fiber.App {
	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
		WriteTimeout:      time.Second * 15,
		ReadTimeout:       time.Second * 15,
		IdleTimeout:       time.Second * 20,
		CaseSensitive:     true,
		StrictRouting:     false,
		ServerHeader:      "Apache 2.0",
		AppName:           "API Gateway",
		BodyLimit:         20 * 1024 * 1024,
	})

	return app
}

func AddRoutes(app *fiber.App, authHandlers *handlers.AuthHandlers, kuboHandlers *handlers.KuboHandlers,
	nftHandlers *handlers.NftHandlers, logger *logger.Logger) {
	app.Use(healthcheck.New())

	v1Router := app.Group("/v1", slogfiber.NewWithConfig(logger.Logger, slogfiber.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      true,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         true,
		WithTraceID:        true,
	}), recover.New())

	addRoutesV1(v1Router, authHandlers, kuboHandlers, nftHandlers, logger)
}

// checkAuthToken утилита для проверки токена
func checkAuthToken(logger *logger.Logger) httpmiddlewares.CheckTokenCallback {
	return func(ctx context.Context, token string) (*tvomodels.TokenData, error) {
		res, err := handlers.AuthHandler.CheckToken(ctx, &dto.CheckTokenRequest{
			Token: token,
		})

		if err != nil {
			logger.Error("idmClient.CheckToken error", "token", token, "error", err)
			return nil, err
		}

		if !res.IsValid {
			logger.Error("idmClient.CheckToken error", "token", token, "error", "invalid token")
			return nil, tvoerrors.ErrInvalidJWT
		}

		return &tvomodels.TokenData{
			UserID:     res.UserId,
			UserRoleID: tvomodels.RoleId(res.RoleId),
			UserPhone:  res.Phone,
			RawToken:   token,
		}, nil
	}
}

// addRoutesV1 добавляем роутинг для версии API v1
func addRoutesV1(v1Router fiber.Router, authHandlers *handlers.AuthHandlers, kuboHandlers *handlers.KuboHandlers,
	nftHandlers *handlers.NftHandlers, logger *logger.Logger) fiber.Router {
	authMiddleware := httpmiddlewares.NewAuthMiddleware(checkAuthToken(logger), false, logger)
	//guestMiddleware := httpmiddlewares.NewAuthMiddleware(checkAuthToken(logger), true, logger)

	auth := v1Router.Group("/auth")

	// публичные методы
	auth.Post("/registration/", httputils.FiberJSONWrapper(authHandlers.Registration))
	auth.Post("/login/", httputils.FiberJSONWrapper(authHandlers.Login))
	auth.Post("/refresh/", httputils.FiberJSONWrapper(authHandlers.Refresh))
	auth.Post("/recovery/", httputils.FiberJSONWrapper(authHandlers.Recovery))
	auth.Post("/ping/", httputils.FiberJSONWrapper(authHandlers.Ping))

	// методы под авторизацией
	authProtected := auth.Group("")
	authProtected = authProtected.Use(authMiddleware)
	authProtected.Post("/logout/", httputils.FiberJSONWrapper(authHandlers.Logout))
	authProtected.Post("/delete_user/", httputils.FiberJSONWrapper(authHandlers.DeleteUser))
	authProtected.Post("/digup_user/", httputils.FiberJSONWrapper(authHandlers.DigupUser))
	authProtected.Post("/update/", httputils.FiberJSONWrapper(authHandlers.UpdateUser))
	authProtected.Post("/change_role/", httputils.FiberJSONWrapper(authHandlers.ChangeRole))
	authProtected.Post("/reset_token/", httputils.FiberJSONWrapper(authHandlers.ResetToken))

	// методы сервиса UPS
	api := v1Router.Group("/api")
	api.Get("/pins", handlers.ListPinsHandler)
	api.Post("/nft_data", httputils.FiberJSONWrapper(nftHandlers.CreateNftData))

	apiProtected := v1Router.Group("", authMiddleware)
	//apiProtected.Post("/nft_data", httputils.FiberJSONWrapper(nftHandlers.CreateNftData))

	apiProtected.Post("/files", handlers.UploadFileHandler)
	// Маршруты для управления закреплением (pin)
	apiProtected.Post("/pins/:cid", handlers.PinCidHandler)
	apiProtected.Delete("/pins/:cid", handlers.UnpinCidHandler)

	return v1Router
}
