package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

type OllamaTagsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "admin" && req.Password == "admin" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AdminLoginResponse{
			Success: true,
			Token:   "admin-authenticated-token",
			Message: "Login successful",
		})
		return
	}

	respondWithError(w, http.StatusUnauthorized, "Invalid admin credentials")
}

func maskAPIKey(key string) string {
	if strings.TrimSpace(key) == "" {
		return ""
	}
	// Return a random number of bullet characters (between 8 and 14) to indicate a value exists without leaking the token
	bulletCount := 8 + rand.Intn(7)
	return strings.Repeat("•", bulletCount)
}

func isMaskedKey(key string) bool {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" {
		return false
	}
	return strings.Contains(trimmed, "•") || strings.Contains(trimmed, "●")
}

func adminSettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		settings := GetSettings()
		settings.GeminiAPIKey = maskAPIKey(settings.GeminiAPIKey)
		settings.MagisteriumAPIKey = maskAPIKey(settings.MagisteriumAPIKey)
		settings.BibleAPIKey = maskAPIKey(settings.BibleAPIKey)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)

	case http.MethodPost, http.MethodPut:
		var settings Settings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid settings payload")
			return
		}

		current := GetSettings()
		if isMaskedKey(settings.GeminiAPIKey) {
			settings.GeminiAPIKey = current.GeminiAPIKey
		}
		if isMaskedKey(settings.MagisteriumAPIKey) {
			settings.MagisteriumAPIKey = current.MagisteriumAPIKey
		}
		if isMaskedKey(settings.BibleAPIKey) {
			settings.BibleAPIKey = current.BibleAPIKey
		}

		if err := SaveSettings(settings); err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save settings: %v", err))
			return
		}

		respSettings := settings
		respSettings.GeminiAPIKey = maskAPIKey(settings.GeminiAPIKey)
		respSettings.MagisteriumAPIKey = maskAPIKey(settings.MagisteriumAPIKey)
		respSettings.BibleAPIKey = maskAPIKey(settings.BibleAPIKey)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"message":  "Settings saved successfully",
			"settings": respSettings,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func adminOllamaModelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/api/tags", ollamaURL))
	if err != nil {
		log.Printf("[adminOllamaModelsHandler] Warning: failed to fetch models from Ollama: %v", err)
		// Fallback model list if Ollama endpoint is unreachable
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"models": []string{"gemma4:cloud", "qwen2.5:4b", "llama3.2:3b", "mistral:latest"},
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read response from Ollama")
		return
	}

	var tagsResp OllamaTagsResponse
	if err := json.Unmarshal(body, &tagsResp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse response from Ollama")
		return
	}

	var modelNames []string
	for _, m := range tagsResp.Models {
		modelNames = append(modelNames, m.Name)
	}

	if len(modelNames) == 0 {
		modelNames = []string{"gemma4:cloud", "qwen2.5:4b", "llama3.2:3b"}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"models": modelNames,
	})
}

func adminGeminiModelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	geminiModels := []string{
		"gemini-3-flash-preview",
		"gemini-3.1-flash-lite",
		"gemini-3.5-flash",
		"gemini-3.5-flash-lite",
		"gemini-3.6-flash",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"models": geminiModels,
	})
}
