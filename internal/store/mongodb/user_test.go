package mongodb

import (
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUser(t *testing.T) {
	s, teardown := testConn(t)
	defer teardown("users")

	id, err := s.Users().New("test_user1", "test_user1@test.com", "test_password1")

	if err != nil {
		t.Fatalf("failed to create test_user1: %s", err)
	}

	if _, ok := id.(primitive.ObjectID); !ok {
		t.Fatalf("failed to assert the created test_user1 insert id")
	}

	_, err = s.Users().New("test_user2", "test_user2@test.com", "test_password2")

	if err != nil {
		t.Fatalf("failed to create test_user2: %s", err)
	}

	//convert the ObjectID to string "5d0575344d9f7ff15e989174"
	uid := id.(primitive.ObjectID).Hex()

	user, err := s.Users().Get(uid)

	if err != nil {
		t.Fatalf("failed to get user by id: %s", err)
	}

	sinceCreated := time.Since(user.CreatedAt)
	if sinceCreated > 3*time.Second || sinceCreated < 0 {
		t.Fatalf("bad user.CreatedAt: %v", user.CreatedAt)
	}

	want := model.User{
		ID:                   id.(primitive.ObjectID),
		Name:                 "test_user1",
		Email:                "test_user1@test.com",
		ResetPasswordExpires: time.Time{},
		ResetPasswordToken:   "",
		VerificationExpires:  time.Time{},
		VerificationToken:    "",
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

	user, err = s.Users().SetAPIKey(uid, apiKey)
	if err != nil {
		t.Fatalf("failed to set api key: %s", err)
	}

	want = model.User{
		ID:                   id.(primitive.ObjectID),
		Name:                 "test_user1",
		Email:                "test_user1@test.com",
		APIKey:               apiKey,
		ResetPasswordExpires: time.Time{},
		ResetPasswordToken:   "",
		VerificationExpires:  time.Time{},
		VerificationToken:    "",
		Verified:             false,
		Password:             user.Password,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            time.Time{},
	}

	if !reflect.DeepEqual(user, want) {
		t.Fatalf("got user %+v want %+v", user, want)
	}
}
