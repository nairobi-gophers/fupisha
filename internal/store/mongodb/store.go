package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nairobi-gophers/fupisha/internal/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

//Store is a mongodb implementation of store interface.
type Store struct {
	db        *mongo.Database
	userStore *userStore
	urlStore  *urlStore
}

//Users returns a user store
func (s *Store) Users() store.UserStore {
	//create a unique index on user email field.
	if _, err := s.db.Collection("users").Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)}); err != nil {
		log.Fatalf("Users: failed to create unique index: %s", err)
	}

	return s.userStore
}

//Urls returns a url store
func (s *Store) Urls() store.URLStore {
	return s.urlStore
}

var _ store.Store = (*Store)(nil)

//Connect connects to a mongodb store and returns an initialized mongo store object.
//address: localhost:27017
func Connect(address, username, password, database string) (*Store, error) {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s&connectTimeoutMS=300000", username, password, address, database)

	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(database)

	s := Store{
		db:        db,
		userStore: &userStore{db: db},
		urlStore:  &urlStore{db: db},
	}

	return &s, nil
}
