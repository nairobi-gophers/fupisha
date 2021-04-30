package provider

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/store"
)

//GenAPIKey generates an api key for third party applications.
func GenAPIKey(uid uuid.UUID, s store.Store) (string, error) {
	ctx := context.Background()

	key := encoding.GenUniqueID()

	//persist the api key to the database before encoding it.
	err := s.Users().SetAPIKey(ctx, uid, key)

	if err != nil {
		log.Fatalf("failed to persist generated api key: %s", err)
	}
	//encode the api key and return it to the caller.
	return encoding.Encode(key), nil
}
