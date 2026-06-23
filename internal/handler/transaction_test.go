package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/go-openapi/testify/v2/require"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionServiceMock struct {
	CreateFn      func(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetByFilterFn func(ctx context.Context, userId bson.ObjectID, req dto.ListTransactionsRequest) (*dto.ListTransactionsResponse, error)
	GetByIdFn     func(ctx context.Context, userId bson.ObjectID, trId string) (*dto.TransactionResponse, error)
	UpdateByIdFn  func(ctx context.Context, userId bson.ObjectID, trId string, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error)
	DeleteByIdFn  func(ctx context.Context, userId bson.ObjectID, trId string) error
}

func (t *TransactionServiceMock) Create(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	return t.CreateFn(ctx, userId, req)
}

func (t *TransactionServiceMock) GetByFilter(ctx context.Context, userId bson.ObjectID, req dto.ListTransactionsRequest) (*dto.ListTransactionsResponse, error) {
	return t.GetByFilterFn(ctx, userId, req)
}

func (t *TransactionServiceMock) GetById(ctx context.Context, userId bson.ObjectID, trId string) (*dto.TransactionResponse, error) {
	return t.GetByIdFn(ctx, userId, trId)
}

func (t *TransactionServiceMock) UpdateById(ctx context.Context, userId bson.ObjectID, trId string, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	return t.UpdateByIdFn(ctx, userId, trId, req)
}

func (t *TransactionServiceMock) DeleteById(ctx context.Context, userId bson.ObjectID, trId string) error {
	return t.DeleteByIdFn(ctx, userId, trId)
}

func TestTransactionHandler_Create_Success(t *testing.T) {
	app := fiber.New()

	mockService := &TransactionServiceMock{
		CreateFn: func(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
			return &dto.TransactionResponse{
				Title:  "Mercado",
				Amount: 25075,
			}, nil
		},
	}

	handler := NewTransactionHandler(mockService)

	app.Post("/transactions", func(c *fiber.Ctx) error {
		c.Locals("user_id", bson.NewObjectID())
		return handler.Create(c)
	})

	body := `{
		"title": "Mercado",
		"amount": 25075,
		"type": "expense",
		"category_id": "abc",
		"transaction_date": "2026-06-01T00:00:00.000Z"
	}`

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result dto.TransactionResponse
	json.NewDecoder(resp.Body).Decode(&result)

	require.Equal(t, "Mercado", result.Title)
	require.Equal(t, 25075, result.Amount)
}

func TestTransactionHandler_Create_Invalid_Title(t *testing.T) {
	app := fiber.New()

	mockService := &TransactionServiceMock{}

	handler := NewTransactionHandler(mockService)

	app.Post("/transactions", func(c *fiber.Ctx) error {
		c.Locals("user_id", bson.NewObjectID())
		return handler.Create(c)
	})

	body := `{
		"title": "",
		"amount": 25075,
		"type": "expense",
		"category_id": "abc",
		"transaction_date": "2026-06-01T00:00:00.000Z"
	}`

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	bodyString := string(bodyBytes)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, "invalid title", bodyString)
}

func TestTransactionHandler_Create_Business_Rule_Error(t *testing.T) {
	app := fiber.New()

	mockService := &TransactionServiceMock{
		CreateFn: func(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
			return nil, apperror.New("")
		},
	}

	handler := NewTransactionHandler(mockService)

	app.Post("/transactions", func(c *fiber.Ctx) error {
		c.Locals("user_id", bson.NewObjectID())
		return handler.Create(c)
	})

	body := `{
		"title": "Mercado",
		"amount": 25075,
		"type": "expense",
		"category_id": "abc",
		"transaction_date": "2026-06-01T00:00:00.000Z"
	}`

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTransactionHandler_Create_Internal_Error(t *testing.T) {
	app := fiber.New()

	mockService := &TransactionServiceMock{
		CreateFn: func(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
			return nil, errors.New("")
		},
	}

	handler := NewTransactionHandler(mockService)

	app.Post("/transactions", func(c *fiber.Ctx) error {
		c.Locals("user_id", bson.NewObjectID())
		return handler.Create(c)
	})

	body := `{
		"title": "Mercado",
		"amount": 25075,
		"type": "expense",
		"category_id": "abc",
		"transaction_date": "2026-06-01T00:00:00.000Z"
	}`

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
