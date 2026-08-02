package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

var (
	redisClient *redis.Client
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
		if err == nil {
			// Ensure passage texts are populated in case they weren't in cached string
			if err := PopulatePassageTexts(devotionData); err != nil {
				log.Printf("[fetchDevotionData] Warning: failed to populate passage texts: %v", err)
			}
			return devotionData, nil
		}
	}

	log.Printf("[fetchDevotionData] Cache miss for %s. Generating devotion concurrently using goroutines...", date)

	devotionData, err := generate_devotion(ctx, date)
	if err != nil {
		return nil, err
	}

	// Populate passage texts using get_bible_passage
	if err := PopulatePassageTexts(devotionData); err != nil {
		log.Printf("[fetchDevotionData] Warning: failed to populate passage texts: %v", err)
	}

	// Re-marshal to string to cache the updated data
	updatedDevotionDataBytes, err := json.Marshal(devotionData)
	if err == nil {
		err = redisClient.Set(ctx, cacheKey, string(updatedDevotionDataBytes), 24*time.Hour).Err()
		if err != nil {
			log.Printf("[fetchDevotionData] Error caching to Redis: %v", err)
		}
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

func generate_devotion(ctx context.Context, date string) (*DevotionData, error) {
	log.Printf("[generate_devotion] Step 1: Scraping mass readings for date: %s", date)
	readingsResp, err := get_mass_readings(date)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mass readings: %w", err)
	}

	devotion := &DevotionData{}

	// Categorize readings
	for _, r := range readingsResp.Readings {
		t := strings.ToLower(r.Type)
		rw := &ReadingWithContext{Citation: r.Citation}
		switch {
		case strings.Contains(t, "second reading") || strings.Contains(t, "reading 2") || strings.Contains(t, "reading ii"):
			devotion.SecondReading = rw
		case strings.Contains(t, "first reading") || strings.Contains(t, "reading 1") || strings.Contains(t, "reading i"):
			devotion.FirstReading = rw
		case strings.Contains(t, "psalm"):
			devotion.ResponsorialPsalm = rw
		case strings.Contains(t, "gospel"):
			devotion.Gospel = rw
		}
	}

	var wg sync.WaitGroup

	// Fetch Pope Quote concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		pq, err := get_pope_quote(date)
		if err != nil {
			log.Printf("[generate_devotion] Warning: failed to fetch Pope quote: %v", err)
			return
		}
		devotion.PopeQuote = pq.Text
	}()

	// Helper to fetch context in goroutine
	fetchContext := func(rw *ReadingWithContext) {
		if rw == nil || rw.Citation == "" {
			return
		}
		wg.Add(1)
		go func(target *ReadingWithContext) {
			defer wg.Done()
			bCtx, err := get_bible_context(ctx, GetBibleContextArgs{Citation: target.Citation})
			if err != nil {
				log.Printf("[generate_devotion] Warning: failed to fetch context for %s: %v", target.Citation, err)
				return
			}
			target.Context = bCtx.Context
		}(rw)
	}

	fetchContext(devotion.FirstReading)
	fetchContext(devotion.ResponsorialPsalm)
	fetchContext(devotion.SecondReading)
	fetchContext(devotion.Gospel)

	wg.Wait()
	log.Printf("[generate_devotion] Finished parallel fetching of context and Pope quote for %s", date)

	return devotion, nil
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
