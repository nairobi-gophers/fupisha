package auth_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/internal/config"
	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/store/mock"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSignup(t *testing.T) {

	testTime, err := time.Parse(time.RFC3339, "2020-02-03T04:05:06Z")
	if err != nil {
		t.Fatal(err)
	}

	testKey := encoding.GenUniqueID()

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	store := mock.Store{
		UserStore: mock.UserStore{
			OnNew: func(name, email, password string) (primitive.ObjectID, error) {
				id, _ := primitive.ObjectIDFromHex("5d0575344d9f7ff15e989174")
				return id, nil
			},
			OnGet: func(id string) (model.User, error) {
				switch id {
				case "5d0575344d9f7ff15e989174":
					id, _ := primitive.ObjectIDFromHex("5d0575344d9f7ff15e989174")
					return model.User{
						ID:                   id,
						Email:                "basebandit@fupisha.io",
						Password:             "str0ngpa55w0rd",
						Name:                 "basebandit",
						CreatedAt:            testTime,
						APIKey:               testKey,
						ResetPasswordExpires: testTime,
						ResetPasswordToken:   "TestResetPasswordToken",
						VerificationExpires:  testTime,
						VerificationToken:    testKey,
					}, nil
				}
				return model.User{}, errors.New("Not found")
			},
			OnGetByEmail: func(email string) (model.User, error) {
				switch email {
				case "basebandit@fupisha.io":
					id, _ := primitive.ObjectIDFromHex("5d0575344d9f7ff15e989174")
					return model.User{
						ID:                   id,
						Email:                "basebandit@fupisha.io",
						Password:             "str0ngpa55w0rd",
						Name:                 "basebandit",
						CreatedAt:            testTime,
						APIKey:               testKey,
						ResetPasswordExpires: testTime,
						ResetPasswordToken:   "TestResetPasswordToken",
						VerificationExpires:  testTime,
						VerificationToken:    testKey,
					}, nil
				}
				return model.User{}, errors.New("Not found")
			},
			OnSetAPIKey: func(id string, key uuid.UUID) (model.User, error) {
				uid, _ := primitive.ObjectIDFromHex(id)
				return model.User{
					ID:                   uid,
					Email:                "basebandit@fupisha.io",
					Password:             "str0ngpa55w0rd",
					Name:                 "basebandit",
					CreatedAt:            testTime,
					APIKey:               testKey,
					ResetPasswordExpires: testTime,
					ResetPasswordToken:   "TestResetPasswordToken",
					VerificationExpires:  testTime,
					VerificationToken:    testKey,
				}, nil
			},
		},
	}

	apiHandler, err := api.New(true, cfg, store)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		body     string
		wantCode int
		wantBody string
	}{
		{
			body:     `{"email":"basebandit@fupisha.io","name":"basebandit","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusCreated,
			wantBody: `{}`,
		},
		{
			body:     `{"email":"admin@fupisha.io","name":"admin","password":"str0ngpa55w0rd_"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"password: must contain English letters and digits only."}`,
		},
		{
			body:     `{"email":"invalid@fupisha","name":"admin","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"email: must be a valid email address."}`,
		},
		{
			body:     `{"email":"noreply@fupisha.io","name":"_","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"name: the length must be between 3 and 32."}`,
		},
	}

	for _, tc := range tests {
		url := "/auth/signup"

		req, err := http.NewRequest("POST", url, ioutil.NopCloser(strings.NewReader(tc.body)))
		if err != nil {
			t.Fatal(err)
		}

		//TODO: Add test for when API version is not set
		req.Header.Set("Api", "v1")

		rr := httptest.NewRecorder()
		apiHandler.ServeHTTP(rr, req)

		if tc.wantCode != rr.Code {
			t.Fatalf("handler returned unexpected status code: want status code %d got %d", tc.wantCode, rr.Code)
		}

		if tc.wantBody != strings.TrimSuffix(rr.Body.String(), "\n") {
			t.Fatalf("handler returned unexpected body: want response body %q got %q", tc.wantBody, strings.TrimSuffix(rr.Body.String(), "\n"))
		}
	}
}
