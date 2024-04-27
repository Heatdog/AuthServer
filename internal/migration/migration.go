package migration

import (
	"context"

	"github.com/Heaterdog/AuthServer/internal/model/user"
	cryptohash "github.com/Heaterdog/AuthServer/pkg/cryptoHash"
	client "github.com/Heaterdog/AuthServer/pkg/postgre"
)

func InitDb(client client.Client) error {
	users := []user.User{
		{
			Login:       "Admin",
			Password:    "1234",
			Role:        "Admin",
			IsConfirmed: true,
		},
		{
			Login:       "Heater",
			Password:    "2345",
			Role:        "Worker",
			IsConfirmed: true,
		},
	}

	q := `
			INSERT INTO Users
				(login, password, role, is_confirmed) 
			VALUES 
				($1, $2, $3, $4)
	`

	for _, el := range users {
		pswd, err := cryptohash.Hash(el.Password)
		if err != nil {
			return err
		}
		client.QueryRow(context.Background(), q, el.Login, string(pswd), el.Role, el.IsConfirmed)
	}
	return nil
}
