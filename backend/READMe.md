# Devotion API

A Go-based API that generates daily Catholic devotionals using **Ollama** and the **Google Agent Development Kit (ADK)**. It fetches daily mass readings and provides historical and theological context for each reading by orchestrating AI agents with specialized tools.

## Technology Stack

- **Language:** Go (Golang) 1.25+
- **AI Orchestration:** [Google Agent Development Kit (ADK)](https://github.com/google/adk)
- **AI Model:** Ollama-based models (e.g., `qwen2.5:4b`, `llama3`, etc.)
- **Database/Cache:** Redis
- **Containerization:** Docker

## Prerequisites

Before running the application, ensure you have the following installed:

1.  **Go**: [Download and install Go](https://go.dev/dl/) (version 1.25 or newer).
2.  **Ollama**: [Install Ollama](https://ollama.com/) and download your preferred model (e.g., `ollama pull qwen2.5:4b`).
3.  **Redis**: [Install Redis](https://redis.io/docs/getting-started/) and ensure it allows connections (default port `6379`).
4.  **Bible API Key**: Obtain an API key from [API.Bible](https://scripture.api.bible/).
5.  **Docker** (Optional): If you prefer to run the application in a container.

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository_url>
cd devotion-api
```

### 2. Environment Configuration

Create a `.env` file in the root directory with the following variables:

```env
# HTTP Server Port
PORT=1024

# Authentication Token
API_TOKEN=your_strong_api_token_here
BIBLE_ID=your_target_bible_id

# CORS Policy
CORS_ALLOWED_ORIGIN=*

# PostgreSQL Database Connection
POSTGRES_HOST=host.docker.internal
POSTGRES_PORT=5432
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=devotion

# Redis Cache Address
REDIS_URL=host.docker.internal:6379

# Ollama Endpoint
OLLAMA_URL=http://host.docker.internal:11434
```

### 3. Running Locally

Ensure your Redis server and Ollama are running, then execute:

```bash
go run .
```

The server will start on port `1024` (or the port specified in your `.env`).

API Endpoint: `http://localhost:1024/devotion`

### 4. Running with Docker

You can also build and run the application using Docker.

**Build the image:**

```bash
docker build -t devotion-api .
```

**Run the container:**

```bash
docker run -p 1024:1024 --env-file .env -e TZ=Asia/Manila devotion-api

```

### 5. Running with Docker Compose

To start the application and Redis together:

```bash
docker-compose up --build
```

The API will be available at `http://localhost:1024`.

To stop the services:

```bash
docker-compose down
```

## API Usage

### Devotion Endpoint
**Endpoint:** `GET /devotion`

**Response Example:**
```json
{
  "first_reading": {
    "citation": "Acts 2:1-11",
    "text": "...",
    "context": "..."
  },
  "responsorial_psalm": {
    "citation": "Psalm 104",
    "text": "..."
  },
  "second_reading": {
    "citation": "1 Corinthians 12:3b-7, 12-13",
    "text": "...",
    "context": "..."
  },
  "gospel": {
    "citation": "John 20:19-23",
    "text": "...",
    "context": "..."
  }
}
```

### Bible Passage Endpoint
**Endpoint:** `GET /bible?passage=John+3:16`

**Response Example:**
```json
{
  "content": "For God so loved the world...",
  "copyright": "..."
}
```