package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	app := setupTestApplication(t)
	h := app.Mux()

	req, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	var expected = "ping"
	if res.Body.String() != expected {
		t.Errorf("expected body of %q but got %q", res.Body.String(), expected)
	}

}
