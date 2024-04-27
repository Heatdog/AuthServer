package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Heaterdog/AuthServer/internal/config"
	"github.com/Heaterdog/AuthServer/internal/migration"
	tokendb "github.com/Heaterdog/AuthServer/internal/repository/token/redis"
	postgreauth "github.com/Heaterdog/AuthServer/internal/repository/user/postgre"
	authservice "github.com/Heaterdog/AuthServer/internal/service/auth"
	tokenservice "github.com/Heaterdog/AuthServer/internal/service/token"
	"github.com/Heaterdog/AuthServer/internal/transport/auth"
	postgre "github.com/Heaterdog/AuthServer/pkg/postgre/db"
	redisStorage "github.com/Heaterdog/AuthServer/pkg/redis"
	"github.com/gorilla/mux"
)

func App() {
	opt := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opt))
	slog.SetDefault(logger)

	ctx := context.Background()
	logger.Info("reading server config files")
	cfg := config.NewConfigStorage(logger)

	logger.Info("connecting to DataBase")
	dbClient, err := postgre.NewPostgreClient(ctx, cfg.Postgre)
	if err != nil {
		logger.Error("connection to PostgreSQL failed", slog.Any("error", err))
	}
	defer dbClient.Close()

	logger.Info("init database with users")
	if err := migration.InitDb(dbClient); err != nil {
		logger.Error("init db failed", slog.Any("error", err))
	}

	logger.Info("connecting to TokenDataBase")
	redisClient, err := redisStorage.NewRedisClient(ctx, &cfg.Redis)
	if err != nil {
		logger.Error("connection to Redis failed", slog.Any("error", err))
	}
	defer redisClient.Close()

	router := mux.NewRouter()

	tokenRepo := tokendb.NewRedisTokenRepository(logger, redisClient)
	tokenService := tokenservice.NewTokenService(logger, tokenRepo, cfg.PasswordKey, cfg.Redis.TokenExpire)

	authRepo := postgreauth.NewUserPostgreRepository(dbClient, logger)
	authService := authservice.NewAuthService(authRepo, tokenService, logger)
	authHandler := auth.NewTokenHandler(logger, authService)
	authHandler.Register(router)

	host := fmt.Sprintf("%s:%d", cfg.Server.IP, cfg.Server.Port)

	logger.Info("listening server", slog.String("host", host))
	if err := http.ListenAndServe(host, router); err != nil {
		logger.Error("server listnenig failed", slog.Any("error", err))
	}
}
