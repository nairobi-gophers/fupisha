package store

import (
	"context"

	"github.com/gofrs/uuid"
)

// Store is composes all the different store abstractions into one single abstraction.
type Store interface {
	UserStore
	URLStore
}

// UserStore is a user data store interface.
type UserStore interface {
	NewUser(ctx context.Context, email, password string) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByVerificationToken(ctx context.Context, token uuid.UUID) (User, error)
	SetUserAPIKey(ctx context.Context, id, key uuid.UUID) error
	SetUserVerified(ctx context.Context, id uuid.UUID) error
}

// URLStore is a url data store interface.
type URLStore interface {
	NewURL(ctx context.Context, userID uuid.UUID, originalURL, shortenedURL string) (URL, error)
	GetURLByID(ctx context.Context, id uuid.UUID) (URL, error)
	GetURLByParam(ctx context.Context, param string) (URL, error)
	GetURLByLongStr(ctx context.Context, longURL string) (URL, error)
}
