package auth

import (
	"context"

	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/provider"
	"github.com/nairobi-gophers/fupisha/store"
)

//Resource defines dependencies for auth handlers.
type Resource struct {
	Store  store.Store
	Config *config.Config
	Mailer *provider.Mailer
}

//NewResource returns a configured authentication resource.
func NewResource(store store.Store, cfg *config.Config, mailer *provider.Mailer) *Resource {
	return &Resource{
		Store:  store,
		Config: cfg,
		Mailer: mailer,
	}
}

//FromContext extracts user id from a Context
func FromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
