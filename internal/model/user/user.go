package user

import "github.com/Heaterdog/AuthServer/internal/model/role"

type User struct {
	ID          string
	Login       string
	Role        role.Role
	IsConfirmed bool
	Password    string
}
