package repository

import "time"

type TransactionUpdate struct {
	Title           *string
	Amount          *int
	TransactionDate *time.Time
}
