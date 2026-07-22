package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestApiTokenMiddleware(t *testing.T) {
	// Set up a dummy handler to wrap
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access Granted"))
	})

	protectedHandler := apiTokenMiddleware(dummyHandler)

	// Save original env value
	origTokenEnv := os.Getenv("API_TOKEN")
	defer os.Setenv("API_TOKEN", origTokenEnv)

	// Set test token
	testToken := "super_secret_test_token"
	os.Setenv("API_TOKEN", testToken)

	// Case 1: Success with correct token
	req := httptest.NewRequest(http.MethodGet, "/devotion", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "Access Granted" {
		t.Errorf("expected body 'Access Granted', got %q", rec.Body.String())
	}

	// Case 2: Unauthorized with missing Authorization header
	req = httptest.NewRequest(http.MethodGet, "/devotion", nil)
	rec = httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}

	// Case 3: Unauthorized with bad token format (missing bearer prefix)
	req = httptest.NewRequest(http.MethodGet, "/devotion", nil)
	req.Header.Set("Authorization", testToken)
	rec = httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}

	// Case 4: Unauthorized with wrong token
	req = httptest.NewRequest(http.MethodGet, "/devotion", nil)
	req.Header.Set("Authorization", "Bearer wrong_token_here")
	rec = httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}

	// Case 5: Internal Server Error when API_TOKEN env is not set
	os.Setenv("API_TOKEN", "")
	req = httptest.NewRequest(http.MethodGet, "/devotion", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec = httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 when API_TOKEN env is empty, got %d", rec.Code)
	}
}
