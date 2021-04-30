package mock

import "github.com/nairobi-gophers/fupisha/internal/store"

// Store is a mock implementation of store.Store.
type Store struct {
	UserStore UserStore
	URLStore  URLStore
}

func (s Store) Users() store.UserStore {
	return s.UserStore
}

func (s Store) Urls() store.URLStore {
	return s.URLStore
}
