package repository

import "time"

type TransactionUpdate struct {
	Title           *string
	Amount          *int
	TransactionDate *time.Time
}

type CategoryUpdate struct {
	Name           *string
	NormalizedName string
}
