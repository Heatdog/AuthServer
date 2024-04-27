package authservice

import (
	"context"
	"fmt"
	"log/slog"

	authmodel "github.com/Heaterdog/AuthServer/internal/model/auth"
	userrepo "github.com/Heaterdog/AuthServer/internal/repository/user"
	tokenservice "github.com/Heaterdog/AuthServer/internal/service/token"
	cryptohash "github.com/Heaterdog/AuthServer/pkg/cryptoHash"
	"github.com/Heaterdog/AuthServer/pkg/jwt"
)

type AuthService interface {
	SignIn(ctx context.Context, user authmodel.AuthRequest) (string, string, error)
	SignUp(ctx context.Context, user authmodel.AuthRequest) error
	Confirm(ctx context.Context, id, role string) error
}

type authService struct {
	logger       *slog.Logger
	repo         userrepo.UserRepository
	tokenService tokenservice.TokenService
}

func NewAuthService(repo userrepo.UserRepository, token tokenservice.TokenService, logger *slog.Logger) AuthService {
	return &authService{
		logger:       logger,
		repo:         repo,
		tokenService: token,
	}
}

func (service *authService) SignIn(ctx context.Context, user authmodel.AuthRequest) (string, string, error) {
	service.logger.Info("sign in", slog.String("user", user.Login))
	service.logger.Debug("get user from repo", slog.String("user", user.Login))

	res, err := service.repo.Get(ctx, user.Login)
	if err != nil {
		service.logger.Warn("user repo error", slog.Any("error", err))
		return "", "", err
	}

	if !res.IsConfirmed {
		err := fmt.Errorf("incomfired profile")
		service.logger.Warn("user repo error", slog.Any("error", err))
		return "", "", err
	}

	service.logger.Debug("verify password", slog.String("user", user.Login))
	if cryptohash.VerifyHash([]byte(res.Password), user.Password) {
		service.logger.Debug("generate tokens", slog.String("user", user.Login))
		accessToken, refreshToken, err := service.tokenService.GenerateToken(ctx, jwt.TokenFileds{
			ID:    res.ID,
			Login: res.Login,
			Role:  string(res.Role),
		})
		if err != nil {
			service.logger.Warn("jwt token generate failed", slog.Any("error", err))
			return "", "", err
		}
		return accessToken, refreshToken, nil
	}

	errStr := fmt.Sprint("wrong password ", slog.Any("error", user.Login))
	service.logger.Info(errStr)
	return "", "", fmt.Errorf(errStr)
}

func (service *authService) SignUp(ctx context.Context, user authmodel.AuthRequest) error {
	service.logger.Info("sign up", slog.String("user", user.Login))

	pswd, err := cryptohash.Hash(user.Password)
	if err != nil {
		service.logger.Warn(err.Error())
		return err
	}

	user.Password = string(pswd)
	return service.repo.Insert(ctx, user)
}

func (service *authService) Confirm(ctx context.Context, id, role string) error {
	service.logger.Info("confirm handler", slog.String("id", id), slog.String("role", role))

}
