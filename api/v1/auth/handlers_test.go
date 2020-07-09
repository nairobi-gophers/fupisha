package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
	//start a test instance of the api server
	testApiServer(t)

	account := struct {
		email    string
		username string
		password string
	}{
		email:    "basebandit@github.com",
		username: "basebandit",
		password: "str0ngp@5520rd0nly_!",
	}

	payload, _ := json.Marshal(account)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	//instantiate our test auth object
	auth := testAuthResource(t)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.HandleSignup)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	//Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	//Check the response body is what we expect.
	expected := "{}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}
