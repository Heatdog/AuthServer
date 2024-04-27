package tokenservice

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tokenrepository "github.com/Heaterdog/AuthServer/internal/repository/token"
	"github.com/Heaterdog/AuthServer/pkg/jwt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=TokenService
type TokenService interface {
	GenerateToken(ctx context.Context, tokenFields jwt.TokenFileds) (accessToken string, refreshToken string, err error)
	VerifyToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)
}

type tokenService struct {
	logger      *slog.Logger
	secretKey   string
	repo        tokenrepository.TokenRepository
	tokenExpire int
}

func NewTokenService(logger *slog.Logger, repo tokenrepository.TokenRepository,
	secretKey string, tokeExpire int) TokenService {
	return &tokenService{
		logger:      logger,
		repo:        repo,
		secretKey:   secretKey,
		tokenExpire: tokeExpire,
	}
}

func (service *tokenService) GenerateToken(ctx context.Context, tokenFields jwt.TokenFileds) (string, string, error) {
	service.logger.Info("generate tokens", slog.Any("user", tokenFields.ID))

	service.logger.Debug("generate access token", slog.Any("user", tokenFields.ID))
	accessToken, err := jwt.GenerateToken(tokenFields, service.secretKey)
	if err != nil {
		service.logger.Error("generate access token failed", slog.Any("error", err))
		return "", "", err
	}

	expire := time.Duration(service.tokenExpire) * time.Hour * 24
	service.logger.Debug("generate refresh token", slog.Any("user", tokenFields.ID))
	refreshToken, err := jwt.GenerateRefreshToken(tokenFields, service.secretKey, expire)
	if err != nil {
		service.logger.Error("generate refresh token failed", slog.Any("error", err))
		return "", "", err
	}

	service.logger.Debug("set token in storage",
		slog.Any("user", tokenFields.ID),
		slog.String("token", refreshToken),
		slog.Any("expire", expire))
	if err = service.repo.SetToken(ctx, tokenFields.ID, refreshToken, expire); err != nil {
		service.logger.Error("set token in repo failed", slog.Any("error", err))
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (service *tokenService) VerifyToken(ctx context.Context, refreshToken string) (string, string, error) {

	service.logger.Info("verify refresh token", slog.String("token", refreshToken))
	fields, err := jwt.VerifyToken(refreshToken, service.secretKey)
	if err != nil {
		service.logger.Warn("incorrect refresh token", slog.String("token", refreshToken))
		return "", "", err
	}

	storagedToken, err := service.repo.GetToken(ctx, fields.ID)
	if err != nil {
		service.logger.Warn("token does not contain", slog.String("token", refreshToken))
		return "", "", err
	}

	service.logger.Debug("got token from storage", slog.String("token", storagedToken))
	if refreshToken != storagedToken {
		service.logger.Warn("tokens are not equal", slog.String("storagedToken", storagedToken),
			slog.String("newToken", refreshToken))
		return "", "", fmt.Errorf("tokens are not equal")
	}

	service.logger.Info("refresh tokens equal, generate new pair of tokens")
	return service.GenerateToken(ctx, *fields)

}
