package postgreauth

import (
	"context"
	"log/slog"

	authmodel "github.com/Heaterdog/AuthServer/internal/model/auth"
	"github.com/Heaterdog/AuthServer/internal/model/user"
	user_repo "github.com/Heaterdog/AuthServer/internal/repository/user"
	client "github.com/Heaterdog/AuthServer/pkg/postgre"
	"github.com/jackc/pgx/v5"
)

type authRespository struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewUserPostgreRepository(dbClient client.Client, logger *slog.Logger) user_repo.UserRepository {
	return &authRespository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *authRespository) Insert(ctx context.Context, user authmodel.AuthRequest) error {
	repo.logger.Debug("insert user in repo")

	q := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
	`
	repo.logger.Debug("user repo query", slog.String("query", q))

	tag, err := repo.dbClient.Exec(ctx, q, user.Login, user.Password)
	if err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return err
	}

	if tag.RowsAffected() != 1 {
		err = pgx.ErrNoRows
		repo.logger.Error("SQL error", slog.Any("error", err))
		return err
	}

	return nil
}

func (repo *authRespository) Get(ctx context.Context, login string) (*user.User, error) {
	repo.logger.Debug("get user")

	q := `
		SELECT id, login, password, role, is_confirmed
		FROM users
		WHERE login = $1
	`

	repo.logger.Debug("user repo query", slog.String("query", q))
	row := repo.dbClient.QueryRow(ctx, q, login)

	var res user.User

	if err := row.Scan(&res.ID, &res.Login, &res.Password, &res.Role, &res.IsConfirmed); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return nil, err
	}

	return &res, nil
}

func (repo *authRespository) Confirm(ctx context.Context, id, role string) error {
	repo.logger.Debug("confirm user")

	q := `
		UPDATE users
		SET is_confirm = TRUE, role = $1
		WHERE id = $2 
	`
	repo.logger.Debug("user repo query", slog.String("query", q))

	tag, err := repo.dbClient.Exec(ctx, q, role, id)
	if err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return err
	}

	if tag.RowsAffected() != 1 {
		err = pgx.ErrNoRows
		repo.logger.Error("SQL error", slog.Any("error", err))
		return err
	}

	return nil
}
