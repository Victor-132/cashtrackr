package repository

import (
	"context"
	"log"
	"sort"

	"github.com/Victor-132/cashtrackr/internal/reports/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Transaction struct {
	coll *mongo.Collection
}

func NewTransactionRepository(coll *mongo.Collection) Transaction {
	return Transaction{coll}
}

func (t *Transaction) GetMonthlySummary(ctx context.Context, userId bson.ObjectID, year, month int) (*dto.MonthlySummary, error) {
	match := bson.D{
		{Key: "$match", Value: bson.M{
			"user_id": userId,
			"$expr": bson.M{
				"$and": []bson.M{
					{
						"$eq": bson.A{
							bson.M{"$year": "$transaction_date"}, year,
						},
					},
					{
						"$eq": bson.A{
							bson.M{"$month": "$transaction_date"}, month,
						},
					},
				},
			}},
		},
	}

	group := bson.D{
		{Key: "$group", Value: bson.M{
			"_id": "$type",
			"total": bson.M{
				"$sum": "$amount",
			},
		}},
	}

	cur, err := t.coll.Aggregate(ctx, mongo.Pipeline{match, group})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer cur.Close(ctx)

	res := dto.MonthlySummary{}

	for cur.Next(ctx) {
		var row struct {
			ID    string `bson:"_id"`
			Total int    `bson:"total"`
		}

		if err := cur.Decode(&row); err != nil {
			return nil, err
		}

		switch row.ID {
		case "income":
			res.Income += row.Total
		case "expense":
			res.Expense += row.Total
		}
	}

	res.Balance = res.Income - res.Expense

	return &res, nil
}

func (t *Transaction) GetExpensesByCategory(ctx context.Context, userId bson.ObjectID, year int) ([]dto.ExpenseByCategory, error) {
	match := bson.D{
		{Key: "$match", Value: bson.M{
			"user_id": userId,
			"type":    "expense",
			"$expr": bson.M{
				"$eq": bson.A{
					bson.M{"$year": "$transaction_date"}, year,
				},
			},
		}},
	}

	lookup := bson.D{
		{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "category_id",
			"foreignField": "_id",
			"as":           "category",
		}},
	}

	unwind := bson.D{
		{Key: "$unwind", Value: bson.M{
			"path": "$category",
		}},
	}

	group := bson.D{
		{Key: "$group", Value: bson.M{
			"_id": "$category._id",
			"category_name": bson.M{
				"$first": "$category.name",
			},
			"total": bson.M{
				"$sum": "$amount",
			},
		}},
	}

	sort := bson.D{
		{Key: "$sort", Value: bson.M{
			"total": -1,
		}},
	}

	cur, err := t.coll.Aggregate(ctx, mongo.Pipeline{match, lookup, unwind, group, sort})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer cur.Close(ctx)

	res := []dto.ExpenseByCategory{}
	if err := cur.All(ctx, &res); err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}

func (t *Transaction) GetMonthlyEvolution(ctx context.Context, userId bson.ObjectID, year int) ([]dto.MonthlyEvolution, error) {
	match := bson.D{
		{Key: "$match", Value: bson.M{
			"user_id": userId,
			"$expr": bson.M{
				"$eq": bson.A{
					bson.M{"$year": "$transaction_date"}, year,
				},
			},
		}},
	}

	group := bson.D{
		{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"month": bson.M{
					"$month": "$transaction_date",
				},
				"type": "$type",
			},
			"total": bson.M{
				"$sum": "$amount",
			},
		}},
	}

	cur, err := t.coll.Aggregate(ctx, mongo.Pipeline{match, group})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer cur.Close(ctx)

	items := map[int]*dto.MonthlyEvolution{}

	for cur.Next(ctx) {
		type Id struct {
			Month int    `bson:"month"`
			Type  string `bson:"type"`
		}

		var row struct {
			Id    Id  `bson:"_id"`
			Total int `bson:"total"`
		}

		if err := cur.Decode(&row); err != nil {
			return nil, err
		}

		if _, ok := items[row.Id.Month]; !ok {
			items[row.Id.Month] = &dto.MonthlyEvolution{Month: row.Id.Month}
		}

		switch row.Id.Type {
		case "income":
			items[row.Id.Month].Income += row.Total
		case "expense":
			items[row.Id.Month].Expense += row.Total
		}

		items[row.Id.Month].Balance = items[row.Id.Month].Income - items[row.Id.Month].Expense
	}

	for i := 1; i <= 12; i++ {
		if _, ok := items[i]; !ok {
			items[i] = &dto.MonthlyEvolution{Month: i}
		}
	}

	var res []dto.MonthlyEvolution

	for _, a := range items {
		res = append(res, *a)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Month < res[j].Month
	})

	return res, nil
}
