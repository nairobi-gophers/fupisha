package tests

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/logging"
	"github.com/nairobi-gophers/fupisha/provider"
	"github.com/nairobi-gophers/fupisha/store/postgres"
)

func TestAuth(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	store, teardown := postgres.NewTestDatabase(t)
	t.Cleanup(teardown)

	mailer, err := provider.NewMailerWithSMTP(cfg, "../../../templates")
	if err != nil {
		t.Fatal(err)
	}

	logger := logging.NewLogger(cfg)
	logger.SetOutput(ioutil.Discard)

	testCfg := &api.ApiConfig{
		Logger:     logger,
		Cfg:        cfg,
		Store:      store,
		Mailer:     mailer,
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
			name:     "Create a new user",
			url:      "/auth/signup",
			method:   "POST",
			body:     `{"email":"parish@fupisha.io","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusCreated,
			wantBody: `{}`,
		},
		{
			name:     "Create a new user with an existing email",
			url:      "/auth/signup",
			method:   "POST",
			body:     `{"email":"parish@fupisha.io","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusConflict,
			wantBody: `{"status":"Conflict","error":"that email is taken"}`,
		},
		{
			name:     "Create a new user with an invalid password",
			url:      "/auth/signup",
			method:   "POST",
			body:     `{"email":"admin@fupisha.io","password":"str0ngpa55w0rd_"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"password: must contain English letters and digits only."}`,
		},
		{
			name:     "Create a new user with an invalid email",
			url:      "/auth/signup",
			method:   "POST",
			body:     `{"email":"invalid@fupisha","password":"str0ngpa55w0rd"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"email: must be a valid email address."}`,
		},
		{
			name:     "Login with a valid email and password",
			url:      "/auth/login",
			method:   "POST",
			body:     `{"email":"parish@fupisha.io","password":"str0ngpa55word"}`,
			wantCode: http.StatusUnauthorized,
			wantBody: `{"status":"Unauthorized","error":"invalid email or password"}`,
		},
		{
			name:     "Login with a valid email and invalid password",
			url:      "/auth/login",
			method:   "POST",
			body:     `{"email":"parish@fupisha.io","password":"str0ngpa55word_"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"status":"Unprocessable Entity","error":"password: must contain English letters and digits only."}`,
		},
	}

	for _, tc := range tests {
		req, err := http.NewRequest(tc.method, tc.url, ioutil.NopCloser(strings.NewReader(tc.body)))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Api", "v1")

		rr := httptest.NewRecorder()
		apiHandler.ServeHTTP(rr, req)

		t.Logf("%v", tc.name)

		if tc.wantCode != rr.Code {
			t.Fatalf("handler returned unexpected status code: want status code %d got %d", tc.wantCode, rr.Code)
		}

		if tc.wantBody != strings.TrimSuffix(rr.Body.String(), "\n") {
			t.Fatalf("handler returned unexpected body: want response body %q\n got %q", tc.wantBody, strings.TrimSuffix(rr.Body.String(), "\n"))
		}
	}

}
