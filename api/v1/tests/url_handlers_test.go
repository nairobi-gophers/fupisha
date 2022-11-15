package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/api/v1/url"
	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/logging"
	"github.com/nairobi-gophers/fupisha/provider"
	"github.com/nairobi-gophers/fupisha/store/postgres"
)

func TestUrl(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	testURL := "https://www.youtube.com/watch?v=ZO3z966AqbU&t=14s"

	baseURL := fmt.Sprintf("%s:%s", cfg.BaseURL, cfg.Port)

	testParam, err := url.Shorten(testURL, cfg.ParamLength)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	testLink := baseURL + testParam

	ctx := context.Background()

	store, teardown := postgres.NewTestDatabase(t)
	t.Cleanup(teardown)

	const (
		testEmail    = "admin@fupisha.io"
		testPassword = "ih@veaStr0ngpassword"
	)

	u, err := store.NewUser(ctx, testEmail, testPassword)
	if err != nil {
		t.Fatalf("could not create test user %q", err)
	}

	_, err = store.NewURL(ctx, u.ID, testURL, testParam)
	if err != nil {
		t.Fatalf("could not insert the short param")
	}

	testSecret := "c4c0f2c42bde58f4d5f453483b3bed2b2915779cacff15526b2560b00748ec36"

	if len(cfg.JWT.Secret) == 0 {
		cfg.JWT.Secret = testSecret
	}

	if cfg.JWT.ExpireDelta == 0 {
		cfg.JWT.ExpireDelta = 6
	}

	jwtService, err := provider.NewJWTService(cfg)
	if err != nil {
		t.Fatal(err)
	}

	testToken, err := jwtService.Encode(u.ID.String())
	if err != nil {
		t.Fatal(err)
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
		name     string
		url      string
		method   string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "Shorten a valid url",
			url:      "/url/shorten",
			method:   "POST",
			body:     fmt.Sprintf(`{"url":"%s"}`, testURL),
			wantCode: http.StatusCreated,
			wantBody: fmt.Sprintf(`{"link":"%s"}`, testLink),
		},
		{
			name:     "Shorten an existing url",
			url:      "/url/shorten",
			method:   "POST",
			body:     fmt.Sprintf(`{"url":"%s"}`, testURL),
			wantCode: http.StatusCreated,
			wantBody: fmt.Sprintf(`{"link":"%s"}`, testLink),
		},
		{
			name:     "Resolve a short url",
			url:      testLink,
			method:   "GET",
			wantCode: http.StatusFound,
		},
	}

	for _, tc := range tests {
		req, err := http.NewRequest(tc.method, tc.url, ioutil.NopCloser(strings.NewReader(tc.body)))
		if err != nil {
			t.Fatal(err)
		}

		//TODO: Add test for when API version is not set
		req.Header.Set("Api", "v1")
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := httptest.NewRecorder()
		apiHandler.ServeHTTP(rr, req)

		t.Log(tc.name)

		if tc.wantCode != rr.Code {
			t.Fatalf("handler returned unexpected status: want status code %d got %d", tc.wantCode, rr.Code)
		}

		if tc.method != "GET" {
			if !strings.Contains(rr.Body.String(), "link") {
				t.Fatalf("handler returned unexpected body: want response body %q got %q", tc.wantBody, strings.TrimSuffix(rr.Body.String(), "\n"))
			}
		}
	}
}
