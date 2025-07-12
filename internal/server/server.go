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

// AddRoutes добавляет роуты к серверу
func AddRoutes(app *fiber.App, authHandlers *handlers.AuthHandlers, kuboHandlers *handlers.KuboHandlers,
	logger *logger.Logger) {
	// helpful middlewares
	app.Use(healthcheck.New())

	// добавляем API v1
	// v1Router := app.Group("/v1", slogfiber.New(logger.Logger), recover.New())
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

	//router := app.Use()

	addRoutesV1(v1Router, authHandlers, kuboHandlers, logger)
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
	logger *logger.Logger) fiber.Router {
	authMiddleware := httpmiddlewares.NewAuthMiddleware(checkAuthToken(logger), false, logger)
	//guestMiddleware := httpmiddlewares.NewAuthMiddleware(checkAuthToken(logger), true, logger)

	// системные урлы

	idm := v1Router.Group("/auth")

	// публичные методы
	idm.Post("/registration/", httputils.FiberJSONWrapper(authHandlers.Registration))
	idm.Post("/login/", httputils.FiberJSONWrapper(authHandlers.Login))
	idm.Post("/refresh/", httputils.FiberJSONWrapper(authHandlers.Refresh))
	idm.Post("/recovery/", httputils.FiberJSONWrapper(authHandlers.Recovery))
	idm.Post("/ping/", httputils.FiberJSONWrapper(authHandlers.Ping))

	// методы под авторизацией
	idmProtected := idm.Group("")
	idmProtected = idmProtected.Use(authMiddleware)
	idmProtected.Post("/logout/", httputils.FiberJSONWrapper(authHandlers.Logout))
	idmProtected.Post("/delete_user/", httputils.FiberJSONWrapper(authHandlers.DeleteUser))
	idmProtected.Post("/digup_user/", httputils.FiberJSONWrapper(authHandlers.DigupUser))
	idmProtected.Post("/update/", httputils.FiberJSONWrapper(authHandlers.UpdateUser))
	idmProtected.Post("/change_role/", httputils.FiberJSONWrapper(authHandlers.ChangeRole))
	idmProtected.Post("/reset_token/", httputils.FiberJSONWrapper(authHandlers.ResetToken))

	// методы сервиса UPS
	api := v1Router.Group("/api/v1")
	api.Get("/pins", handlers.ListPinsHandler)

	apiProtected := v1Router.Group("/api/v1", authMiddleware)
	apiProtected.Post("/files", handlers.UploadFileHandler)
	// Маршруты для управления закреплением (pin)
	apiProtected.Post("/pins/:cid", handlers.PinCidHandler)
	apiProtected.Delete("/pins/:cid", handlers.UnpinCidHandler)

	return v1Router
}
