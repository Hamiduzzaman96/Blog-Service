package domain

import (
	"errors"
	"strings"
)

const (
	RoleUser   = "USER"
	RoleAuthor = "AUTHOR"
)

type User struct {
	ID       uint
	Email    string
	Password string
	Role     string
}

func NewUser(email, password string) (*User, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return nil, errors.New("email and password cannot be empty")
	}

	u := &User{
		Email:    email,
		Password: password,
		Role:     RoleUser,
	}

	return u, nil
}

func (u *User) PromoteToAuthor() {
	u.Role = RoleAuthor
}
