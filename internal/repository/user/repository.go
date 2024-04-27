package userrepo

import (
	"context"

	authmodel "github.com/Heaterdog/AuthServer/internal/model/auth"
	"github.com/Heaterdog/AuthServer/internal/model/user"
)

type UserRepository interface {
	Insert(ctx context.Context, user authmodel.AuthRequest) error
	Get(ctx context.Context, login string) (*user.User, error)
	Confirm(ctx context.Context, id, role string) error
}
