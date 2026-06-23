package service

import (
	"context"
	"errors"
	"testing"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"github.com/go-openapi/testify/v2/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionRepositoryMock struct {
	CreateFn      func(ctx context.Context, tr model.Transaction) (string, error)
	GetByFilterFn func(ctx context.Context, trFilter repository.TransactionFilter) (*repository.PaginatedTransactions, error)
	GetByIdFn     func(ctx context.Context, id, userId bson.ObjectID) (*model.Transaction, error)
	UpdateByIdFn  func(ctx context.Context, trId, userId bson.ObjectID, req repository.TransactionUpdate) (*model.Transaction, error)
	DeleteByIdFn  func(ctx context.Context, trId, userId bson.ObjectID) error
}

func (t *TransactionRepositoryMock) Create(ctx context.Context, tr model.Transaction) (string, error) {
	return t.CreateFn(ctx, tr)
}

func (t *TransactionRepositoryMock) GetByFilter(ctx context.Context, trFilter repository.TransactionFilter) (*repository.PaginatedTransactions, error) {
	return t.GetByFilterFn(ctx, trFilter)
}

func (t *TransactionRepositoryMock) GetById(ctx context.Context, id, userId bson.ObjectID) (*model.Transaction, error) {
	return t.GetByIdFn(ctx, id, userId)
}

func (t *TransactionRepositoryMock) UpdateById(ctx context.Context, trId, userId bson.ObjectID, req repository.TransactionUpdate) (*model.Transaction, error) {
	return t.UpdateByIdFn(ctx, userId, trId, req)
}

func (t *TransactionRepositoryMock) DeleteById(ctx context.Context, trId, userId bson.ObjectID) error {
	return t.DeleteByIdFn(ctx, userId, trId)
}

type CategoryRepositoryMock struct {
	CreateFn              func(ctx context.Context, ct model.Category) (string, error)
	GetByNormalizedNameFn func(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error)
	GetByFilterFn         func(ctx context.Context, cf repository.CategoryFilter) (*repository.PaginatedCategories, error)
	UpdateByIdFn          func(ctx context.Context, id, userId bson.ObjectID, req repository.CategoryUpdate) (*model.Category, error)
	GetByIdFn             func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error)
	DeleteByIdFn          func(ctx context.Context, id, userId bson.ObjectID) error
}

func (c CategoryRepositoryMock) Create(ctx context.Context, ct model.Category) (string, error) {
	return c.CreateFn(ctx, ct)
}

func (c CategoryRepositoryMock) GetByNormalizedName(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
	return c.GetByNormalizedNameFn(ctx, userId, name)
}

func (c CategoryRepositoryMock) GetByFilter(ctx context.Context, cf repository.CategoryFilter) (*repository.PaginatedCategories, error) {
	return c.GetByFilterFn(ctx, cf)
}

func (c CategoryRepositoryMock) UpdateById(ctx context.Context, id, userId bson.ObjectID, req repository.CategoryUpdate) (*model.Category, error) {
	return c.UpdateByIdFn(ctx, id, userId, req)
}

func (c CategoryRepositoryMock) GetById(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
	return c.GetByIdFn(ctx, id, userId)
}

func (c CategoryRepositoryMock) DeleteById(ctx context.Context, id, userId bson.ObjectID) error {
	return c.DeleteByIdFn(ctx, id, userId)
}

func TestTransactionService_Create_Success(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return &model.Category{
				ID:     categoryId,
				UserID: userId,
				Name:   "Alimentação",
			}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		CreateFn: func(ctx context.Context, tr model.Transaction) (string, error) {
			return bson.NewObjectID().Hex(), nil
		},
	}

	svc := NewTransactionService(trRepo, ctRepo)

	req := dto.CreateTransactionRequest{
		Title:           "Restaurante",
		Description:     "Almoço na empresa",
		Amount:          50000,
		Type:            "expense",
		CategoryID:      categoryId.Hex(),
		TransactionDate: time.Now().UTC(),
	}

	res, err := svc.Create(ctx, userId, req)

	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, req.Title, res.Title)
	require.Equal(t, req.Amount, res.Amount)
	require.Equal(t, req.Type, res.Type)
	require.Equal(t, req.TransactionDate, res.TransactionDate)
}

