package mongodb

import (
	"context"
	"time"

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
func (u *urlStore) New(userID, originalURL, shortenedURL string) (interface{}, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	url := model.URL{
		ID:           primitive.NewObjectID(),
		User:         uid,
		OriginalURL:  originalURL,
		ShortenedURL: shortenedURL,
		CreatedAt:    time.Now(),
	}

	result, err := u.db.Collection("urls").InsertOne(u.ctx, url)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

//Get finds a url by id
func (u *urlStore) Get(id string) (model.URL, error) {
	url := model.URL{}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return url, err
	}

	if err := u.db.Collection("urls").FindOne(u.ctx, bson.M{"_id": docID}).Decode(&url); err != nil {
		return url, err
	}

	return url, nil
}
