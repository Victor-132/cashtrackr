package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID              bson.ObjectID   `bson:"_id,omitempty"`
	UserID          bson.ObjectID   `bson:"user_id"`
	Title           string          `bson:"title"`
	Description     string          `bson:"description,omitempty"`
	Amount          int             `bson:"amount"`
	Type            TransactionType `bson:"type"`
	TransactionDate time.Time       `bson:"transaction_date"`
	CreatedAt       time.Time       `bson:"created_at"`
	UpdatedAt       time.Time       `bson:"updated_at"`
}
