# Devotion Monorepo

Welcome to **Devotion**! This monorepo merges the frontend web application and backend API into a single workspace.

## Directory Structure

```
devotion/
├── frontend/        # Vue 3 + TypeScript + Vite + Cloudflare Pages (Wrangler)
├── backend/         # Go 1.25 REST API + Google ADK + Ollama / Gemini + Redis
├── README.md        # Monorepo setup guide & documentation
└── SETUP.md         # Detailed step-by-step setup guide
```

---

## 🛠️ Prerequisites

Make sure you have the following installed on your machine before getting started:

| Tool | Version / Details | Purpose |
| --- | --- | --- |
| **Node.js** | `v18.x` or newer | Frontend development & package management |
| **Go** | `v1.25.0` or newer | Backend API service |
| **Redis** | Default port `6379` | Response caching |
| **Ollama** | Local server (`http://localhost:11434`) | Local AI model generation (e.g. `qwen2.5:4b`) |
| **Google Gemini API Key** | Optional / Env | Google ADK AI orchestration |
| **API.Bible Key** | Optional / Env | Bible text lookup integration |

---

## ⚙️ Environment Configuration

### 1. Root / Backend Environment Variables
Create a `.env` file in the root project directory (or inside `backend/.env`):

```env
# Server Config
PORT=8080
API_TOKEN=your_secure_api_token
CORS_ALLOWED_ORIGIN=*

# AI & Storage Services
OLLAMA_URL=http://localhost:11434
REDIS_URL=localhost:6379
GEMINI_API_KEY=your_gemini_api_key

# External Scripture API
BIBLE_API_KEY=your_bible_api_key
BIBLE_ID=your_target_bible_id
```

### 2. Frontend Environment Variables
Create a `.env` or `.env.local` inside `frontend/`:

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TOKEN=your_secure_api_token
```

---

## 🚀 Quick Start Guide

### Step 1: Start the Backend Service

```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod download

# (Optional) Verify backend tests
go test ./...

# Run the Go API server
go run .
```
The API server will launch at `http://localhost:8080`.

*(Alternatively, run backend + Redis via Docker Compose from `backend/`:)*
```bash
cd backend
docker-compose up --build
```

---

### Step 2: Start the Frontend Application

Open a new terminal session and execute:

```bash
# Navigate to frontend directory
cd frontend

# Install npm packages
npm install

# Start Vite dev server
npm run dev
```
The frontend dev server will be available at `http://localhost:5173`.

---

## 📚 Related Documentation

- **Backend & Auth Guide**: See [backend/FRONTEND_AUTH_GUIDE.md](file:///Users/raenard/Documents/devotion/backend/FRONTEND_AUTH_GUIDE.md) for full endpoint specifications, response schemas, and Bearer token instructions.
- **Frontend Cloudflare Deployment**: See [frontend/README.md](file:///Users/raenard/Documents/devotion/frontend/README.md) for deployment scripts using `wrangler`.
- **Step-by-step Setup**: See [SETUP.md](file:///Users/raenard/Documents/devotion/SETUP.md) for complete detailed instructions.
