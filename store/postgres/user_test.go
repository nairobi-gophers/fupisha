package postgres

import (
	"context"
	"reflect"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/nairobi-gophers/fupisha/store"
)

func TestUser(t *testing.T) {

	s, err := newTestDatabase(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// wantName := "test_user"
	wantEmail := "test_user@test.com"
	wantPassword := "test_password"

	u, err := s.Users().New(ctx, wantEmail, wantPassword)

	if err != nil {
		t.Fatalf("failed to create test_user1: %s", err)
	}

	sinceCreatedAt := time.Since(u.CreatedAt)
	// fmt.Printf("sinceCreatedAt: %v, createdAt:%v", sinceCreatedAt, u.CreatedAt)

	if sinceCreatedAt > 3*time.Second {
		t.Fatalf("bad user.CreatedAt: %v", u.CreatedAt)
	}

	beforeVerificationExpiry := time.Until(u.VerificationExpires)

	//We are checking how many minutes we have until the verification token expires.
	//It cannot be 60 since some seconds elapse between creating the token and the point at which we are verifying it.
	//It should be less than 60.
	if beforeVerificationExpiry.Minutes() == 60 || beforeVerificationExpiry.Minutes() > 60 {
		t.Fatalf("bad user.VerificationExpires: %v", u.VerificationExpires)
	}

	want := store.User{
		ID:                   u.ID,
		Email:                wantEmail,
		ResetPasswordExpires: u.ResetPasswordExpires,
		ResetPasswordToken:   u.ResetPasswordToken,
		VerificationExpires:  u.VerificationExpires,
		VerificationToken:    u.VerificationToken,
		Verified:             false,
		Password:             u.Password,
		CreatedAt:            u.CreatedAt,
		UpdatedAt:            u.UpdatedAt,
	}

	if !reflect.DeepEqual(u, want) {
		t.Fatalf("got %+v want %+v", u, want)
	}

	if _, err := u.Compare(u.Password, wantPassword); err != nil {
		t.Fatalf("failed to compare password: %s", err)
	}
}
