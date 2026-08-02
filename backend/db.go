package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Settings struct {
	ContextProvider     string `json:"context_provider"`
	FactCheckerProvider string `json:"fact_checker_provider"`
	GeminiAPIKey        string `json:"gemini_api_key"`
	MagisteriumAPIKey   string `json:"magisterium_api_key"`
	BibleAPIKey         string `json:"bible_api_key"`
	GeminiModel         string `json:"gemini_model"`
	OllamaModel         string `json:"ollama_model"`
	ContextModel        string `json:"context_model"`
	ContextInstruction  string `json:"context_instruction"`
	EnableFactChecker   bool   `json:"enable_fact_checker"`
}

var (
	db             *sql.DB
	memorySettings Settings
)

const defaultInstruction = `As a Catholic Bible scholar guided by Catholic Social Teaching and Pope Leo XIV's encyclical Magnifica Humanitas, provide a detailed and insightful context for the passage {{citation}}.
{{passage_text}}

Focus on the historical setting, literal, moral, allegorical and anagogical senses of the scripture, upholding human dignity, truth, and spiritual growth. Relate this passage to the Catechism of the Catholic Church (CCC) or Catholic doctrine if possible. Additionally, explore the typology (e.g., how an Old Testament passage prefigures New Testament realities, or how a New Testament passage fulfills Old Testament foreshadowing). Where relevant, mention early Church writers, Doctors of the Church (such as Saint Justin Martyr, Saint Ignatius, Saint Irenaeus, or Saint Augustine), and papal teachings. Use clear, accessible English. Structure your response with clean formatting and typography (using Markdown bold, italics, or lists where helpful to make the details highly readable). Limit your response to a maximum of 15 sentences.

ETHICAL & ETHNICITY DIRECTIVES (Magnifica Humanitas):
- Respect human dignity and truth: do not invent facts, speculative doctrines, or non-authentic interpretations.
- Maintain transparency: this context is an AI-assisted tool intended to aid reflection, pointing the believer back to prayer, Sacred Tradition, and real community life.

CRITICAL CITATION RULES:

DEFAULT TO GENERAL DOCTRINAL REFERENCES: Because exact paragraph numbers are highly prone to hallucination, do not attempt to guess specific CCC paragraph numbers. Instead, name the specific Catholic doctrine, section, or theme being referenced (e.g., "The Catechism's teaching on Grace," or "The CCC section on the Sacrament of Baptism"). Only include an exact paragraph number if you are executing a verified internal data match; otherwise, omit the number entirely and use the doctrinal name.

VERIFY TEXTUAL ALIGNMENT: Before outputting a specific CCC paragraph number or Bible verse, verify that the text you are writing explicitly matches the actual content of that number/verse. If there is a mismatch between the description and the number, delete the number and keep only the general text description.

STRICT BIBLE CITATION: Any verse cited must exactly match the text and citation of that verse. If unsure of the exact verse number, refer to the chapter generally.`

func getDefaultSettings() Settings {
	return Settings{
		ContextProvider:     "ollama",
		FactCheckerProvider: "gemini",
		GeminiAPIKey:        "",
		MagisteriumAPIKey:   "",
		BibleAPIKey:         "",
		GeminiModel:         "gemini-3.1-flash-lite",
		OllamaModel:         "gemma4:cloud",
		ContextModel:        "gemma4:cloud",
		ContextInstruction:  defaultInstruction,
		EnableFactChecker:   true,
	}
}

func isDBNotExistError(err error) bool {
	if err == nil {
		return false
	}
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "3D000" { // invalid_catalog_name
			return true
		}
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "3d000")
}

func parseConnDetails(connStr string) (targetDB string, maintenanceConnStr string, err error) {
	if strings.HasPrefix(connStr, "postgres://") || strings.HasPrefix(connStr, "postgresql://") {
		u, parseErr := url.Parse(connStr)
		if parseErr != nil {
			return "", "", parseErr
		}
		targetDB = strings.TrimPrefix(u.Path, "/")
		if targetDB == "" || targetDB == "postgres" {
			return targetDB, connStr, nil
		}
		u.Path = "/postgres"
		maintenanceConnStr = u.String()
		return targetDB, maintenanceConnStr, nil
	}

	re := regexp.MustCompile(`dbname=([^\s]+)`)
	matches := re.FindStringSubmatch(connStr)
	if len(matches) > 1 {
		targetDB = matches[1]
	}
	if targetDB == "" || targetDB == "postgres" {
		return targetDB, connStr, nil
	}
	maintenanceConnStr = re.ReplaceAllString(connStr, "dbname=postgres")
	return targetDB, maintenanceConnStr, nil
}

