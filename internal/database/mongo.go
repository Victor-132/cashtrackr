package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongo() (*mongo.Database, error) {
	var opts *options.ClientOptions

	uri := os.Getenv("MONGODB_URI")
	if uri != "" {
		opts = options.Client().
			ApplyURI(uri)
	}

	c, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "cashtrackr_dev"
	}

	return c.Database(dbName), nil
}
