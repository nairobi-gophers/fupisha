package mongodb

import (
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
)

func TestUser(t *testing.T) {
	s, teardown := testConn(t)
	defer teardown("users")

	id, err := s.Users().New("test_user1", "test_user1@test.com", "test_password1")

	if err != nil {
		t.Fatalf("failed to create test_user1: %s", err)
	}

	_, err = s.Users().New("test_user2", "test_user1@test.com", "test_password2")

	if err == nil {
		t.Fatalf("no error on creating account with an existing email")
	}

	_, err = s.Users().New("test_user3", "test_user3@test.com", "test_password3")

	if err != nil {
		t.Fatalf("failed to create test_user3: %s", err)
	}

	//convert the ObjectID to string "5d0575344d9f7ff15e989174"
	// uid := id.Hex()

	user, err := s.Users().Get(id)

	if err != nil {
		t.Fatalf("failed to get user by id: %s", err)
	}

	sinceCreated := time.Since(user.CreatedAt)
	if sinceCreated > 3*time.Second || sinceCreated < 0 {
		t.Fatalf("bad user.CreatedAt: %v", user.CreatedAt)
	}

	beforeVerificationExpiry := user.VerificationExpires.Sub(time.Now())

	//We are checking how many minutes we have until the verification token expires.
	//It cannot be 60 since some seconds elapse between creating the token and the point at which we are verifying it.
	//It should be less than 60.
	if beforeVerificationExpiry.Minutes() == 60 || beforeVerificationExpiry.Minutes() > 60 || beforeVerificationExpiry.Minutes() < 0 {
		t.Fatalf("bad user.VerificationExpires: %v", user.VerificationExpires)
	}

	want := model.User{
		ID:                   id,
		Name:                 "test_user1",
		Email:                "test_user1@test.com",
		ResetPasswordExpires: time.Time{},
		ResetPasswordToken:   "",
		VerificationExpires:  user.VerificationExpires,
		VerificationToken:    user.VerificationToken,
		Verified:             false,
		Password:             user.Password,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            time.Time{},
	}

	if !reflect.DeepEqual(user, want) {
		t.Fatalf("got user %+v want %+v", user, want)
	}

	if _, err := user.Compare(user.Password, "test_password1"); err != nil {
		t.Fatalf("failed to compare password: %s", err)
	}

	apiKey, _ := uuid.FromString("5bcd34d1-6bc2-464e-b0ab-6ca76f3c6f1b")

	err = s.Users().SetAPIKey(id, apiKey)
	if err != nil {
		t.Fatalf("failed to set api key: %s", err)
	}

	want = model.User{
		ID:                   id,
		Name:                 "test_user1",
		Email:                "test_user1@test.com",
		APIKey:               apiKey,
		ResetPasswordExpires: time.Time{},
		ResetPasswordToken:   "",
		VerificationExpires:  user.VerificationExpires,
		VerificationToken:    user.VerificationToken,
		Verified:             false,
		Password:             user.Password,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            time.Time{},
	}

	got, err := s.Users().Get(id)
	if err != nil {
		t.Fatalf("failed to get user by id: %s", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got user %+v want %+v", got, want)
	}

	got, err = s.Users().GetByEmail("test_user1@test.com")
	if err != nil {
		t.Fatalf("failed to get user by name: %s", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got user %v want %v", got, want)
	}

}
