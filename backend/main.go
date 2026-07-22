package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

var (
	redisClient       *redis.Client
	adkRunner         *runner.Runner
	adkSessionService session.Service
)

func main() {
	log.Println("[main] Devotion API Started")
	ctx := context.Background()
	initializeEnvironment()

	// Initialize Database
	InitDB()

	// Initialize Redis
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("[main] Failed to connect to Redis: %v", err)
	} else {
		log.Println("[main] Successfully connected to Redis")
	}

	settings := GetSettings()
	gemini_api_key := settings.GeminiAPIKey
	if gemini_api_key == "" {
		log.Printf("[main] Warning: GEMINI_API_KEY not set in DB settings. Please configure in /admin.")
		gemini_api_key = "dummy-key-pending-admin-setup"
	}

	geminiModelName := settings.GeminiModel
	if geminiModelName == "" {
		geminiModelName = "gemini-3.1-flash-lite"
	}

	model, err := gemini.NewModel(ctx, geminiModelName, &genai.ClientConfig{
		APIKey: gemini_api_key,
	})
	if err != nil {
		log.Printf("[main] Warning: Failed to create default Gemini model: %v", err)
	}

	if model != nil {
		catholicAgent, err := llmagent.New(llmagent.Config{
			Name:        "CatholicAssistant",
			Model:       model,
			Description: "An expert Catholic assistant providing daily devotionals.",
			Instruction: "You are a wise and compassionate Catholic assistant. Your mission is to help the faithful prepare for the day by providing the daily mass readings, deep historical and theological context, and the inspiring 'Words of the Popes' from Vatican News.\n\nFollow these steps faithfully:\n1. Fetch the daily mass readings using the `GetMassReadings` tool.\n2. For every reading found (First Reading, Psalm, Second Reading, Gospel), obtain its historical and spiritual context using the `GetBibleContext` tool.\n3. Retrieve 'The words of the Popes' for the given date using the `GetPopeQuote` tool.\n4. Synthesize all this information into a beautifully structured daily devotion. When populating the `context` field for each reading in the `SubmitDevotion` tool call, you MUST use the exact, unmodified output returned by the `GetBibleContext` tool for that reading. Do not summarize, rewrite, or synthesize your own context for the readings; strictly respect and copy the output of the `GetBibleContext` tool.\n5. FINAL STEP: Call the `SubmitDevotion` tool with the complete JSON data. Leave the `text` field blank for each reading as it will be automatically populated from a word-for-word Bible API. This is your most important duty. Do not just output the JSON as text; you MUST call the tool.",
			Tools:       get_tools(),
		})
		if err != nil {
			log.Printf("[main] Failed to create ADK agent: %v", err)
		} else {
			adkSessionService = session.InMemoryService()
			adkRunner, err = runner.New(runner.Config{
				AppName:        "DevotionAPI",
				Agent:          catholicAgent,
				SessionService: adkSessionService,
			})
			if err != nil {
				log.Printf("[main] Failed to create ADK runner: %v", err)
			}
		}
	}

	http.Handle("/devotion", corsMiddleware(apiTokenMiddleware(http.HandlerFunc(devotionHandler))))
	http.Handle("/bible", corsMiddleware(apiTokenMiddleware(http.HandlerFunc(biblePassageHandler))))
	http.Handle("/context", corsMiddleware(apiTokenMiddleware(http.HandlerFunc(bibleContextHandler))))

	// Admin API Endpoints
	http.Handle("/api/admin/login", corsMiddleware(http.HandlerFunc(adminLoginHandler)))
	http.Handle("/api/admin/settings", corsMiddleware(http.HandlerFunc(adminSettingsHandler)))
	http.Handle("/api/admin/ollama-models", corsMiddleware(http.HandlerFunc(adminOllamaModelsHandler)))
	http.Handle("/api/admin/gemini-models", corsMiddleware(http.HandlerFunc(adminGeminiModelsHandler)))

	// Start Cron Jobs
	startCronJobs(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[main] Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("[main] %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func devotionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		// Basic validation: ensure it parses
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.Printf("[devotionHandler] Invalid date format: %s", date)
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	devotionData, err := fetchDevotionData(ctx, date)
	if err != nil {
		log.Printf("[devotionHandler] Error fetching devotion data: %v", err)
		http.Error(w, "Error fetching devotion data", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(devotionData)
	if err != nil {
		log.Printf("[devotionHandler] Error marshaling response: %v", err)
		http.Error(w, "Error producing JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func fetchDevotionData(ctx context.Context, date string) (*DevotionData, error) {
	cacheKey := fmt.Sprintf("devotion:%s", date)

	// Check Redis cache
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Printf("[fetchDevotionData] Cache hit for %s", date)
		devotionData, err := ExtractDevotionData(val)
		if err != nil {
			return nil, err
		}

		// Ensure passage texts are populated in case they weren't in the cached string
		if err := PopulatePassageTexts(devotionData); err != nil {
			log.Printf("[fetchDevotionData] Warning: failed to populate passage texts: %v", err)
		}

		return devotionData, nil
	}

	log.Printf("[fetchDevotionData] Cache miss for %s. Generating devotion using ADK...", date)

	devotionDataStr, err := generate_devotion(ctx, date)
	if err != nil {
		return nil, err
	}

	devotionData, err := ExtractDevotionData(devotionDataStr)
	if err != nil {
		return nil, err
	}

	// Populate passage texts using get_bible_passage
	if err := PopulatePassageTexts(devotionData); err != nil {
		log.Printf("[fetchDevotionData] Warning: failed to populate passage texts: %v", err)
	}

	// Re-marshal to string to cache the updated data
	updatedDevotionDataStr, err := json.Marshal(devotionData)
	if err == nil {
		devotionDataStr = string(updatedDevotionDataStr)
	}

	// Cache the result
	err = redisClient.Set(ctx, cacheKey, devotionDataStr, 24*time.Hour).Err()
	if err != nil {
		log.Printf("[fetchDevotionData] Error caching to Redis: %v", err)
	}

	return devotionData, nil
}

func startCronJobs(ctx context.Context) {
	c := cron.New()

	// Run every 12 AM (midnight)
	_, err := c.AddFunc("0 0 * * *", func() {
		today := time.Now().Format("2006-01-02")
		log.Printf("[cron] Prepopulating cache for today: %s", today)
		_, err := fetchDevotionData(ctx, today)
		if err != nil {
			log.Printf("[cron] Error prepopulating cache for %s: %v", today, err)
		} else {
			log.Printf("[cron] Successfully prepopulated cache for %s", today)
		}
	})

	if err != nil {
		log.Printf("[cron] Error scheduling cron job: %v", err)
		return
	}

	c.Start()
	log.Println("[cron] Scheduled daily prepopulation job for 12 AM")
}

func generate_devotion(ctx context.Context, date string) (string, error) {
	sessionID := fmt.Sprintf("devotion-%s", strings.ReplaceAll(date, "-", ""))
	userID := "system"

	// Ensure session exists
	_, err := adkSessionService.Create(ctx, &session.CreateRequest{
		AppName:   "DevotionAPI",
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Printf("[generate_devotion] Warning calling session Create: %v", err)
	}

	var lastText strings.Builder
	var lastToolResult string

	for event, err := range adkRunner.Run(ctx, userID, sessionID, genai.NewContentFromText(fmt.Sprintf("Generate devotion for %s. Use the GetMassReadings tool with the specific date.", date), "user"), agent.RunConfig{}) {
		if err != nil {
			return "", fmt.Errorf("error during agent run: %w", err)
		}

		if event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					lastText.WriteString(part.Text)
				}
				if part.FunctionResponse != nil && part.FunctionResponse.Name == "SubmitDevotion" {
					b, _ := json.Marshal(part.FunctionResponse.Response)
					lastToolResult = string(b)
				}
			}
		}
	}

	if lastToolResult != "" {
		return lastToolResult, nil
	}
	return lastText.String(), nil
}

func initializeEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[initializeEnvironment] Error loading .env file")
	}
}

func toMap(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var res any
	err = json.Unmarshal(b, &res)
	return res, err
}

func biblePassageHandler(w http.ResponseWriter, r *http.Request) {
	passage := r.URL.Query().Get("passage")
	if passage == "" {
		http.Error(w, "Missing passage parameter", http.StatusBadRequest)
		return
	}
	passageText, copyRight, err := get_bible_passage(passage)
	if err != nil {
		http.Error(w, "Error getting bible passage", http.StatusInternalServerError)
		return
	}
	response := BibleContent{
		Content:   passageText,
		Copyright: copyRight,
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error producing JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func bibleContextHandler(w http.ResponseWriter, r *http.Request) {
	passage := r.URL.Query().Get("passage")
	if passage == "" {
		passage = r.URL.Query().Get("citation")
	}
	if passage == "" {
		http.Error(w, "Missing passage or citation parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	res, err := get_bible_context(ctx, GetBibleContextArgs{Citation: passage})
	if err != nil {
		log.Printf("[bibleContextHandler] Error getting context: %v", err)
		http.Error(w, "Error getting bible context", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Printf("[bibleContextHandler] Error marshaling response: %v", err)
		http.Error(w, "Error producing JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
