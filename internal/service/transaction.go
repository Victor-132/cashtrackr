package service

import (
	"context"
	"log"
	"math"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/model"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionService interface {
	Create(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetByFilter(ctx context.Context, userId bson.ObjectID, req dto.ListTransactionsRequest) (*dto.ListTransactionsResponse, error)
	GetById(ctx context.Context, userId bson.ObjectID, trId string) (*dto.TransactionResponse, error)
	UpdateById(ctx context.Context, userId bson.ObjectID, trId string, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error)
	DeleteById(ctx context.Context, userId bson.ObjectID, trId string) error
}

type Transaction struct {
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository, categoryRepo repository.CategoryRepository) TransactionService {
	return &Transaction{transactionRepo, categoryRepo}
}

func (t *Transaction) Create(ctx context.Context, userId bson.ObjectID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	ctID, err := bson.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ct, err := t.categoryRepo.GetById(ctx, ctID, userId)
	if err != nil {
		return nil, err
	}

	if ct == nil {
		err = apperror.New("category not found")
		log.Println(err)
		return nil, err
	}

	tr := model.Transaction{
		UserID:          userId,
		CategoryID:      ct.ID,
		Title:           req.Title,
		Description:     req.Description,
		Amount:          req.Amount,
		Type:            model.TransactionType(req.Type),
		TransactionDate: req.TransactionDate,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	trId, err := t.transactionRepo.Create(ctx, tr)
	if err != nil {
		return nil, err
	}

	res := dto.TransactionResponse{
		ID:              trId,
		Title:           tr.Title,
		Description:     tr.Description,
		Amount:          tr.Amount,
		Type:            string(tr.Type),
		Category:        ct.Name,
		TransactionDate: tr.TransactionDate,
		CreatedAt:       tr.CreatedAt,
		UpdatedAt:       tr.UpdatedAt,
	}

	return &res, nil
}

func (t *Transaction) GetByFilter(ctx context.Context, userId bson.ObjectID, req dto.ListTransactionsRequest) (*dto.ListTransactionsResponse, error) {
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

	ret, err := t.transactionRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	list := []dto.TransactionResponse{}
	for _, tr := range ret.Transactions {
		ct, err := t.categoryRepo.GetById(ctx, tr.CategoryID, tr.UserID)
		if err != nil {
			return nil, err
		}

		list = append(list, dto.TransactionResponse{
			ID:              tr.ID.Hex(),
			Title:           tr.Title,
			Description:     tr.Description,
			Amount:          tr.Amount,
			Type:            string(tr.Type),
			Category:        ct.Name,
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

func (t *Transaction) GetById(ctx context.Context, userId bson.ObjectID, trId string) (*dto.TransactionResponse, error) {
	id, err := bson.ObjectIDFromHex(trId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tr, err := t.transactionRepo.GetById(ctx, id, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if tr == nil {
		return nil, nil
	}

	ct, err := t.categoryRepo.GetById(ctx, tr.CategoryID, tr.UserID)
	if err != nil {
		return nil, err
	}

	res := dto.TransactionResponse{
		ID:              tr.ID.Hex(),
		Title:           tr.Title,
		Description:     tr.Description,
		Amount:          tr.Amount,
		Type:            string(tr.Type),
		Category:        ct.Name,
		TransactionDate: tr.TransactionDate,
		CreatedAt:       tr.CreatedAt,
		UpdatedAt:       tr.UpdatedAt,
	}

	return &res, nil
}

func (t *Transaction) UpdateById(ctx context.Context, userId bson.ObjectID, trId string, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	var ct *model.Category
	if req.CategoryID != nil {
		categoryId, err := bson.ObjectIDFromHex(*req.CategoryID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		category, err := t.categoryRepo.GetById(ctx, categoryId, userId)
		if err != nil {
			return nil, err
		}

		if category == nil {
			err = apperror.New("category not found")
			log.Println(err)
			return nil, err
		}

		ct = category
	}

	id, err := bson.ObjectIDFromHex(trId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	upd := repository.TransactionUpdate{
		CategoryID:      &ct.ID,
		Title:           req.Title,
		Amount:          req.Amount,
		TransactionDate: req.TransactionDate,
	}

	tr, err := t.transactionRepo.UpdateById(ctx, id, userId, upd)
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
		Category:        ct.Name,
		TransactionDate: tr.TransactionDate,
		CreatedAt:       tr.CreatedAt,
		UpdatedAt:       tr.UpdatedAt,
	}

	return &ret, nil
}

func (t *Transaction) DeleteById(ctx context.Context, userId bson.ObjectID, trId string) error {
	id, err := bson.ObjectIDFromHex(trId)
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.transactionRepo.DeleteById(ctx, id, userId)

	return err
}
