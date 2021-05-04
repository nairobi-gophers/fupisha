package mock

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/store"
)

//URLStore is a mock implementation of store.urlStore
type URLStore struct {
	OnNew        func(ctx context.Context, userID uuid.UUID, originalURL, shortenedURL string) (store.URL, error)
	OnGet        func(ctx context.Context, id uuid.UUID) (store.URL, error)
	OnGetByParam func(ctx context.Context, param string) (store.URL, error)
}

func (u URLStore) New(ctx context.Context, userID uuid.UUID, originalURL, shortenedURL string) (store.URL, error) {
	return u.OnNew(ctx, userID, originalURL, shortenedURL)
}

func (u URLStore) Get(ctx context.Context, id uuid.UUID) (store.URL, error) {
	return u.OnGet(ctx, id)
}

func (u URLStore) GetByParam(ctx context.Context, param string) (store.URL, error) {
	return u.OnGetByParam(ctx, param)
}
