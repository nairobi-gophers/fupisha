package store

import (
	"context"

	"github.com/gofrs/uuid"
)

//Store is a data store interface.
type Store interface {
	Users() UserStore
	Urls() URLStore
}

//UserStore is a user data store interface.
type UserStore interface {
	New(ctx context.Context, email, password string) (User, error)
	Get(ctx context.Context, id uuid.UUID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByVerificationToken(ctx context.Context, token uuid.UUID) (User, error)
	SetAPIKey(ctx context.Context, id, key uuid.UUID) error
	SetVerified(ctx context.Context, id uuid.UUID) error
}

//URLStore is a url data store interface.
type URLStore interface {
	New(ctx context.Context, userID uuid.UUID, originalURL, shortenedURL string) (URL, error)
	Get(ctx context.Context, id uuid.UUID) (URL, error)
	GetByParam(ctx context.Context, param string) (URL, error)
	GetByURL(ctx context.Context, longURL string) (URL, error)
}
