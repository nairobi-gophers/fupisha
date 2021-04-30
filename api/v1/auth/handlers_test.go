package auth_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/store"
	"github.com/nairobi-gophers/fupisha/store/mock"
)

func TestHandleSignup(t *testing.T) {

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	store := mock.Store{
		UserStore: mock.UserStore{
			OnNew: func(ctx context.Context, email, password string) (store.User, error) {

				user := store.User{
					Email:    email,
					Password: password,
				}

				if err := user.HashPassword(); err != nil {
					return store.User{}, err
				}

				return user, nil
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
			body:     `{"email":"parish@fupisha.io","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusCreated,
			wantBody: `{}`,
		},
		{
			body:     `{"email":"admin@fupisha.io","password":"str0ngpa55w0rd_"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"password: must contain English letters and digits only."}`,
		},
		{
			body:     `{"email":"invalid@fupisha","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"email: must be a valid email address."}`,
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
	testSecret := "c4c0f2c42bde58f4d5f453483b3bed2b2915779cacff15526b2560b00748ec36"

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.JWT.Secret) == 0 {
		cfg.JWT.Secret = testSecret
	}

	if cfg.JWT.ExpireDelta == 0 {
		cfg.JWT.ExpireDelta = 6
	}

	uid := "JBxPVPeqS72SqJzESF2TVw"

	testID, err := encoding.Decode(uid)
	if err != nil {
		t.Fatal(err)
	}

	// jwtService, err := provider.NewJWTService(cfg.JWT.Secret)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// testToken, err := jwtService.Encode(testID.String())
	// if err != nil {
	// 	t.Fatal(err)
	// }

	store := mock.Store{
		UserStore: mock.UserStore{
			OnGet: func(ctxt context.Context, id uuid.UUID) (store.User, error) {
				switch id {
				case testID:
					u := store.User{
						ID:                  testID,
						Email:               "parish@fupisha.io",
						Password:            "str0ngpa55w0rd",
						CreatedAt:           testTime,
						UpdatedAt:           testTime,
						VerificationExpires: testTime.Add(time.Minute * 60),
						VerificationToken:   testKey,
					}

					u.HashPassword()

					return u, nil
				}
				return store.User{}, errors.New("Not found")
			},
			OnGetByEmail: func(ctx context.Context, email string) (store.User, error) {
				switch email {
				case "parish@fupisha.io":

					usr := store.User{
						ID:                  testID,
						Email:               "parish@fupisha.io",
						Password:            "str0ngpa55w0rd",
						CreatedAt:           testTime,
						APIKey:              testKey,
						VerificationExpires: testTime.Add(time.Minute * 60),
						VerificationToken:   testKey,
					}

					usr.HashPassword()

					return usr, nil
				}
				return store.User{}, errors.New("Not found")
			},
			OnSetAPIKey: func(ctx context.Context, id, key uuid.UUID) error {

				switch id {

				case testID:
					u := store.User{
						ID:                  testID,
						Email:               "parish@fupisha.io",
						Password:            "str0ngpa55w0rd",
						CreatedAt:           testTime,
						APIKey:              testKey,
						VerificationExpires: testTime.Add(time.Minute * 60),
						VerificationToken:   testKey,
					}

					u.HashPassword()

					return nil
				}

				return errors.New("Not found")
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
		//TODO: Test for valid jwt token upon successful login.Add check for this test case. Note: token field is variable, it changes with each request call.
		// {
		// 	body:     `{"email":"parish@fupisha.io","password":"str0ngpa55w0rd"}`,
		// 	wantCode: http.StatusOK,
		// 	wantBody: fmt.Sprintf(`{"email":"parish@fupisha.io","id":"%s",
		// 	"token":"%s"}`, testID, testToken),
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
			t.Fatalf("handler returned unexpected body: want response body %q\n got %q", tc.wantBody, strings.TrimSuffix(rr.Body.String(), "\n"))
		}
	}
}