func TestTransactionService_Create_Category_Not_Found(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return nil, nil
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewTransactionService(trRepo, ctRepo)

	req := dto.CreateTransactionRequest{
		Title:           "Restaurante",
		Description:     "Almoço na empresa",
		Amount:          50000,
		Type:            "expense",
		CategoryID:      categoryId.Hex(),
		TransactionDate: time.Now().UTC(),
	}

	_, err := svc.Create(ctx, userId, req)

	require.Error(t, err)
	require.Equal(t, "category not found", err.(apperror.AppError).Error())
}

func TestTransactionService_Create_Category_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return nil, errors.New("")
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewTransactionService(trRepo, ctRepo)

	req := dto.CreateTransactionRequest{
		Title:           "Restaurante",
		Description:     "Almoço na empresa",
		Amount:          50000,
		Type:            "expense",
		CategoryID:      categoryId.Hex(),
		TransactionDate: time.Now().UTC(),
	}

	_, err := svc.Create(ctx, userId, req)

	require.Error(t, err)
}

func TestTransactionService_Create_Transaction_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
			return &model.Category{
				ID:     categoryId,
				UserID: userId,
				Name:   "Alimentação",
			}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		CreateFn: func(ctx context.Context, tr model.Transaction) (string, error) {
			return "", errors.New("")
		},
	}

	svc := NewTransactionService(trRepo, ctRepo)

	req := dto.CreateTransactionRequest{
		Title:           "Restaurante",
		Description:     "Almoço na empresa",
		Amount:          50000,
		Type:            "expense",
		CategoryID:      categoryId.Hex(),
		TransactionDate: time.Now().UTC(),
	}

	_, err := svc.Create(ctx, userId, req)

	require.Error(t, err)
}

func TestTransactionService_Update_Success(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	transactionId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, usrId bson.ObjectID) (*model.Category, error) {
			return &model.Category{
				ID:     categoryId,
				UserID: usrId,
				Name:   "Alimentação",
			}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		UpdateByIdFn: func(ctx context.Context, trId, usrId bson.ObjectID, req repository.TransactionUpdate) (*model.Transaction, error) {
			return &model.Transaction{
				ID:              trId,
				UserID:          usrId,
				CategoryID:      categoryId,
				Title:           *req.Title,
				Amount:          *req.Amount,
				TransactionDate: *req.TransactionDate,
			}, nil
		},
	}

	svc := NewTransactionService(trRepo, ctRepo)

	ctId := categoryId.Hex()
	title := "Feijão com pimenta"
	amount := 51000
	date := time.Now().UTC()
	req := dto.UpdateTransactionRequest{
		CategoryID:      &ctId,
		Title:           &title,
		Amount:          &amount,
		TransactionDate: &date,
	}

	res, err := svc.UpdateById(ctx, userId, transactionId.Hex(), req)

	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, *req.Title, res.Title)
	require.Equal(t, *req.Amount, res.Amount)
	require.Equal(t, *req.TransactionDate, res.TransactionDate)
}

func TestTransactionService_Update_Category_Not_Found(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	transactionId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, usrId bson.ObjectID) (*model.Category, error) {
			return nil, nil
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewTransactionService(trRepo, ctRepo)

	ctId := categoryId.Hex()
	title := "Feijão com pimenta"
	amount := 51000
	date := time.Now().UTC()
	req := dto.UpdateTransactionRequest{
		CategoryID:      &ctId,
		Title:           &title,
		Amount:          &amount,
		TransactionDate: &date,
	}

	_, err := svc.UpdateById(ctx, userId, transactionId.Hex(), req)

	require.Error(t, err)
	require.Equal(t, "category not found", err.(apperror.AppError).Error())
}

