package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func apiTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configuredToken := os.Getenv("API_TOKEN")
		if configuredToken == "" {
			log.Println("[apiTokenMiddleware] WARNING: API_TOKEN environment variable is not set")
			respondWithError(w, http.StatusInternalServerError, "Server configuration error")
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format. Expected 'Bearer <token>'")
			return
		}

		clientToken := parts[1]
		if clientToken != configuredToken {
			respondWithError(w, http.StatusUnauthorized, "Invalid API token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
