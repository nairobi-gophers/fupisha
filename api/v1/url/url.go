package url

import (
	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/store"
)

//Resource defines dependencies for url handlers.
type Resource struct {
	Store  store.Store
	Config *config.Config
}

//NewResource returns a configures url resource.
func NewResource(store store.Store, cfg *config.Config) *Resource {
	return &Resource{
		Store:  store,
		Config: cfg,
	}
}
