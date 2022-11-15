package store

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/store"
)

// Store is composes all the different store abstractions into one single abstraction.
type Store interface {
	UserStore
	URLStore
}

// UserStore is a user data store interface.
type UserStore interface {
	NewUser(ctx context.Context, email, password string) (store.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (store.User, error)
	GetUserByEmail(ctx context.Context, email string) (store.User, error)
	GetUserByVerificationToken(ctx context.Context, token uuid.UUID) (store.User, error)
	SetUserAPIKey(ctx context.Context, id, key uuid.UUID) error
	SetUserVerified(ctx context.Context, id uuid.UUID) error
}

// URLStore is a url data store interface.
type URLStore interface {
	NewURL(ctx context.Context, userID uuid.UUID, originalURL, shortenedURL string) (store.URL, error)
	GetURLByID(ctx context.Context, id uuid.UUID) (store.URL, error)
	GetURLByParam(ctx context.Context, param string) (store.URL, error)
	GetURLByLongStr(ctx context.Context, longURL string) (store.URL, error)
}
