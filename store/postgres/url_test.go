package postgres

import (
	"context"
	"testing"
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
	shortURL := "https://fp.org/l2g4th9Urety2i"

	_, err = s.Urls().New(ctx, u.ID, originalURL, shortURL)
	if err != nil {
		t.Fatalf("failed to create url: %s", err)
	}

}
