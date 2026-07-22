# Complete Monorepo Setup & Development Guide

This document outlines the detailed steps required to configure, develop, test, and deploy both the **Frontend** web app and **Backend** API service in the **Devotion** repository.

---

## Workspace Layout

```
devotion/
├── frontend/             # Vue 3 + TypeScript + Vite (Cloudflare Pages)
│   ├── src/              # App components & routing
│   ├── package.json      # Dependencies and scripts
│   └── vite.config.ts    # Vite configuration
├── backend/              # Go 1.25 REST API
│   ├── cmd/              # Utility commands
│   ├── main.go           # Server startup & routing setup
│   ├── go.mod            # Go modules definition
│   ├── Dockerfile        # Container build definition
│   └── docker-compose.yml# Multi-container setup (API + Redis)
├── README.md             # Repository summary and quickstart
└── SETUP.md              # Detailed step-by-step setup guide (this file)
```

---

## 1. Environment Requirements

Ensure the following runtimes and services are available on your development host:

- **Go**: `v1.25.0+`
- **Node.js**: `v18.0+` & `npm`
- **Redis Server**: Local service on `localhost:6379` or Docker container
- **Ollama**: Running locally at `http://localhost:11434` (Model: `qwen2.5:4b` or equivalent)

---

## 2. Environment Configuration File (.env)

### Backend Environment File (`.env` or `backend/.env`)
```env
# HTTP Server Port
PORT=8080

# Authentication Token
API_TOKEN=your_secure_api_token

# CORS Policy
CORS_ALLOWED_ORIGIN=*

# Redis Cache Address
REDIS_URL=localhost:6379

# Ollama Endpoint
OLLAMA_URL=http://localhost:11434

# Google Gemini API Key (Required for Google ADK AI orchestration)
GEMINI_API_KEY=your_gemini_api_key

# Bible API Settings
BIBLE_API_KEY=your_bible_api_key
BIBLE_ID=your_target_bible_id
```

### Frontend Environment File (`frontend/.env`)
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TOKEN=your_secure_api_token
```

---

## 3. Running Backend Locally

1. Open terminal and enter `backend/`:
   ```bash
   cd backend
   ```
2. Install Go dependencies:
   ```bash
   go mod download
   ```
3. Run backend unit tests:
   ```bash
   go test ./...
   ```
4. Start the server:
   ```bash
   go run .
   ```
5. Test API connection:
   ```bash
   curl -H "Authorization: Bearer your_secure_api_token" http://localhost:8080/devotion
   ```

### Running Backend with Docker Compose
From `backend/`:
```bash
docker-compose up --build
```

---

## 4. Running Frontend Locally

1. Open terminal and enter `frontend/`:
   ```bash
   cd frontend
   ```
2. Install Node dependencies:
   ```bash
   npm install
   ```
3. Start the dev server:
   ```bash
   npm run dev
   ```
4. Access the web app at `http://localhost:5173`.

---

## 5. Build & Deployment Steps

### Frontend Deployment (Cloudflare Pages)
```bash
cd frontend
npm run build
npx wrangler pages deploy dist --project-name devotion
```

### Backend Deployment (Docker Container)
```bash
cd backend
docker build -t devotion-api .
docker run -d -p 8080:8080 --env-file .env devotion-api
```
