# Complete Monorepo Setup & Development Guide

This document outlines the detailed steps required to configure, develop, test, and deploy both the **Frontend** web app and **Backend** API service in the **Devotion** repository.

---

## Workspace Layout

```
devotion/
├── frontend/             # Vue 3 + TypeScript + Vite (Cloudflare Pages)
│   ├── src/              # App components, routing, & Admin view (/admin)
│   ├── package.json      # Dependencies and scripts
│   └── vite.config.ts    # Vite configuration
├── backend/              # Go 1.25 REST API
│   ├── main.go           # Server startup & routing setup
│   ├── db.go             # PostgreSQL connection & settings persistence
│   ├── admin_handler.go  # Admin auth & settings endpoints
│   ├── bible_context.go  # Dynamic Bible context generator (Ollama / Gemini)
│   ├── go.mod            # Go modules definition
│   ├── Dockerfile        # Container build definition
│   └── docker-compose.yml# Multi-container setup (API + Redis + PostgreSQL)
├── README.md             # Repository summary and quickstart
└── SETUP.md              # Detailed step-by-step setup guide (this file)
```

---

## 1. Environment Requirements

Ensure the following runtimes and services are available on your development host:

- **Go**: `v1.25.0+`
- **Node.js**: `v18.0+` & `npm`
- **PostgreSQL**: Running locally on `localhost:5432` or via Docker container
- **Redis Server**: Local service on `localhost:6379` or via Docker container
- **Ollama**: Running locally at `http://localhost:11434`

---

## 2. Environment Configuration File (.env)

The `.env` file contains system connection details (database, redis, ollama url). AI configuration settings (Gemini API Key, active provider, model, and context instruction) are stored in PostgreSQL and managed directly via the `/admin` settings page.

### Backend Environment File (`.env` or `backend/.env`)
```env
# HTTP Server Port
PORT=8080

# Authentication Token
API_TOKEN=your_secure_api_token

# CORS Policy
CORS_ALLOWED_ORIGIN=*

# PostgreSQL Database Connection
DATABASE_URL=postgres://postgres:postgres@localhost:5432/devotion?sslmode=disable
# Alternatively set individual PG host/user/password/db:
# POSTGRES_HOST=localhost
# POSTGRES_PORT=5432
# POSTGRES_USER=postgres
# POSTGRES_PASSWORD=postgres
# POSTGRES_DB=devotion

# Redis Cache Address
REDIS_URL=localhost:6379

# Ollama Endpoint
OLLAMA_URL=http://localhost:11434
```

### Frontend Environment File (`frontend/.env`)
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TOKEN=your_secure_api_token
```

---

## 3. Admin Settings Page (`/admin`)

The Admin settings page allows live customization of AI generation settings.

1. Navigate to `http://localhost:5173/#/admin` (or `/admin`).
2. Log in with credentials:
   - **Username**: `admin`
   - **Password**: `admin`
3. Features configured in Admin Settings:
   - **AI Provider**: Switch between **Ollama (Local)** and **Google Gemini**.
   - **Gemini API Key**: Input & store `GEMINI_API_KEY` in PostgreSQL DB (read exclusively from database).
   - **Bible API Key**: Input & store `BIBLE_API_KEY` in PostgreSQL DB (read exclusively from database).
   - **Model Selection**: Select from available Ollama or Gemini models.
   - **Context Instruction**: Modify prompt instructions given to the AI, with support for `{{citation}}` and `{{passage_text}}` dynamic placeholders.

---

## 4. Running Backend Locally

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

### Running Backend with Docker Compose (API + Redis + PostgreSQL)
From `backend/`:
```bash
docker-compose up --build
```

---

## 5. Running Frontend Locally

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
5. Access Admin Settings at `http://localhost:5173/#/admin`.

---

## 6. Build & Deployment Steps

### Frontend Deployment (Cloudflare Pages)

Before deploying, ensure you configure the required environment variables.

> [!IMPORTANT]
> **Production Environment Variables**:
> - **GitHub Actions**: Add `VITE_API_TOKEN` (matching the backend `API_TOKEN`) and `GEMINI_API_KEY` to your repository **Secrets**, and `VITE_API_BASE_URL` to your repository **Variables** (`Settings > Secrets and variables > Actions`).
> - **Cloudflare Pages**: Add `VITE_API_TOKEN` and `VITE_API_BASE_URL` in the Cloudflare Pages project console under **Settings > Environment variables**.

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
