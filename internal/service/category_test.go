package service

import (
	"context"
	"errors"
	"testing"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"github.com/go-openapi/testify/v2/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCategoryService_Create_Success(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByNormalizedNameFn: func(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
			return nil, nil
		},
		CreateFn: func(ctx context.Context, ct model.Category) (string, error) {
			return "", nil
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewCategoryService(ctRepo, trRepo)

	req := dto.CreateCategoryRequest{
		Name: "Alimentação",
	}

	res, err := svc.Create(ctx, userId, req)

	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, req.Name, res.Name)
}

func TestCategoryService_Create_Duplicated_Name(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByNormalizedNameFn: func(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
			return &model.Category{}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewCategoryService(ctRepo, trRepo)

	req := dto.CreateCategoryRequest{
		Name: "Alimentação",
	}

	_, err := svc.Create(ctx, userId, req)

	require.Error(t, err)
	require.Equal(t, "category already exists", err.(apperror.AppError).Error())
}

func TestCategoryService_Create_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByNormalizedNameFn: func(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
			return nil, nil
		},
		CreateFn: func(ctx context.Context, ct model.Category) (string, error) {
			return "", errors.New("")
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewCategoryService(ctRepo, trRepo)

	req := dto.CreateCategoryRequest{
		Name: "Alimentação",
	}

	_, err := svc.Create(ctx, userId, req)

	require.Error(t, err)
}

func TestCategoryService_Update_Success(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByNormalizedNameFn: func(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
			return nil, nil
		},
		UpdateByIdFn: func(ctx context.Context, id, userId bson.ObjectID, req repository.CategoryUpdate) (*model.Category, error) {
			return &model.Category{
				Name: *req.Name,
			}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewCategoryService(ctRepo, trRepo)

	name := "Combustível"
	req := dto.UpdateCategoryRequest{
		Name: &name,
	}

	res, err := svc.UpdateById(ctx, userId, categoryId.Hex(), req)

	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, *req.Name, res.Name)
}

func TestCategoryService_Update_Duplicated_Name(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByNormalizedNameFn: func(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
			return &model.Category{}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewCategoryService(ctRepo, trRepo)

	name := "Combustível"
	req := dto.UpdateCategoryRequest{
		Name: &name,
	}

	_, err := svc.UpdateById(ctx, userId, categoryId.Hex(), req)

	require.Error(t, err)
	require.Equal(t, "category already exists", err.(apperror.AppError).Error())
}

func TestCategoryService_Update_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByNormalizedNameFn: func(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
			return nil, nil
		},
		UpdateByIdFn: func(ctx context.Context, id, userId bson.ObjectID, req repository.CategoryUpdate) (*model.Category, error) {
			return nil, errors.New("")
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewCategoryService(ctRepo, trRepo)

	name := "Combustível"
	req := dto.UpdateCategoryRequest{
		Name: &name,
	}

	_, err := svc.UpdateById(ctx, userId, categoryId.Hex(), req)

	require.Error(t, err)
}

func TestCategoryService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return &model.Category{}, nil
		},
		DeleteByIdFn: func(ctx context.Context, id, userId bson.ObjectID) error {
			return nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		GetByFilterFn: func(ctx context.Context, trFilter repository.TransactionFilter) (*repository.PaginatedTransactions, error) {
			return &repository.PaginatedTransactions{}, nil
		},
	}

	svc := NewCategoryService(ctRepo, trRepo)

	err := svc.DeleteById(ctx, userId, categoryId.Hex())

	require.NoError(t, err)
}

func TestCategoryService_Delete_Category_Not_Found(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return nil, nil
		},
		DeleteByIdFn: func(ctx context.Context, id, userId bson.ObjectID) error {
			return nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		GetByFilterFn: func(ctx context.Context, trFilter repository.TransactionFilter) (*repository.PaginatedTransactions, error) {
			return &repository.PaginatedTransactions{}, nil
		},
	}

	svc := NewCategoryService(ctRepo, trRepo)

	err := svc.DeleteById(ctx, userId, categoryId.Hex())

	require.Error(t, err)
	require.Equal(t, "category not found", err.(apperror.AppError).Error())
}

func TestCategoryService_Delete_Linked_Transactions(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return &model.Category{}, nil
		},
		DeleteByIdFn: func(ctx context.Context, id, userId bson.ObjectID) error {
			return nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		GetByFilterFn: func(ctx context.Context, trFilter repository.TransactionFilter) (*repository.PaginatedTransactions, error) {
			return &repository.PaginatedTransactions{TotalItems: 1}, nil
		},
	}

	svc := NewCategoryService(ctRepo, trRepo)

	err := svc.DeleteById(ctx, userId, categoryId.Hex())

	require.Error(t, err)
	require.Equal(t, "category is linked to one or more transaction", err.(apperror.AppError).Error())
}

func TestCategoryService_Delete_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return &model.Category{}, nil
		},
		DeleteByIdFn: func(ctx context.Context, id, userId bson.ObjectID) error {
			return errors.New("")
		},
	}

	trRepo := &TransactionRepositoryMock{
		GetByFilterFn: func(ctx context.Context, trFilter repository.TransactionFilter) (*repository.PaginatedTransactions, error) {
			return &repository.PaginatedTransactions{}, nil
		},
	}

	svc := NewCategoryService(ctRepo, trRepo)

	err := svc.DeleteById(ctx, userId, categoryId.Hex())

	require.Error(t, err)
}
