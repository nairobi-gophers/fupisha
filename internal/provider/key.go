package provider

import (
	"log"

	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/store"
)

//GenAPIKey generates an api key for third party applications.
func GenAPIKey(uid string, s store.Store) (string, error) {
	key := encoding.GenUniqueID()

	//persist the api key to the database before encoding it.
	_, err := s.Users().SetAPIKey(uid, key)

	if err != nil {
		log.Fatalf("failed to persist generated api key: %s", err)
	}
	//encode the api key and return it to the caller.
	return encoding.Encode(key), nil
}
