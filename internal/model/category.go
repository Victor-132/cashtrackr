package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Category struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	UserID         bson.ObjectID `bson:"user_id"`
	Name           string        `bson:"name"`
	NormalizedName string        `bson:"normalized_name"`
	CreatedAt      time.Time     `bson:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at"`
}
