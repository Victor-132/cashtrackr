package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Victor-132/cashtrackr/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CategoryRepository struct {
	coll *mongo.Collection
}

func NewCategoryRepository(coll *mongo.Collection) CategoryRepository {
	return CategoryRepository{coll}
}

func (c *CategoryRepository) Create(ctx context.Context, ct model.Category) (string, error) {
	res, err := c.coll.InsertOne(ctx, ct)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (c *CategoryRepository) GetByNormalizedName(ctx context.Context, userId bson.ObjectID, name string) (*model.Category, error) {
	filter := bson.M{
		"user_id":         userId,
		"normalized_name": name,
	}

	var ret *model.Category
	if err := c.coll.FindOne(ctx, filter).Decode(&ret); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			return nil, err
		}
	}

	return ret, nil
}
