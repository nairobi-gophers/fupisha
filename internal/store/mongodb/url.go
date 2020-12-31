package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

type urlStore struct {
	db *mongo.Database
}

//New creates a new url document.
func (u *urlStore) New(ctx context.Context, userID, originalURL, shortenedURL string) (string, error) {
	//Lets check if its a valid ObjectID
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	url := model.URL{
		ID:           primitive.NewObjectID().Hex(),
		User:         userID,
		OriginalURL:  originalURL,
		ShortenedURL: shortenedURL,
		CreatedAt:    time.Now(),
	}

	result, err := u.db.Collection("urls").InsertOne(ctx, url)
	if err != nil {
		return "", err
	}

	resultID, ok := result.InsertedID.(string)

	if !ok {
		return "", errors.New("invalid object id")
	}

	return resultID, nil
}

//Get finds a url by id
func (u *urlStore) Get(ctx context.Context, id string) (model.URL, error) {
	url := model.URL{}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return url, err
	}

	if err := u.db.Collection("urls").FindOne(ctx, bson.M{"_id": docID}).Decode(&url); err != nil {
		return url, err
	}

	return url, nil
}
