package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMagisteriumChatHandler_Validation(t *testing.T) {
	// 1. Invalid method
	req := httptest.NewRequest(http.MethodGet, "/magisterium/chat", nil)
	w := httptest.NewRecorder()
	magisteriumChatHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %d", w.Code)
	}

	// 2. Empty payload
	req = httptest.NewRequest(http.MethodPost, "/magisterium/chat", bytes.NewBufferString("{}"))
	w = httptest.NewRecorder()
	magisteriumChatHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for empty payload, got %d", w.Code)
	}

	// 3. Missing API key with streaming
	payload := MagisteriumChatRequest{
		Messages: []MagisteriumChatMessage{
			{Role: "user", Content: "What is grace?"},
		},
		Stream: true,
	}
	body, _ := json.Marshal(payload)
	req = httptest.NewRequest(http.MethodPost, "/magisterium/chat", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	magisteriumChatHandler(w, req)

	// Since MAGISTERIUM_API_KEY might be empty in test env, check that if apiKey is empty,
	// streaming returns text/event-stream with error event
	settings := GetSettings()
	if settings.MagisteriumAPIKey == "" {
		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/event-stream") {
			t.Errorf("expected Content-Type text/event-stream, got %s", contentType)
		}
		if !strings.Contains(w.Body.String(), `"type":"error"`) {
			t.Errorf("expected error event in body, got %s", w.Body.String())
		}
	}
}

func TestPublicConversationsEndpoints(t *testing.T) {
	// 1. Save Public Conversation
	savePayload := SavePublicConversationRequest{
		ID:         "test-session-123",
		AuthorName: "St. Thomas",
		Title:      "Faith and Reason",
		Messages:   json.RawMessage(`[{"role":"user","content":"How does reason relate to faith?"}]`),
	}
	body, _ := json.Marshal(savePayload)
	req := httptest.NewRequest(http.MethodPost, "/magisterium/conversations/save", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	publicConversationsSaveHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK saving conversation, got %d: %s", w.Code, w.Body.String())
	}

	// 2. List Public Conversations
	req = httptest.NewRequest(http.MethodGet, "/magisterium/conversations/public", nil)
	w = httptest.NewRecorder()
	publicConversationsListHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing conversations, got %d", w.Code)
	}

	var listResp struct {
		Conversations []PublicConversation `json:"conversations"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}

	found := false
	for _, c := range listResp.Conversations {
		if c.ID == "test-session-123" {
			found = true
			if c.AuthorName != "St. Thomas" || c.Title != "Faith and Reason" {
				t.Errorf("unexpected conversation data: %+v", c)
			}
			break
		}
	}
	if !found {
		t.Errorf("expected conversation test-session-123 to be in list")
	}

	// 3. Delete Public Conversation (via POST with JSON body)
	delPayload := DeletePublicConversationRequest{ID: "test-session-123"}
	delBody, _ := json.Marshal(delPayload)
	req = httptest.NewRequest(http.MethodPost, "/magisterium/conversations/delete", bytes.NewBuffer(delBody))
	w = httptest.NewRecorder()
	publicConversationsDeleteHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK deleting conversation, got %d: %s", w.Code, w.Body.String())
	}

	// 4. Verify conversation is removed from list
	req = httptest.NewRequest(http.MethodGet, "/magisterium/conversations/public", nil)
	w = httptest.NewRecorder()
	publicConversationsListHandler(w, req)

	listResp.Conversations = nil
	json.Unmarshal(w.Body.Bytes(), &listResp)
	for _, c := range listResp.Conversations {
		if c.ID == "test-session-123" {
			t.Errorf("expected conversation test-session-123 to be removed from list")
		}
	}

	// 5. Delete via DELETE method with query param on /magisterium/conversations/public
	SavePublicConversation("test-session-456", "St. Augustine", "Confessions", []byte(`[]`))
	req = httptest.NewRequest(http.MethodDelete, "/magisterium/conversations/public?id=test-session-456", nil)
	w = httptest.NewRecorder()
	publicConversationsListHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK deleting via DELETE on public handler, got %d: %s", w.Code, w.Body.String())
	}

	// 6. Validation: missing ID
	req = httptest.NewRequest(http.MethodPost, "/magisterium/conversations/delete", bytes.NewBufferString("{}"))
	w = httptest.NewRecorder()
	publicConversationsDeleteHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for missing ID, got %d", w.Code)
	}

	// 7. Validation: invalid method
	req = httptest.NewRequest(http.MethodGet, "/magisterium/conversations/delete", nil)
	w = httptest.NewRecorder()
	publicConversationsDeleteHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %d", w.Code)
	}
}
