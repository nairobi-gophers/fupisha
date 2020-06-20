package mongodb

import (
	"context"

	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

type urlStore struct {
	db  *mongo.Database
	ctx context.Context
}

//New creates a new url document.
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

//Get finds a url by id
func (u *urlStore) Get(id string) (model.URL, error) {
	url := model.URL{}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return url, err
	}

	if err := u.db.Collection("urls").FindOne(u.ctx, bson.M{"_id": docID}).Decode(url); err != nil {
		return url, err
	}

	return url, nil
}
