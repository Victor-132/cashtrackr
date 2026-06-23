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
