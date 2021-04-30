package mock

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/internal/store"
)

// UserStore is a mock implementation of store.UserStore.
type UserStore struct {
	OnNew        func(ctx context.Context, email, password string) (store.User, error)
	OnGet        func(ctx context.Context, id uuid.UUID) (store.User, error)
	OnGetByEmail func(ctx context.Context, email string) (store.User, error)
	OnSetAPIKey  func(ctx context.Context, id, key uuid.UUID) error
}

func (s UserStore) New(ctx context.Context, email, password string) (store.User, error) {
	return s.OnNew(ctx, email, password)
}

func (s UserStore) Get(ctx context.Context, id uuid.UUID) (store.User, error) {
	return s.OnGet(ctx, id)
}

func (s UserStore) GetByEmail(ctx context.Context, email string) (store.User, error) {
	return s.OnGetByEmail(ctx, email)
}

func (s UserStore) SetAPIKey(ctx context.Context, id, key uuid.UUID) error {
	return s.OnSetAPIKey(ctx, id, key)
}
