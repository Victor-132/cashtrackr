package service

import (
	"context"

	"github.com/Victor-132/cashtrackr/internal/reports/dto"
	"github.com/Victor-132/cashtrackr/internal/reports/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Transaction struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) Transaction {
	return Transaction{repo}
}

func (t *Transaction) GetMonthlySummary(ctx context.Context, userId bson.ObjectID, year, month int) (*dto.MonthlySummary, error) {
	return t.repo.GetMonthlySummary(ctx, userId, year, month)
}

func (t *Transaction) GetExpensesByCategory(ctx context.Context, userId bson.ObjectID, year int) ([]dto.ExpenseByCategory, error) {
	return t.repo.GetExpensesByCategory(ctx, userId, year)
}

func (t *Transaction) GetMonthlyEvolution(ctx context.Context, userId bson.ObjectID, year int) ([]dto.MonthlyEvolution, error) {
	return t.repo.GetMonthlyEvolution(ctx, userId, year)
}
