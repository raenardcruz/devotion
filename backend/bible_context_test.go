package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBibleContextHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectJSON     bool
		expectedCite   string
	}{
		{
			name:           "Valid citation (Psalms shortcut)",
			query:          "?passage=Psalms+23:1",
			expectedStatus: http.StatusOK,
			expectJSON:     true,
			expectedCite:   "Psalms 23:1",
		},
		{
			name:           "Valid citation using citation param (Psalms shortcut)",
			query:          "?citation=Psalms+121:1",
			expectedStatus: http.StatusOK,
			expectJSON:     true,
			expectedCite:   "Psalms 121:1",
		},
		{
			name:           "Missing parameters",
			query:          "",
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/context"+tt.query, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(bibleContextHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectJSON {
				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %q", contentType)
				}

				var resp BibleContext
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Errorf("failed to unmarshal JSON: %v", err)
				}

				if resp.Citation != tt.expectedCite {
					t.Errorf("expected citation %q, got %q", tt.expectedCite, resp.Citation)
				}
			}
		})
	}
}
