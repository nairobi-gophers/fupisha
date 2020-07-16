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

func TestHandleSignup(t *testing.T) {

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
			body:     `{"email":"parish@fupisha.io","name":"parish","password":"str0ngpa55w0rd"}`,
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

func TestHandleLogin(t *testing.T) {

	testTime, err := time.Parse(time.RFC3339, "2020-02-03T04:05:06Z")
	if err != nil {
		t.Fatal(err)
	}

	testKey := encoding.GenUniqueID()

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	// jwtService, err := provider.NewJWTService(cfg.JWT.Secret)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// testToken, err := jwtService.Encode("5d0575344d9f7ff15e989174")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	store := mock.Store{
		UserStore: mock.UserStore{
			OnGet: func(id string) (model.User, error) {
				switch id {
				case "5d0575344d9f7ff15e989174":
					id, _ := primitive.ObjectIDFromHex("5d0575344d9f7ff15e989174")
					return model.User{
						ID:                   id,
						Email:                "parish@fupisha.io",
						Password:             "str0ngpa55w0rd",
						Name:                 "parish",
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
				case "parish@fupisha.io":
					id, _ := primitive.ObjectIDFromHex("5d0575344d9f7ff15e989174")
					usr := model.User{
						ID:                   id,
						Email:                "parish@fupisha.io",
						Password:             "str0ngpa55w0rd",
						Name:                 "parish",
						CreatedAt:            testTime,
						APIKey:               testKey,
						ResetPasswordExpires: testTime,
						ResetPasswordToken:   "TestResetPasswordToken",
						VerificationExpires:  testTime,
						VerificationToken:    testKey,
					}

					usr.HashPassword()

					return usr, nil
				}
				return model.User{}, errors.New("Not found")
			},
			OnSetAPIKey: func(id string, key uuid.UUID) (model.User, error) {
				uid, _ := primitive.ObjectIDFromHex(id)
				return model.User{
					ID:                   uid,
					Email:                "parish@fupisha.io",
					Password:             "str0ngpa55w0rd",
					Name:                 "parish",
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

	apiHandler, err := api.New(false, cfg, store)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		body     string
		wantCode int
		wantBody string
	}{
		//TODO: Add check for this test case. Note: token field is variable.
		// {
		// 	body:     `{"email":"parish@fupisha.io","password":"str0ngpa55w0rd"}`,
		// 	wantCode: http.StatusOK,
		// 	wantBody: fmt.Sprintf(`{"email":"parish@fupisha.io","name":"parish",
		// 	"id":"5d0575344d9f7ff15e989174",
		// 	"token":"%s"}`, testToken),
		// },
		{
			body:     `{"email":"parish@fupisha.io","password":"str0ngpa55word"}`,
			wantCode: http.StatusUnauthorized,
			wantBody: `{"status":"Unauthorized","error":"invalid email or password"}`,
		},
		{
			body:     `{"email":"parish@fupisha.io","password":"str0ngpa55word_"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"password: must contain English letters and digits only."}`,
		},
	}

	for _, tc := range tests {
		url := "/auth/login"

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
