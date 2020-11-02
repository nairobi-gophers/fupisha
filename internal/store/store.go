package store

import (
	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
)

//Store is a data store interface.
type Store interface {
	Users() UserStore
	Urls() URLStore
}

//UserStore is a user data store interface.
type UserStore interface {
	New(name, email, password string) (string, error)
	Get(id string) (model.User, error)
	GetByEmail(email string) (model.User, error)
	SetAPIKey(id string, key uuid.UUID) error
}

//URLStore is a url data store interface.
type URLStore interface {
	New(userID, originalURL, shortenedURL string) (string, error)
	Get(id string) (model.URL, error)
}
