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

	s, teardown := NewTestDatabase(t)
	t.Cleanup(teardown)

	ctx := context.Background()

	wantEmail := "test_user@test.com"
	wantPassword := "test_password"

	u, err := s.Users().New(ctx, wantEmail, wantPassword)

	if err != nil {
		t.Fatalf("failed to create test_user1: %s", err)
	}

	got, err := s.Users().GetByEmail(ctx, wantEmail)
	if err != nil {
		t.Fatal(err)
	}

	sinceCreatedAt := time.Since(got.CreatedAt)

	if sinceCreatedAt > 3*time.Second || sinceCreatedAt < 0 {
		t.Fatalf("bad user.CreatedAt: %v", got.CreatedAt)
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
		Verified:             u.Verified,
		Password:             u.Password,
		CreatedAt:            u.CreatedAt,
		UpdatedAt:            u.UpdatedAt,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}

	if _, err := got.Compare(got.Password, wantPassword); err != nil {
		t.Fatalf("failed to compare password: %s", err)
	}

	u1, err := s.Users().GetByVerificationToken(ctx, u.VerificationToken)
	if err != nil {
		t.Fatalf("bad user verification token: %v", u.VerificationToken)
	}

	if !reflect.DeepEqual(u1, want) {
		t.Fatalf("got %+v want %+v", u1, want)
	}

	if err := s.Users().SetVerified(ctx, u1.ID); err != nil {
		t.Fatalf("failed to update the verified field: %s", err)
	}

	got1, err := s.Users().GetByEmail(ctx, u1.Email)
	if err != nil {
		t.Fatal(err)
	}

	verified := true

	if got1.Verified != verified {
		t.Fatalf("got %t want %t", got1.Verified, verified)
	}

}
