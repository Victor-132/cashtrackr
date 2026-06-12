package repository

import "github.com/Victor-132/cashtrackr/internal/model"

type PaginatedTransactions struct {
	Transactions []model.Transaction
	TotalItems   int
}
