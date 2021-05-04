package postgres

import (
	"context"
	"reflect"
	"testing"

	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/store"
)

func TestURL(t *testing.T) {
	s, err := newTestDatabase(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	u, err := s.Users().New(ctx, "test_user@test.com", "test_password")

	if err != nil {
		t.Fatalf("failed to create user: %s", err)
	}

	originalURL := "http://highscalability.com/blog/2016/1/25/design-of-a-modern-cache.html"

	param, err := encoding.GenUniqueParam("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 6)

	if err != nil {
		t.Fatalf("failed to generate url param: %s", err)
	}

	url, err := s.Urls().New(ctx, u.ID, originalURL, param)
	if err != nil {
		t.Fatalf("failed to create url: %s", err)
	}

	want := store.URL{
		ID:           url.ID,
		Owner:        u.ID,
		OriginalURL:  originalURL,
		ShortenedURL: param,
		CreatedAt:    url.CreatedAt,
		UpdatedAt:    url.UpdatedAt,
	}

	got, err := s.Urls().GetByParam(ctx, param)
	if err != nil {
		t.Fatalf("failed to retrieve url param: %s", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v\n want %+v\n", got, want)
	}

	url2, err := s.Urls().Get(ctx, url.ID)
	if err != nil {
		t.Fatalf("failed to retrieve url by id: %s", err)
	}

	if !reflect.DeepEqual(url, url2) {
		t.Fatalf("got %+v\n want %+v\n", got, want)
	}
}
