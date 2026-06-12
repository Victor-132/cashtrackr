package dto

import apperror "github.com/Victor-132/cashtrackr/internal/app_error"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

func (c ChangePasswordRequest) Validate() error {
	if len(c.NewPassword) < MinPasswordLength || len(c.NewPassword) > MaxPasswordLength {
		return apperror.New("password must be between 8 and 72 characters long")
	}

	return nil
}
