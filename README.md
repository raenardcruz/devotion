# Devotion Monorepo

Welcome to **Devotion**! This monorepo merges the frontend web application and backend API into a single workspace.

## Directory Structure

```
devotion/
├── frontend/        # Vue 3 + TypeScript + Vite + Cloudflare Pages (Wrangler)
├── backend/         # Go 1.25 REST API + Google ADK + Ollama / Gemini + Redis + PostgreSQL
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
| **PostgreSQL** | Default port `5432` | Storage for system settings & AI configurations |
| **Redis** | Default port `6379` | Response caching |
| **Ollama** | Local server (`http://localhost:11434`) | Local AI model generation (e.g. `gemma4:cloud` / `qwen2.5:4b`) |

---

## ⚙️ Environment Configuration

Environment files now strictly contain system connections (DB connection, Redis, Ollama URL). All AI provider selections, Gemini API Key, model dropdown choices, and prompt instructions are stored in PostgreSQL and configured via the **Admin Settings Page (`/admin`)**.

### 1. Root / Backend Environment Variables (`.env`)
Create a `.env` file in `backend/.env` or the workspace root:

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

### 2. Frontend Environment Variables (`frontend/.env`)
Create a `.env` file inside `frontend/`:

```env
VITE_API_BASE_URL=http://localhost:1024
VITE_API_TOKEN=your_strong_api_token_here
```

---

## 🔑 Admin Settings Page (`/admin`)

The Admin settings dashboard allows live customization of AI providers and context generation.

- **URL**: `http://localhost:5173/#/admin`
- **Username**: `admin`
- **Password**: `admin`

**Settings Available:**
1. **AI Generation Provider**: Toggle between **Ollama (Local)** and **Google Gemini**.
2. **Gemini API Key**: Saved to & read strictly from PostgreSQL database.
3. **Bible API Key**: Saved to & read strictly from PostgreSQL database.
4. **Model Selection**: Dynamic dropdown populated from Ollama tags or Gemini model options.
5. **Bible Context Instruction**: Custom prompt template for scripture context generation.

---

## 🚀 Quick Start Guide

### Step 1: Start the Backend Service

```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod download

# Verify backend tests
go test ./...

# Run the Go API server
go run .
```
The API server will launch at `http://localhost:1024`.

*(Alternatively, run backend + Redis + PostgreSQL via Docker Compose:)*
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
The frontend dev server will be available at `http://localhost:5173`. Access `/admin` to open the Admin Settings portal.

---

## 📚 Related Documentation

- **Step-by-step Setup**: See [SETUP.md](file:///Users/raenard/Documents/devotion/SETUP.md) for complete detailed instructions.
- **Backend & Auth Guide**: See [backend/FRONTEND_AUTH_GUIDE.md](file:///Users/raenard/Documents/devotion/backend/FRONTEND_AUTH_GUIDE.md) for endpoint specifications and Bearer token instructions.
