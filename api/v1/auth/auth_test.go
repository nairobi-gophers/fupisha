package auth_test

import (
	"testing"

	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/api/v1/auth"
	"github.com/nairobi-gophers/fupisha/internal/config"
)

func testAuthResource(t *testing.T) *auth.Resource {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	store, err := cfg.GetStore()
	if err != nil {
		t.Fatal(err)
	}

	authResource := auth.NewResource(store, cfg)

	return authResource
}

func testApiServer(t *testing.T) {
	// router := chi.NewRouter()
	api, err := api.NewServer()
	if err != nil {
		t.Fatal(err)
	}

	api.Start()
}
