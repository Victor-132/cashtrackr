package service

import (
	"context"
	"log"
	"strings"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/auth"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return AuthService{repo}
}

func (a *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	usr, err := a.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		err = apperror.New("user not found")
		log.Println(err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(req.Password)); err != nil {
		err = apperror.New("invalid credentials")
		log.Println(err)
		return nil, err
	}

	token, err := auth.GenerateJWT(usr.ID.Hex())
	if err != nil {
		return nil, err
	}

	res := dto.LoginResponse{
		Token: token,
	}

	return &res, nil
}
