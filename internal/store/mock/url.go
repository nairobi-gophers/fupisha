package mock

import "github.com/nairobi-gophers/fupisha/internal/store/model"

//URLStore is a mock implementation of store.urlStore
type URLStore struct {
	OnNew func(userID, originalURL, shortenedURL string) (interface{}, error)
	OnGet func(id string) (model.URL, error)
}

func (u URLStore) New(userID, originalURL, shortenedURL string) (interface{}, error) {
	return u.OnNew(userID, originalURL, shortenedURL)
}

func (u URLStore) Get(id string) (model.URL, error) {
	return u.OnGet(id)
}
