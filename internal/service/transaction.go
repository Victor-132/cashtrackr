package service

import (
	"context"
	"time"

	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TrasnactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionRepository(repo repository.TransactionRepository) TrasnactionService {
	return TrasnactionService{repo}
}

func (t *TrasnactionService) Create(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	tr := model.Transaction{
		UserID:          userId,
		Title:           req.Title,
		Description:     req.Description,
		Amount:          req.Amount,
		Type:            model.TransactionType(req.Type),
		TransactionDate: req.TransactionDate,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	trId, err := t.repo.Create(ctx, tr)
	if err != nil {
		return nil, err
	}

	res := dto.TransactionResponse{
		ID:              trId,
		Title:           tr.Title,
		Description:     tr.Description,
		Amount:          tr.Amount,
		Type:            string(tr.Type),
		TransactionDate: tr.TransactionDate,
		CreatedAt:       tr.CreatedAt,
	}

	return &res, nil
}
