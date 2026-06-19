package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
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

func (c *CategoryRepository) GetByFilter(ctx context.Context, cf CategoryFilter) (*PaginatedCategories, error) {
	filter := bson.M{"user_id": cf.UserID}

	skip := (cf.Page - 1) * cf.Limit

	opt := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(cf.Limit)).
		SetSkip(int64(skip))

	cur, err := c.coll.Find(ctx, filter, opt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.Category
	if err := cur.All(ctx, &list); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			return nil, err
		}
	}

	total, err := c.coll.CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res := PaginatedCategories{
		Categories: list,
		TotalItems: int(total),
	}

	return &res, nil
}

func (c *CategoryRepository) UpdateById(ctx context.Context, id, userId bson.ObjectID, req CategoryUpdate) (*model.Category, error) {
	filter := bson.M{
		"_id":     id,
		"user_id": userId,
	}

	set := bson.M{}
	if req.Name != nil {
		set["name"] = *req.Name
		set["normalized_name"] = req.NormalizedName
	}

	if len(set) > 0 {
		set["updated_at"] = time.Now().UTC()
	}

	upd := bson.M{"$set": set}

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ret *model.Category
	if err := c.coll.FindOneAndUpdate(ctx, filter, upd, opt).Decode(&ret); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			return nil, err
		}
	}

	return ret, nil
}

func (c *CategoryRepository) GetById(ctx context.Context, id, userId bson.ObjectID) (*model.Category, error) {
	filter := bson.M{
		"_id":     id,
		"user_id": userId,
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

func (c *CategoryRepository) DeleteById(ctx context.Context, id, userId bson.ObjectID) error {
	filter := bson.M{
		"_id":     id,
		"user_id": userId,
	}

	res, err := c.coll.DeleteOne(ctx, filter)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.DeletedCount == 0 {
		err = apperror.New("category not found")
		log.Println(err)
		return err
	}

	return nil
}
