package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type userStore struct {
	db  *mongo.Database
	ctx context.Context
}

//New creates a new user document
func (s userStore) New(name, email, password string) (string, error) {

	tkn := encoding.GenUniqueID()

	fmt.Println("TOKEN", tkn)
	id := primitive.NewObjectID()

	var insertedID string //zero value

	user := model.User{
		ID:                  id.Hex(),
		Name:                name,
		Email:               email,
		Password:            password,
		VerificationToken:   tkn,
		VerificationExpires: time.Now().Add(time.Minute * 60), //60 mins
		CreatedAt:           time.Now(),
	}

	if err := user.HashPassword(); err != nil {
		return insertedID, err
	}

	result, err := s.db.Collection("users").InsertOne(s.ctx, user)
	if err != nil {
		return insertedID, err
	}

	resultID, ok := result.InsertedID.(string)

	if !ok {
		return insertedID, errors.New("invalid object id")
	}

	insertedID = resultID

	return insertedID, nil
}

//Get finds a user by id
func (s userStore) Get(id string) (model.User, error) {

	user := model.User{}

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	if err := s.db.Collection("users").FindOne(s.ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

//GetByEmail retrieve an existing user with the given email
func (s userStore) GetByEmail(email string) (model.User, error) {
	user := model.User{}

	if err := s.db.Collection("users").FindOne(s.ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

//SetAPIKey sets the api key for the given user id.
func (s userStore) SetAPIKey(id string, key uuid.UUID) error {

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{"apiKey": key},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := s.db.Collection("users").FindOneAndUpdate(s.ctx, filter, update, &opt)

	if result.Err() != nil {
		return result.Err()
	}

	// doc := model.User{}
	// decodeErr := result.Decode(&doc)

	// return doc, decodeErr
	return nil
}
