package dto

import (
	"net/mail"
	"strings"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *CreateUserRequest) Validate() error {
	if len(c.Password) < MinPasswordLength || len(c.Password) > MaxPasswordLength {
		return apperror.New("password must be between 8 and 72 characters long")
	}

	c.Email = strings.TrimSpace(strings.ToLower(c.Email))

	_, err := mail.ParseAddress(c.Email)

	if c.Email == "" || len(c.Email) > 254 || err != nil {
		return apperror.New("invalid email")
	}

	if strings.TrimSpace(c.Name) == "" {
		return apperror.New("invalid name")
	}

	return nil
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

func (u *UpdateUserRequest) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return apperror.New("invalid name")
	}

	return nil
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