func TestTransactionService_Update_Transaction_Not_Found(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	transactionId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, usrId bson.ObjectID) (*model.Category, error) {
			return &model.Category{
				ID:     categoryId,
				UserID: usrId,
				Name:   "Alimentação",
			}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		UpdateByIdFn: func(ctx context.Context, trId, usrId bson.ObjectID, req repository.TransactionUpdate) (*model.Transaction, error) {
			return nil, nil
		},
	}

	svc := NewTransactionService(trRepo, ctRepo)

	ctId := categoryId.Hex()
	title := "Feijão com pimenta"
	amount := 51000
	date := time.Now().UTC()
	req := dto.UpdateTransactionRequest{
		CategoryID:      &ctId,
		Title:           &title,
		Amount:          &amount,
		TransactionDate: &date,
	}

	res, err := svc.UpdateById(ctx, userId, transactionId.Hex(), req)

	require.NoError(t, err)
	require.Nil(t, res)
}

func TestTransactionService_Update_Category_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	transactionId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, usrId bson.ObjectID) (*model.Category, error) {
			return nil, errors.New("")
		},
	}

	trRepo := &TransactionRepositoryMock{}

	svc := NewTransactionService(trRepo, ctRepo)

	ctId := categoryId.Hex()
	title := "Feijão com pimenta"
	amount := 51000
	date := time.Now().UTC()
	req := dto.UpdateTransactionRequest{
		CategoryID:      &ctId,
		Title:           &title,
		Amount:          &amount,
		TransactionDate: &date,
	}

	_, err := svc.UpdateById(ctx, userId, transactionId.Hex(), req)

	require.Error(t, err)
}

func TestTransactionService_Update_Transaction_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	transactionId := bson.NewObjectID()
	categoryId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{
		GetByIdFn: func(ctx context.Context, id, usrId bson.ObjectID) (*model.Category, error) {
			return &model.Category{
				ID:     categoryId,
				UserID: usrId,
				Name:   "Alimentação",
			}, nil
		},
	}

	trRepo := &TransactionRepositoryMock{
		UpdateByIdFn: func(ctx context.Context, trId, usrId bson.ObjectID, req repository.TransactionUpdate) (*model.Transaction, error) {
			return nil, errors.New("")
		},
	}

	svc := NewTransactionService(trRepo, ctRepo)

	ctId := categoryId.Hex()
	title := "Feijão com pimenta"
	amount := 51000
	date := time.Now().UTC()
	req := dto.UpdateTransactionRequest{
		CategoryID:      &ctId,
		Title:           &title,
		Amount:          &amount,
		TransactionDate: &date,
	}

	_, err := svc.UpdateById(ctx, userId, transactionId.Hex(), req)

	require.Error(t, err)
}

func TestTransactionService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	transactionId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{}

	trRepo := &TransactionRepositoryMock{
		DeleteByIdFn: func(ctx context.Context, trId, userId bson.ObjectID) error {
			return nil
		},
	}

	svc := NewTransactionService(trRepo, ctRepo)

	err := svc.DeleteById(ctx, userId, transactionId.Hex())

	require.NoError(t, err)
}

func TestTransactionService_Delete_Repository_Error(t *testing.T) {
	ctx := context.Background()
	userId := bson.NewObjectID()
	transactionId := bson.NewObjectID()

	ctRepo := &CategoryRepositoryMock{}

	trRepo := &TransactionRepositoryMock{
		DeleteByIdFn: func(ctx context.Context, trId, userId bson.ObjectID) error {
			return errors.New("")
		},
	}

	svc := NewTransactionService(trRepo, ctRepo)

	err := svc.DeleteById(ctx, userId, transactionId.Hex())

	require.Error(t, err)
}
