package tokenrepository

import (
	"context"
	"time"
)

type TokenRepository interface {
	SetToken(ctx context.Context, userId, token string, expire time.Duration) error
	GetToken(ctx context.Context, userId string) (string, error)
}
