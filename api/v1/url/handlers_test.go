package url_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/api/v1/url"
	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/logging"
	"github.com/nairobi-gophers/fupisha/provider"
	"github.com/nairobi-gophers/fupisha/store"
	"github.com/nairobi-gophers/fupisha/store/mock"
)

func TestHandleShortenURL(t *testing.T) {

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	testURL := "https://www.youtube.com/watch?v=ZO3z966AqbU&t=14s"

	if cfg.ParamLength == 0 {
		cfg.ParamLength = 6
	}

	testLink, err := url.Shorten(testURL, cfg.BaseURL, cfg.ParamLength)
	if err != nil {
		t.Fatal(err)
	}

	testTime, err := time.Parse(time.RFC3339, "2020-02-03T04:05:06Z")
	if err != nil {
		t.Fatal(err)
	}

	uid := "JBxPVPeqS72SqJzESF2TVw"

	testID, err := encoding.Decode(uid)
	if err != nil {
		t.Fatal(err)
	}

	testSecret := "c4c0f2c42bde58f4d5f453483b3bed2b2915779cacff15526b2560b00748ec36"

	if len(cfg.JWT.Secret) == 0 {
		cfg.JWT.Secret = testSecret
	}

	if cfg.JWT.ExpireDelta == 0 {
		cfg.JWT.ExpireDelta = 6
	}

	jwtService, err := provider.NewJWTService(cfg.JWT.Secret)
	if err != nil {
		t.Fatal(err)
	}

	testToken, err := jwtService.Encode(testID.String())
	if err != nil {
		t.Fatal(err)
	}

	testKey := encoding.GenUniqueID()

	store := mock.Store{
		UserStore: mock.UserStore{
			OnGet: func(ctx context.Context, id uuid.UUID) (store.User, error) {
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
		},
		URLStore: mock.URLStore{
			OnNew: func(ctx context.Context, userID uuid.UUID, originalURL, shortenedURL string) (store.URL, error) {

				switch userID {
				case testID:
					link := store.URL{
						ID:           encoding.GenUniqueID(),
						Owner:        testID,
						OriginalURL:  originalURL,
						ShortenedURL: shortenedURL,
						CreatedAt:    testTime,
						UpdatedAt:    testTime,
					}
					return link, nil
				}
				return store.URL{}, errors.New("user ID not found")
			},
			OnGet: func(ctx context.Context, id uuid.UUID) (store.URL, error) {
				switch id {
				case testID:
					return store.URL{
						ID:           testID,
						OriginalURL:  testURL,
						ShortenedURL: testLink,
						CreatedAt:    testTime,
						UpdatedAt:    testTime,
					}, nil
				}
				return store.URL{}, errors.New("user ID not found")
			},
		},
	}

	logger := logging.NewLogger(cfg)
	logger.SetOutput(ioutil.Discard)

	testCfg := &api.ApiConfig{
		Logger:     logger,
		Cfg:        cfg,
		Store:      store,
		EnableCORS: false,
	}

	apiHandler, err := api.New(testCfg)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		body     string
		wantCode int
		wantBody string
	}{
		{
			body:     fmt.Sprintf(`{"url":"%s"}`, testURL),
			wantCode: http.StatusCreated,
			wantBody: fmt.Sprintf(`{"link":"%s"}`, testLink),
		},
	}

	for _, tc := range tests {
		url := "/url/shorten"

		req, err := http.NewRequest("POST", url, ioutil.NopCloser(strings.NewReader(tc.body)))
		if err != nil {
			t.Fatal(err)
		}

		//TODO: Add test for when API version is not set
		req.Header.Set("Api", "v1")
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := httptest.NewRecorder()
		apiHandler.ServeHTTP(rr, req)

		if tc.wantCode != rr.Code {
			t.Fatalf("handler returned unexpected status: want status code %d got %d", tc.wantCode, rr.Code)
		}

		if !strings.Contains(rr.Body.String(), "link") {
			t.Fatalf("handler returned unexpected body: want response body %q got %q", tc.wantBody, strings.TrimSuffix(rr.Body.String(), "\n"))
		}

	}
}
