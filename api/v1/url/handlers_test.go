package url

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/internal/config"
	"github.com/nairobi-gophers/fupisha/internal/store/mock"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandleShortenURL(t *testing.T) {

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	testURL := "https://www.youtube.com/watch?v=ZO3z966AqbU&t=14s"

	testLink := shortenURL(testURL, cfg.BaseURL, cfg.ParamLength)

	testTime, err := time.Parse(time.RFC3339, "2020-02-03T04:05:06Z")
	if err != nil {
		t.Fatal(err)
	}

	store := mock.Store{
		URLStore: mock.URLStore{
			OnNew: func(userID, originalURL, shortenedURL string) (interface{}, error) {
				id, _ := primitive.ObjectIDFromHex("5d0575344d9f7ff15e969178")
				return id, nil
			},
			OnGet: func(id string) (model.URL, error) {
				urlID, _ := primitive.ObjectIDFromHex(id)
				return model.URL{
					CreatedAt:    testTime,
					ID:           urlID,
					OriginalURL:  testURL,
					ShortenedURL: testLink,
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

		rr := httptest.NewRecorder()
		apiHandler.ServeHTTP(rr, req)

		if tc.wantCode != rr.Code {
			t.Fatalf("handler returned unexpected status")
		}

		if tc.wantBody != strings.TrimSuffix(rr.Body.String(), "\n") {
			t.Fatalf("handler returned unexpected body: want response body %q got %q", tc.wantBody, strings.TrimSuffix(rr.Body.String(), "\n"))
		}
	}
}
