package service

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionRepository(repo repository.TransactionRepository) TransactionService {
	return TransactionService{repo}
}

func (t *TransactionService) Create(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
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
		UpdatedAt:       tr.UpdatedAt,
	}

	return &res, nil
}

func (t *TransactionService) GetByFilter(ctx context.Context, userId bson.ObjectID, req dto.ListTransactionsRequest) (*dto.ListTransactionsResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	filter := repository.TransactionFilter{
		UserID:    userId,
		Page:      page,
		Limit:     limit,
		Type:      req.Type,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	ret, err := t.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	list := []dto.TransactionResponse{}
	for _, tr := range ret.Transactions {
		list = append(list, dto.TransactionResponse{
			ID:              tr.ID.Hex(),
			Title:           tr.Title,
			Description:     tr.Description,
			Amount:          tr.Amount,
			Type:            string(tr.Type),
			TransactionDate: tr.TransactionDate,
			CreatedAt:       tr.CreatedAt,
			UpdatedAt:       tr.UpdatedAt,
		})
	}

	res := dto.ListTransactionsResponse{
		Data:       list,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: ret.TotalItems,
		TotalPages: int(math.Ceil(float64(ret.TotalItems) / float64(filter.Limit))),
	}

	return &res, nil
}

func (t *TransactionService) GetById(ctx context.Context, userId bson.ObjectID, trId string) (*dto.TransactionResponse, error) {
	id, err := bson.ObjectIDFromHex(trId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tr, err := t.repo.GetById(ctx, id, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if tr == nil {
		return nil, nil
	}

	res := dto.TransactionResponse{
		ID:              tr.ID.Hex(),
		Title:           tr.Title,
		Description:     tr.Description,
		Amount:          tr.Amount,
		Type:            string(tr.Type),
		TransactionDate: tr.TransactionDate,
		CreatedAt:       tr.CreatedAt,
		UpdatedAt:       tr.UpdatedAt,
	}

	return &res, nil
}

func (t *TransactionService) UpdateById(ctx context.Context, userId bson.ObjectID, trId string, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	id, err := bson.ObjectIDFromHex(trId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	upd := repository.TransactionUpdate{
		Title:           req.Title,
		Amount:          req.Amount,
		TransactionDate: req.TransactionDate,
	}

	tr, err := t.repo.UpdateById(ctx, id, userId, upd)
	if err != nil {
		return nil, err
	}

	if tr == nil {
		return nil, nil
	}

	ret := dto.TransactionResponse{
		ID:              tr.ID.Hex(),
		Title:           tr.Title,
		Description:     tr.Description,
		Amount:          tr.Amount,
		Type:            string(tr.Type),
		TransactionDate: tr.TransactionDate,
		CreatedAt:       tr.CreatedAt,
		UpdatedAt:       tr.UpdatedAt,
	}

	return &ret, nil
}
