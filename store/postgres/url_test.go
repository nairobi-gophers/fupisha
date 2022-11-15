package postgres

import (
	"context"
	"reflect"
	"testing"

	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/store"
)

func TestURL(t *testing.T) {
	s, teardown := NewTestDatabase(t)
	t.Cleanup(teardown)

	ctx := context.Background()

	u, err := s.NewUser(ctx, "test_user@test.com", "test_password")

	if err != nil {
		t.Fatalf("failed to create user: %s", err)
	}

	originalURL := "http://highscalability.com/blog/2016/1/25/design-of-a-modern-cache.html"

	param, err := encoding.GenUniqueParam("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 6)

	if err != nil {
		t.Fatalf("failed to generate url param: %s", err)
	}

	url, err := s.NewURL(ctx, u.ID, originalURL, param)
	if err != nil {
		t.Fatalf("failed to create url: %s", err)
	}

	want := store.URL{
		ID:                url.ID,
		Owner:             u.ID,
		OriginalURL:       originalURL,
		ShortenedURLParam: param,
		CreatedAt:         url.CreatedAt,
		UpdatedAt:         url.UpdatedAt,
	}

	got, err := s.GetURLByParam(ctx, param)
	if err != nil {
		t.Fatalf("failed to retrieve url param: %s", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v\n want %+v\n", got, want)
	}

	url2, err := s.GetURLByID(ctx, url.ID)
	if err != nil {
		t.Fatalf("failed to retrieve url by id: %s", err)
	}

	if !reflect.DeepEqual(url, url2) {
		t.Fatalf("got %+v\n want %+v\n", got, want)
	}

	url3, err := s.GetURLByLongStr(ctx, originalURL)
	if err != nil {
		t.Fatal(err)
	}

	if url3.OriginalURL != originalURL {
		t.Fatalf("got %v want %v\n", url3.OriginalURL, originalURL)
	}
}
