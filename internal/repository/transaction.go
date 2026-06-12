package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Victor-132/cashtrackr/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TransactionRepository struct {
	coll *mongo.Collection
}

func NewTransactionRepository(coll *mongo.Collection) TransactionRepository {
	return TransactionRepository{coll}
}

func (t *TransactionRepository) Create(ctx context.Context, tr model.Transaction) (string, error) {
	res, err := t.coll.InsertOne(ctx, tr)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (t *TransactionRepository) GetByUserId(ctx context.Context, trFilter TransactionFilter) (*PaginatedTransactions, error) {
	filter := bson.M{"user_id": trFilter.UserID}

	if strings.TrimSpace(trFilter.Type) != "" {
		filter["type"] = trFilter.Type
	}

	if trFilter.StartDate != nil || trFilter.EndDate != nil {
		dateFilter := bson.M{}

		if trFilter.StartDate != nil {
			dateFilter["$gte"] = *trFilter.StartDate
		}

		if trFilter.EndDate != nil {
			dateFilter["$lt"] = *trFilter.EndDate
		}

		filter["transaction_date"] = dateFilter
	}

	skip := (trFilter.Page - 1) * trFilter.Limit

	opt := options.Find().
		SetSort(bson.M{"transaction_date": -1}).
		SetLimit(int64(trFilter.Limit)).
		SetSkip(int64(skip))

	cur, err := t.coll.Find(ctx, filter, opt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.Transaction
	if err := cur.All(ctx, &list); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			return nil, err
		}
	}

	total, err := t.coll.CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res := PaginatedTransactions{
		Transactions: list,
		TotalItems:   int(total),
	}

	return &res, nil
}
