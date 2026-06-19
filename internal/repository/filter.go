package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionFilter struct {
	UserID     bson.ObjectID
	Page       int
	Limit      int
	Type       string
	CategoryID *bson.ObjectID
	StartDate  *time.Time
	EndDate    *time.Time
}

type CategoryFilter struct {
	UserID bson.ObjectID
	Page   int
	Limit  int
}
