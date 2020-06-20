package mongodb

import (
	"context"
	"strconv"

	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

type urlStore struct {
	db  *mongo.Database
	ctx context.Context
}

func (u *urlStore) New(userID, originalURL, shortenedURL, customAlias, target string) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	url := model.URL{
		ID:          primitive.NewObjectID(),
		User:        uid,
		OriginalURL: originalURL,
		CustomAlias: customAlias,
		Target:      target,
	}

	if _, err := u.db.Collection("urls").InsertOne(u.ctx, url); err != nil {
		return err
	}

	return nil
}

func (u *urlStore) Get(id int) (model.URL, error) {
	url := model.URL{}

	docID, err := primitive.ObjectIDFromHex(strconv.Itoa(id))
	if err != nil {
		return url, err
	}

	if err := u.db.Collection("urls").FindOne(u.ctx, bson.M{"_id": docID}).Decode(url); err != nil {
		return url, err
	}

	return url, nil
}
