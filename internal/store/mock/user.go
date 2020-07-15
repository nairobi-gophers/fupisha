package mock

import (
	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserStore is a mock implementation of store.UserStore.
type UserStore struct {
	OnNew        func(name, email, password string) (primitive.ObjectID, error)
	OnGet        func(id string) (model.User, error)
	OnGetByEmail func(email string) (model.User, error)
	OnSetAPIKey  func(id string, key uuid.UUID) (model.User, error)
}

func (s UserStore) New(name, email, password string) (primitive.ObjectID, error) {
	return s.OnNew(name, email, password)
}

func (s UserStore) Get(id string) (model.User, error) {
	return s.OnGet(id)
}

func (s UserStore) GetByEmail(email string) (model.User, error) {
	return s.OnGetByEmail(email)
}

func (s UserStore) SetAPIKey(id string, key uuid.UUID) (model.User, error) {
	return s.OnSetAPIKey(id, key)
}
