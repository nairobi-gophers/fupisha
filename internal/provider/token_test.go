package provider

import (
	"strings"
	"testing"
	"time"

	"github.com/nairobi-gophers/fupisha/internal/encoding"
)

func TestEncodeDecode(t *testing.T) {
	userID := "5d0575344d9f7ff15e989174"

	secret := encoding.GenHexKey(32)

	s, err := NewJWTService(secret)
	if err != nil {
		t.Fatalf("failed to create a new service: %s", err)
	}

	tokenString, err := s.Encode(userID)
	if err != nil {
		t.Fatalf("failed to create a token: %s", err)
	}

	gotUserID, gotIssuedAt, err := s.Decode(tokenString)
	if err != nil {
		t.Fatalf("failed to verify a token: %s", err)
	}

	if gotUserID != userID {
		t.Fatalf("bad verified userID: got %v; want %v", gotUserID, userID)
	}

	now := time.Now()
	since := now.Sub(gotIssuedAt)
	if since < 0 || since > 3*time.Second {
		t.Fatalf("bad issuedAt of the verified token: %v; now: %v", gotIssuedAt, now)
	}

	badSecret := strings.Repeat("1", 64)

	s, err = NewJWTService(badSecret)
	if err != nil {
		t.Fatalf("failed to create a new service: %s", err)
	}

	_, _, err = s.Decode(tokenString)
	if err == nil {
		t.Fatalf("no error on decoding with bad secret")
	}
}
