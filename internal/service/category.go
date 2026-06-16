package service

import (
	"context"
	"log"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/normalize"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return CategoryService{repo}
}

func (c *CategoryService) Create(ctx context.Context, userId bson.ObjectID, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	normalizedName := normalize.CategoryName(req.Name)

	exists, err := c.repo.GetByNormalizedName(ctx, userId, normalizedName)
	if err != nil {
		return nil, err
	}

	if exists != nil {
		err = apperror.New("category already exists")
		log.Println(err)
		return nil, err
	}

	ct := model.Category{
		UserID:         userId,
		Name:           req.Name,
		NormalizedName: normalizedName,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	ctId, err := c.repo.Create(ctx, ct)
	if err != nil {
		return nil, err
	}

	res := dto.CategoryResponse{
		ID:        ctId,
		Name:      ct.Name,
		CreatedAt: ct.CreatedAt,
		UpdatedAt: ct.UpdatedAt,
	}

	return &res, nil
}
