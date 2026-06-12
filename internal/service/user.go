package service

import (
	"context"
	"log"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{repo}
}

func (u *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	exists, err := u.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exists != nil {
		err = apperror.New("email already exists")
		log.Println(err)
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	id, err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	res := dto.UserResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}

	return &res, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, userID bson.ObjectID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	usr, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		err = apperror.New("user not found")
		log.Println(err)
		return nil, err
	}

	usr.Name = req.Name
	usr.UpdatedAt = time.Now().UTC()

	if err := u.repo.Update(ctx, *usr); err != nil {
		return nil, err
	}

	res := dto.UserResponse{
		ID:    usr.ID.Hex(),
		Name:  usr.Name,
		Email: usr.Email,
	}

	return &res, nil
}

func (u *UserService) ChangePassword(ctx context.Context, userID bson.ObjectID, req dto.ChangePasswordRequest) error {
	usr, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if usr == nil {
		err = apperror.New("user not found")
		log.Println(err)
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(req.OldPassword)); err != nil {
		err = apperror.New("invalid credentials")
		log.Println(err)
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		log.Println(err)
		return err
	}

	usr.PasswordHash = string(hash)
	usr.UpdatedAt = time.Now().UTC()

	if err := u.repo.Update(ctx, *usr); err != nil {
		return err
	}

	return nil
}
