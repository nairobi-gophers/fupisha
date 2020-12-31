package store

import (
	"context"

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
	New(ctx context.Context, name, email, password string) (string, error)
	Get(ctx context.Context, id string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	SetAPIKey(ctx context.Context, id string, key uuid.UUID) error
}

//URLStore is a url data store interface.
type URLStore interface {
	New(ctx context.Context, userID, originalURL, shortenedURL string) (string, error)
	Get(ctx context.Context, id string) (model.URL, error)
}
