package models

import (
	"errors"
	"time"
)

type Role string

const (
	RoleAdmin       = "admin"
	RoleUserManager = "usermanager"
	RoleUser        = "user"
)

var (
	ErrInvalidRole = errors.New("invalid role")

	roleMap = map[string]Role{
		"admin":       RoleAdmin,
		"usermanager": RoleUserManager,
		"user":        RoleUser,
	}
)

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
