package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionUpdate struct {
	CategoryID      *bson.ObjectID
	Title           *string
	Amount          *int
	TransactionDate *time.Time
}

type CategoryUpdate struct {
	Name           *string
	NormalizedName string
}
