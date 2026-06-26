package dto

import "go.mongodb.org/mongo-driver/v2/bson"

type MonthlySummary struct {
	Income  int `json:"income"`
	Expense int `json:"expense"`
	Balance int `json:"balance"`
}

type ExpenseByCategory struct {
	CategoryID   bson.ObjectID `json:"category_id" bson:"_id"`
	CategoryName string        `json:"category_name" bson:"category_name"`
	Total        int64         `json:"total" bson:"total"`
}

type MonthlyEvolution struct {
	Month   int `json:"month"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
	Balance int `json:"balance"`
}
