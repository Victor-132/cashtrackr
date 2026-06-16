package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateIndexes(db *mongo.Database) error {
	users := db.Collection("users")

	_, err := users.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("idx_users_email_unique"),
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	transactions := db.Collection("transactions")

	_, err = transactions.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "transaction_date", Value: -1},
			},
			Options: options.Index().
				SetName("idx_transactions_user_date"),
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	categories := db.Collection("categories")

	_, err = categories.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "normalized_name", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("idx_categories_user_normalized_name_unique"),
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
