package repository

import (
	"context"
	"log"

	"github.com/Victor-132/cashtrackr/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