func ensureDatabaseExists(connStr string) error {
	targetDB, maintenanceConnStr, err := parseConnDetails(connStr)
	if err != nil || targetDB == "" || targetDB == "postgres" {
		return err
	}

	// Try pinging or connecting to target DB first
	targetDBConn, err := sql.Open("postgres", connStr)
	if err == nil {
		pingErr := targetDBConn.Ping()
		targetDBConn.Close()
		if pingErr == nil {
			return nil
		}
		if !isDBNotExistError(pingErr) {
			return pingErr
		}
	}

	log.Printf("[InitDB] Target database %q does not exist. Attempting to create it...", targetDB)

	mDB, err := sql.Open("postgres", maintenanceConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to maintenance db: %w", err)
	}
	defer mDB.Close()

	if err := mDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping maintenance db: %w", err)
	}

	createQuery := fmt.Sprintf("CREATE DATABASE %s;", pq.QuoteIdentifier(targetDB))
	if _, err := mDB.Exec(createQuery); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("failed to create database %s: %w", targetDB, err)
		}
	}

	log.Printf("[InitDB] Database %q created successfully", targetDB)
	return nil
}

func InitDB() {
	memorySettings = getDefaultSettings()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		host := os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("POSTGRES_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("POSTGRES_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("POSTGRES_PASSWORD")
		if password == "" {
			password = "postgres"
		}
		dbname := os.Getenv("POSTGRES_DB")
		if dbname == "" {
			dbname = "devotion"
		}
		sslmode := os.Getenv("POSTGRES_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	}

	if err := ensureDatabaseExists(connStr); err != nil {
		log.Printf("[InitDB] Database creation notice: %v", err)
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("[InitDB] Warning: Failed to open postgres connection: %v. Using default in-memory settings.", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Printf("[InitDB] Warning: Could not ping PostgreSQL: %v. System will use fallbacks until DB is available.", err)
		return
	}

	log.Println("[InitDB] Successfully connected to PostgreSQL")

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS settings (
		id INT PRIMARY KEY DEFAULT 1,
		context_provider VARCHAR(50) NOT NULL DEFAULT 'ollama',
		gemini_api_key TEXT NOT NULL DEFAULT '',
		bible_api_key TEXT NOT NULL DEFAULT '',
		gemini_model VARCHAR(100) NOT NULL DEFAULT 'gemini-3.1-flash-lite',
		ollama_model VARCHAR(100) NOT NULL DEFAULT 'gemma4:cloud',
		context_model VARCHAR(100) NOT NULL DEFAULT 'gemma4:cloud',
		context_instruction TEXT NOT NULL DEFAULT '',
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT single_row CHECK (id = 1)
	);
	
	CREATE TABLE IF NOT EXISTS public_chat_conversations (
		id VARCHAR(100) PRIMARY KEY,
		author_name VARCHAR(100) NOT NULL,
		title VARCHAR(200) NOT NULL,
		messages JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Printf("[InitDB] Error creating tables: %v", err)
		return
	}

	// Migrations for missing columns on existing tables
	migrations := []string{
		`ALTER TABLE settings ADD COLUMN IF NOT EXISTS magisterium_api_key TEXT NOT NULL DEFAULT '';`,
		`ALTER TABLE settings ADD COLUMN IF NOT EXISTS bible_api_key TEXT NOT NULL DEFAULT '';`,
		`ALTER TABLE settings ADD COLUMN IF NOT EXISTS gemini_model VARCHAR(100) NOT NULL DEFAULT 'gemini-3.1-flash-lite';`,
		`ALTER TABLE settings ADD COLUMN IF NOT EXISTS ollama_model VARCHAR(100) NOT NULL DEFAULT 'gemma4:cloud';`,
		`ALTER TABLE settings ADD COLUMN IF NOT EXISTS enable_fact_checker BOOLEAN NOT NULL DEFAULT true;`,
		`ALTER TABLE settings ADD COLUMN IF NOT EXISTS fact_checker_provider VARCHAR(50) NOT NULL DEFAULT 'gemini';`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			log.Printf("[InitDB] Migration notice: %v", err)
		}
	}

	// Insert initial default settings if empty
	insertDefaultQuery := `
	INSERT INTO settings (id, context_provider, fact_checker_provider, gemini_api_key, magisterium_api_key, bible_api_key, gemini_model, ollama_model, context_model, context_instruction, enable_fact_checker)
	VALUES (1, 'ollama', 'gemini', '', '', '', 'gemini-3.1-flash-lite', 'gemma4:cloud', 'gemma4:cloud', $1, true)
	ON CONFLICT (id) DO NOTHING;`

	if _, err := db.Exec(insertDefaultQuery, defaultInstruction); err != nil {
		log.Printf("[InitDB] Error inserting default settings: %v", err)
	}
}

func GetSettings() Settings {
	if db == nil || db.Ping() != nil {
		return memorySettings
	}

	query := `SELECT context_provider, fact_checker_provider, gemini_api_key, magisterium_api_key, bible_api_key, gemini_model, ollama_model, context_model, context_instruction, enable_fact_checker FROM settings WHERE id = 1`
	row := db.QueryRow(query)

	var s Settings
	err := row.Scan(&s.ContextProvider, &s.FactCheckerProvider, &s.GeminiAPIKey, &s.MagisteriumAPIKey, &s.BibleAPIKey, &s.GeminiModel, &s.OllamaModel, &s.ContextModel, &s.ContextInstruction, &s.EnableFactChecker)
	if err != nil {
		log.Printf("[GetSettings] Warning scanning row: %v. Returning current memory settings.", err)
		return memorySettings
	}

	if s.ContextProvider == "" {
		s.ContextProvider = "ollama"
	}
	if s.FactCheckerProvider == "" {
		s.FactCheckerProvider = "gemini"
	}
	if s.GeminiModel == "" {
		s.GeminiModel = "gemini-3.1-flash-lite"
	}
	if s.OllamaModel == "" {
		s.OllamaModel = "gemma4:cloud"
	}
	if s.ContextProvider == "gemini" {
		s.ContextModel = s.GeminiModel
	} else {
		s.ContextModel = s.OllamaModel
	}
	if s.ContextInstruction == "" {
		s.ContextInstruction = defaultInstruction
	}

	memorySettings = s
	return s
}

func SaveSettings(s Settings) error {
	if s.ContextProvider == "" {
		s.ContextProvider = "ollama"
	}
	if s.FactCheckerProvider == "" {
		s.FactCheckerProvider = "gemini"
	}
	if s.GeminiModel == "" {
		s.GeminiModel = "gemini-3.1-flash-lite"
	}
	if s.OllamaModel == "" {
		s.OllamaModel = "gemma4:cloud"
	}
	if s.ContextProvider == "gemini" {
		s.ContextModel = s.GeminiModel
	} else {
		s.ContextModel = s.OllamaModel
	}
	if s.ContextInstruction == "" {
		s.ContextInstruction = defaultInstruction
	}

	memorySettings = s

	if db == nil || db.Ping() != nil {
		log.Println("[SaveSettings] Database unavailable, saved settings to in-memory store")
		return nil
	}

	query := `
	INSERT INTO settings (id, context_provider, fact_checker_provider, gemini_api_key, magisterium_api_key, bible_api_key, gemini_model, ollama_model, context_model, context_instruction, enable_fact_checker, updated_at)
	VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO UPDATE SET
		context_provider = EXCLUDED.context_provider,
		fact_checker_provider = EXCLUDED.fact_checker_provider,
		gemini_api_key = EXCLUDED.gemini_api_key,
		magisterium_api_key = EXCLUDED.magisterium_api_key,
		bible_api_key = EXCLUDED.bible_api_key,
		gemini_model = EXCLUDED.gemini_model,
		ollama_model = EXCLUDED.ollama_model,
		context_model = EXCLUDED.context_model,
		context_instruction = EXCLUDED.context_instruction,
		enable_fact_checker = EXCLUDED.enable_fact_checker,
		updated_at = CURRENT_TIMESTAMP;`

	_, err := db.Exec(query, s.ContextProvider, s.FactCheckerProvider, s.GeminiAPIKey, s.MagisteriumAPIKey, s.BibleAPIKey, s.GeminiModel, s.OllamaModel, s.ContextModel, s.ContextInstruction, s.EnableFactChecker)
	if err != nil {
		log.Printf("[SaveSettings] Error saving settings to postgres: %v", err)
		return err
	}

	log.Println("[SaveSettings] Settings updated successfully in PostgreSQL")
	return nil
}

type PublicConversation struct {
	ID         string      `json:"id"`
	AuthorName string      `json:"author_name"`
	Title      string      `json:"title"`
	Messages   interface{} `json:"messages"`
	CreatedAt  string      `json:"created_at"`
}

var memoryPublicConversations []PublicConversation

func SavePublicConversation(id, authorName, title string, messagesJSON []byte) error {
	var msgs interface{}
	json.Unmarshal(messagesJSON, &msgs)

	pub := PublicConversation{
		ID:         id,
		AuthorName: authorName,
		Title:      title,
		Messages:   msgs,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	memoryPublicConversations = append([]PublicConversation{pub}, memoryPublicConversations...)

	if db == nil || db.Ping() != nil {
		log.Println("[SavePublicConversation] DB unavailable, saved to memory")
		return nil
	}

	query := `
	INSERT INTO public_chat_conversations (id, author_name, title, messages, created_at)
	VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO UPDATE SET
		author_name = EXCLUDED.author_name,
		title = EXCLUDED.title,
		messages = EXCLUDED.messages;`

	_, err := db.Exec(query, id, authorName, title, string(messagesJSON))
	if err != nil {
		log.Printf("[SavePublicConversation] Error saving conversation: %v", err)
		return err
	}

	return nil
}

func GetPublicConversations() ([]PublicConversation, error) {
	if db == nil || db.Ping() != nil {
		return memoryPublicConversations, nil
	}

	query := `SELECT id, author_name, title, messages, created_at FROM public_chat_conversations ORDER BY created_at DESC LIMIT 50`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("[GetPublicConversations] Query error: %v", err)
		return memoryPublicConversations, nil
	}
	defer rows.Close()

	var result []PublicConversation
	for rows.Next() {
		var pc PublicConversation
		var msgsRaw string
		var t time.Time
		if err := rows.Scan(&pc.ID, &pc.AuthorName, &pc.Title, &msgsRaw, &t); err == nil {
			pc.CreatedAt = t.Format(time.RFC3339)
			json.Unmarshal([]byte(msgsRaw), &pc.Messages)
			result = append(result, pc)
		}
	}

	return result, nil
}
