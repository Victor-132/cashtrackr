package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ConnectMongo(dbName string) (*mongo.Database, error) {
	c, err := mongo.Connect()
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return c.Database(dbName), nil
}
