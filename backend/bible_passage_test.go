package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func TestPassageSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal passage",
			input:    "Genesis 1:1",
			expected: "Genesis 1:1",
		},
		{
			name:     "With 'and'",
			input:    "Genesis 1:1 and 2",
			expected: "Genesis 1:1,2",
		},
		{
			name:     "With spaces around colon",
			input:    "Genesis 1 : 1",
			expected: "Genesis 1:1",
		},
		{
			name:     "With spaces around hyphen",
			input:    "Genesis 1:1 - 5",
			expected: "Genesis 1:1-5",
		},
		{
			name:     "Complex case",
			input:    "James 1:12-18, 20 and 22",
			expected: "James 1:12-18,20,22",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := passage_sanitize(tt.input)
			if got != tt.expected {
				t.Errorf("passage_sanitize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetBibleBookShort(t *testing.T) {
	// bible_books is a global variable in bible_constant.go
	tests := []struct {
		name     string
		bookName string
		expected string
	}{
		{
			name:     "Genesis",
			bookName: "Genesis",
			expected: "GEN",
		},
		{
			name:     "Exodus",
			bookName: "Exodus",
			expected: "EXO",
		},
		{
			name:     "Matthew",
			bookName: "Matthew",
			expected: "MAT",
		},
		{
			name:     "Unknown book",
			bookName: "NonExistent",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := get_bible_book_short(tt.bookName)
			if got != tt.expected {
				t.Errorf("get_bible_book_short(%q) = %q, want %q", tt.bookName, got, tt.expected)
			}
		})
	}
}

func TestGetBiblePassage(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println("[initializeEnvironment] Error loading .env file")
	}

	api_key := os.Getenv("BIBLE_API_KEY")
	if api_key == "" {
		t.Skip("BIBLE_API_KEY not set, skipping integration test")
	}

	tests := []struct {
		name    string
		passage string
		wantErr bool
	}{
		{
			name:    "Valid single verse",
			passage: "John 3:7b-15",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Getting bible passage for", tt.passage)

			text, copyright, err := get_bible_passage(tt.passage)
			if (err != nil) != tt.wantErr {
				t.Errorf("get_bible_passage(%q) error = %v, wantErr %v", tt.passage, err, tt.wantErr)
				return
			}

			fmt.Println("Text:", text)
			fmt.Println("Copyright:", copyright)

			if !tt.wantErr {
				if text == "" {
					t.Errorf("get_bible_passage(%q) returned empty text", tt.passage)
				}
				if copyright == "" {
					t.Errorf("get_bible_passage(%q) returned empty copyright", tt.passage)
				}
			}
		})
	}
}
func TestStripVerseSuffix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1", "1"},
		{"1a", "1"},
		{"1ab", "1"},
		{"7b", "7"},
		{"15", "15"},
		{"1cd", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := strip_verse_suffix(tt.input)
			if got != tt.expected {
				t.Errorf("strip_verse_suffix(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBiblePassageCaching(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println("[TestBiblePassageCaching] Error loading .env file")
	}

	// Initialize a local redis client for testing
	testRedisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	
	ctx := context.Background()
	if err := testRedisClient.Ping(ctx).Err(); err != nil {
		t.Skip("Redis is not running on localhost:6379, skipping cache integration test")
	}
	
	// Temporarily set the global redisClient to our test client
	oldRedisClient := redisClient
	redisClient = testRedisClient
	defer func() {
		redisClient = oldRedisClient
	}()

	testPassage := "TestCachePassage 1:1"
	cacheKey := "bible_passage:" + testPassage

	// Ensure cache is clean
	redisClient.Del(ctx, cacheKey)
	defer redisClient.Del(ctx, cacheKey)

	// Step 1: Pre-populate cache directly with dummy data
	dummyContent := "This is cached content"
	dummyCopyright := "Cached Copyright"
	cachedData := BibleContent{
		Content:   dummyContent,
		Copyright: dummyCopyright,
	}
	cachedVal, err := json.Marshal(cachedData)
	if err != nil {
		t.Fatalf("failed to marshal dummy data: %v", err)
	}

	err = redisClient.Set(ctx, cacheKey, string(cachedVal), 10*time.Second).Err()
	if err != nil {
		t.Fatalf("failed to prepopulate cache: %v", err)
	}

	// Step 2: Call get_bible_passage and verify it hits the cache and returns our dummy data
	content, copyright, err := get_bible_passage(testPassage)
	if err != nil {
		t.Fatalf("expected get_bible_passage to succeed from cache, got error: %v", err)
	}

	if content != dummyContent {
		t.Errorf("expected content %q, got %q", dummyContent, content)
	}
	if copyright != dummyCopyright {
		t.Errorf("expected copyright %q, got %q", dummyCopyright, copyright)
	}

	// Step 3: Test cache write if API keys are set
	api_key := os.Getenv("BIBLE_API_KEY")
	bible_id := os.Getenv("BIBLE_ID")
	if api_key != "" && bible_id != "" {
		testWritePassage := "John 3:16"
		writeCacheKey := "bible_passage:" + testWritePassage
		redisClient.Del(ctx, writeCacheKey)
		defer redisClient.Del(ctx, writeCacheKey)

		// Call get_bible_passage
		content, _, err := get_bible_passage(testWritePassage)
		if err == nil && content != "" {
			// Check if it's now in Redis
			val, err := redisClient.Get(ctx, writeCacheKey).Result()
			if err != nil {
				t.Errorf("expected cache key to be populated after get_bible_passage call, got error: %v", err)
			} else {
				var cached BibleContent
				if err := json.Unmarshal([]byte(val), &cached); err != nil {
					t.Errorf("failed to unmarshal cached data: %v", err)
				} else if cached.Content != content {
					t.Errorf("expected cached content to match API response, cached: %q, response: %q", cached.Content, content)
				}
			}
		}
	}
}
