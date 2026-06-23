package service

import (
	"context"

	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/repository"
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
