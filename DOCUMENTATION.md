# Travel Planner — Technical Documentation (Generated)

> Repo: https://github.com/ashanda/travel-planner  
> Generated: 2026-01-13 (Asia/Colombo)

## 1) What this repository contains

From the repository root, the main components are:
- `trip-planner` — backend service (Go)  
- `trip-planner-web` — frontend web app (Vue/Nuxt-style)  
- `nginx` — reverse proxy / edge configuration  
- `docker-compose.yml` — local/prod orchestration  
- `.env.example` — environment variable template  

(These folders/files are visible in the repo root.) 

## 2) High-level architecture

This project is a typical 3-tier web system:

- **Client (Browser)**
  - Loads the web UI from the frontend container.
- **Frontend (trip-planner-web)**
  - Serves UI and calls backend APIs (usually under `/api/*`).
- **Backend API (trip-planner)**
  - Handles auth, itinerary generation, storage, integrations.
- **Storage**
  - Local DB (often SQLite for dev) or MySQL/Postgres for production (depends on config).
- **External APIs (optional)**
  - Maps/Places, AI provider, weather, etc.

## 3) Running the project

### 3.1 Run with Docker Compose (recommended)
1. Create env file:
   - Copy `.env.example` → `.env`
2. Build + start:
   ```bash
   docker compose up -d --build
   ```
3. Open:
   - Frontend: http://localhost (or your domain)
   - API: http://localhost/api (commonly)

> Note: exact ports/hosts depend on `docker-compose.yml` and `.env`.

### 3.2 Local dev (typical)
- Backend (Go):
  ```bash
  cd trip-planner
  go mod tidy
  go run .
  ```
- Frontend (Node/Nuxt):
  ```bash
  cd trip-planner-web
  npm install
  npm run dev
  ```

## 4) Environment variables

This repo includes a `.env.example` in the root, intended to document required env vars.

Common env variables for this stack usually include:
- `PORT` (API port)
- `JWT_SECRET`
- `COOKIE_DOMAIN`
- `DATABASE_URL` (or `DB_PATH` for SQLite)
- `BASE_URL` / `API_URL` for the frontend

**Action for you:** open `.env.example` and ensure every value is set in `.env` for your environment.

## 5) Request flow (typical)

1. User opens site (frontend served via Nginx).
2. Frontend calls API endpoints (e.g., `/api/auth/login`, `/api/plans`).
3. API validates JWT / cookie session.
4. API reads/writes database.
5. API returns JSON.
6. Frontend renders UI.

## 6) Suggested API surface (documentation template)

Because the repository’s file viewer is currently failing to render code in-browser, below is a **safe template** you can quickly align with your actual routes.

### Auth
- `POST /api/auth/register` — create user
- `POST /api/auth/login` — sign in
- `POST /api/auth/logout` — sign out
- `GET  /api/auth/me` — current user

### Trips / Plans
- `POST /api/plans` — generate a new plan
- `GET  /api/plans` — list saved plans
- `GET  /api/plans/:id` — get plan
- `DELETE /api/plans/:id` — delete plan

### Places
- `GET /api/places/search?q=...`
- `GET /api/places/:id`

## 7) Security checklist (recommended)

- Use **HTTP-only cookies** for sessions if possible.
- Ensure CORS settings are correct (restrict origins).
- Set secure cookie flags in production: `Secure`, `SameSite=Lax/Strict`.
- Keep secrets out of git (`.env` should not be committed).

## 8) Production deployment notes

Typical pattern:
- Nginx terminates TLS (or use Cloudflare / load balancer)
- Nginx routes:
  - `/` → frontend container
  - `/api` → backend container
- Persistent volumes:
  - DB volume
  - file uploads / cache (if used)

## 9) Next improvement: auto-generated docs

If you want fully accurate endpoint documentation:
- Add Swagger/OpenAPI for the Go API
  - go-swagger / swaggo (Gin-friendly)
- Add typed API client for frontend
  - generate TS client from OpenAPI

---

## Appendix A — Diagrams

This documentation comes with Mermaid diagrams (see `diagrams/` files).
