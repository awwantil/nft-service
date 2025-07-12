// main.go
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"main/internal/config"
	jwtManager "main/internal/lib/jwt"
	"main/internal/repository/postgresql"
	"main/internal/server"
	rediscache "main/tools/pkg/cache/redis"
	coreconfig "main/tools/pkg/core_config"
	"main/tools/pkg/database"
	"main/tools/pkg/logger"
	"os"
	"os/signal"

	"main/internal/handlers"
)

func main() {
	_, ctx := errgroup.WithContext(context.Background())

	go func() {
		if recoveryMessage := recover(); recoveryMessage != nil {
			fmt.Println(recoveryMessage)
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// создаем объект конфига приложения
	var cfg config.Config

	// инициализация конфига
	if err := coreconfig.Load(&cfg, ""); err != nil {
		log.Panic("Can't load config file", err)
	}

	// создаем логгер
	logger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Panic("logger initialization error ", err)
	}

	// подключаемся к БД
	db, err := database.NewClient(ctx, &cfg.Database)
	if err != nil {
		log.Panic("Failed to database connection", err)
	}
	defer db.Close()

	// создаем клиент для кеша
	cacheClient, err := rediscache.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Panic("cache error: ", err)
	}
	defer func() {
		if err = cacheClient.Close(); err != nil {
			logger.Error("Error closing Redis connection:", "error", err)
		}
	}()

	// инициализируем репозитории
	userRepository := postgresql.NewUserRepository(db, cfg.Secret)
	tokenRepository := postgresql.NewUserTokenRepository(db, cfg.App.Debug)
	roleRepository := postgresql.NewRoleRepository(db)
	nftDataRepository := postgresql.NewNftDataRepository(db)
	jwt := jwtManager.NewJWTManager(&cfg.JWT)

	logger.Info("Create server")

	app := server.NewServer()
	logger.Info("Creating internal handlers")
	authHandlers := handlers.NewAuthHandlers(logger, jwt, userRepository, tokenRepository, roleRepository, cacheClient, cfg.Secret)
	kuboHandlers := handlers.NewKuboHandlers(logger)
	nftDataHandlers := handlers.NewNftHandlers(logger, nftDataRepository)

	// добавляем роуты для экземпляра сервера
	server.AddRoutes(app, authHandlers, kuboHandlers, nftDataHandlers, logger)

	logger.Info("Service api gateway starts", "address", cfg.App.Addr)
	if err = app.Listen(cfg.App.Addr); err != nil {
		logger.Error("listen error", "error", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	if err = app.Shutdown(); err != nil {
		logger.Info("api-gateway service shutdown")
		return
	}

	logger.Info("api-gateway service was stopped")

}
