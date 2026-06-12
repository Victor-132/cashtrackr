package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Victor-132/cashtrackr/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(coll *mongo.Collection) UserRepository {
	return UserRepository{coll}
}

func (u *UserRepository) Create(ctx context.Context, user model.User) (string, error) {
	res, err := u.coll.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (u *UserRepository) Update(ctx context.Context, user model.User) error {
	upd := bson.M{"$set": user}

	_, err := u.coll.UpdateByID(ctx, user.ID, upd)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (u *UserRepository) GetByID(ctx context.Context, userID bson.ObjectID) (*model.User, error) {
	filter := bson.M{"_id": userID}

	var usr model.User

	if err := u.coll.FindOne(ctx, filter).Decode(&usr); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		log.Println(err)
		return nil, err
	}

	return &usr, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	filter := bson.M{"email": email}

	var usr model.User

	if err := u.coll.FindOne(ctx, filter).Decode(&usr); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		log.Println(err)
		return nil, err
	}

	return &usr, nil
}
