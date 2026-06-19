package service

import (
	"context"
	"log"
	"math"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/normalize"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CategoryService struct {
	categoryRepo    repository.CategoryRepository
	transactionRepo repository.TransactionRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, transactionRepository repository.TransactionRepository) CategoryService {
	return CategoryService{categoryRepo, transactionRepository}
}

func (c *CategoryService) Create(ctx context.Context, userId bson.ObjectID, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	normalizedName := normalize.CategoryName(req.Name)

	exists, err := c.categoryRepo.GetByNormalizedName(ctx, userId, normalizedName)
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

	ctId, err := c.categoryRepo.Create(ctx, ct)
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

func (c *CategoryService) GetByFilter(ctx context.Context, userId bson.ObjectID, req dto.ListCategoriesRequest) (*dto.ListCategoriesResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	filter := repository.CategoryFilter{
		UserID: userId,
		Page:   page,
		Limit:  limit,
	}

	ret, err := c.categoryRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	list := []dto.CategoryResponse{}
	for _, ct := range ret.Categories {
		list = append(list, dto.CategoryResponse{
			ID:        ct.ID.Hex(),
			Name:      ct.Name,
			CreatedAt: ct.CreatedAt,
			UpdatedAt: ct.UpdatedAt,
		})
	}

	res := dto.ListCategoriesResponse{
		Data:       list,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: ret.TotalItems,
		TotalPages: int(math.Ceil(float64(ret.TotalItems) / float64(filter.Limit))),
	}

	return &res, nil
}

func (c *CategoryService) UpdateById(ctx context.Context, userId bson.ObjectID, categoryId string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	id, err := bson.ObjectIDFromHex(categoryId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	normalizedName := ""

	if req.Name != nil {
		normalizedName = normalize.CategoryName(*req.Name)
		exists, err := c.categoryRepo.GetByNormalizedName(ctx, userId, normalizedName)
		if err != nil {
			return nil, err
		}

		if exists != nil && exists.ID != id {
			err = apperror.New("category already exists")
			log.Println(err)
			return nil, err
		}
	}

	up := repository.CategoryUpdate{
		Name:           req.Name,
		NormalizedName: normalizedName,
	}

	ct, err := c.categoryRepo.UpdateById(ctx, id, userId, up)
	if err != nil {
		return nil, err
	}

	if ct == nil {
		return nil, nil
	}

	res := dto.CategoryResponse{
		ID:        ct.ID.Hex(),
		Name:      ct.Name,
		CreatedAt: ct.CreatedAt,
		UpdatedAt: ct.UpdatedAt,
	}

	return &res, nil
}

func (c *CategoryService) GetById(ctx context.Context, userId bson.ObjectID, categoryId string) (*dto.CategoryResponse, error) {
	id, err := bson.ObjectIDFromHex(categoryId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ret, err := c.categoryRepo.GetById(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, nil
	}

	res := dto.CategoryResponse{
		ID:        ret.ID.Hex(),
		Name:      ret.Name,
		CreatedAt: ret.CreatedAt,
		UpdatedAt: ret.UpdatedAt,
	}

	return &res, nil
}

func (c *CategoryService) DeleteById(ctx context.Context, userId bson.ObjectID, categoryId string) error {
	id, err := bson.ObjectIDFromHex(categoryId)
	if err != nil {
		log.Println(err)
		return err
	}

	ct, err := c.categoryRepo.GetById(ctx, id, userId)
	if err != nil {
		return err
	}

	if ct == nil {
		err = apperror.New("category not found")
		log.Println(err)
		return err
	}

	filter := repository.TransactionFilter{UserID: userId, CategoryID: &ct.ID}
	res, err := c.transactionRepo.GetByFilter(ctx, filter)
	if err != nil {
		return err
	}

	if res.TotalItems > 0 {
		err = apperror.New("category is linked to one or more transaction")
		log.Println(err)
		return err
	}

	if err := c.categoryRepo.DeleteById(ctx, id, userId); err != nil {
		return err
	}

	return nil
}
