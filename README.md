<div align="center">

# Rekko

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)
![Vue Version](https://img.shields.io/badge/Vue-3.5-42b883?style=flat&logo=vue.js)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?style=flat&logo=typescript)
![TailwindCSS](https://img.shields.io/badge/Tailwind-4.1-38B2AC?style=flat&logo=tailwindcss)
![CockroachDB](https://img.shields.io/badge/CockroachDB-23.2-6933FF?style=flat&logo=cockroachlabs)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)

Stock recommendation engine powered by analyst ratings and price targets.

</div>

## Table of Contents

- [Project Overview](#project-overview)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start (Docker)](#quick-start-docker)
- [Local Development](#local-development)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
  - [Database Only](#database-only)
- [Deploy to Railway](#deploy-to-railway)
- [API Documentation](#api-documentation)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [Stock Endpoints](#stock-endpoints)
  - [Recommendation Endpoints](#recommendation-endpoints)
  - [Dashboard Endpoint](#dashboard-endpoint)
  - [Sync Endpoint](#sync-endpoint)
- [Project Structure](#project-structure)
- [Recommendation Algorithm](#recommendation-algorithm)
  - [Scoring Factors](#scoring-factors)
  - [Rating Values](#rating-values)
  - [Action Scores](#action-scores)
- [Environment Variables](#environment-variables)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## Project Overview

Rekko aggregates analyst recommendations from external sources and applies a weighted scoring algorithm to identify the most promising investment opportunities. The system provides:

- Real-time data synchronization from external stock APIs
- Intelligent ranking of stocks based on analyst upgrades, price targets, and brokerage actions
- A responsive web interface for browsing, searching, and filtering stock recommendations
- RESTful API for programmatic access to stock data and recommendations

## Tech Stack

| Layer | Technology |
|-------|------------|
| **Backend** | Go 1.24 with Gin framework |
| **Frontend** | Vue 3 + TypeScript + Pinia + Tailwind CSS 4 |
| **Database** | CockroachDB v23.2 |
| **UI Components** | shadcn-vue, Radix Vue, TanStack Table, TanStack Vue Query |
| **API Style** | RESTful with OpenAPI/Swagger documentation |
| **Build Tools** | Docker, Docker Compose |
| **Package Manager** | pnpm (frontend) |

## Features

- **Data Synchronization**: Fetch and store stock data from external APIs with pagination support
- **Search & Filter**: Search stocks by ticker, company name, or analyst action
- **Sorting**: Sort stocks by various criteria (ticker, company, action, target price, date)
- **Stock Details**: View detailed information for individual stocks
- **Smart Recommendations**: Intelligent scoring algorithm based on:
  - Rating upgrades (e.g., Hold to Buy)
  - Target price increases
  - Brokerage consensus and action types
  - Analyst credibility signals
- **Dashboard Analytics**: Aggregated statistics including total stocks, action distribution, top brokerages, and recent daily activity
- **Pagination**: Efficient handling of large datasets with server-side pagination
- **Responsive UI**: Mobile-friendly interface built with Tailwind CSS and shadcn-vue, featuring dark/light theme support and a collapsible sidebar navigation
- **Interactive API Documentation**: Auto-generated Swagger UI for exploring and testing all endpoints directly from the browser

## Prerequisites

### For Docker Deployment (Recommended)
- Docker 20.10+
- Docker Compose 2.0+

### For Local Development
- Go 1.24+
- Node.js 20+
- pnpm 8+ (`npm install -g pnpm`)
- CockroachDB (can run via Docker)

## Quick Start (Docker)

The fastest way to get the application running is with Docker Compose.

### 1. Clone the repository

```bash
git clone <repository-url>
cd rekko
```

### 2. Configure environment variables

```bash
cp .env.example .env
```

Edit the `.env` file and add your API token:

```bash
# Required: Add your KarenAI API token
KARENAI_AUTH_TOKEN=your_auth_token_here
```

### 3. Start all services

```bash
docker-compose up -d
```

This will start:
- **CockroachDB** on port 26257 (DB UI on port 8081)
- **Backend API** on port 8080
- **Frontend** on port 3000

### 4. Verify services are running

```bash
docker-compose ps
```

You should see all three services in "running" state.

### 5. Access the application

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/swagger/ |
| CockroachDB UI | http://localhost:8081 |

### 6. Populate the database

The database starts empty. You need to sync data from the external API.

**Option A: Via the UI**
- Open http://localhost:3000
- Click the "Sync Data" button in the header

**Option B: Via curl**
```bash
curl -X POST http://localhost:8080/api/v1/sync
```

Expected response:
```json
{
  "status": true,
  "message": "Sync completed successfully",
  "data": {
    "count": 150
  }
}
```

### 7. Stop services

```bash
docker-compose down
```

To also remove the database volume:
```bash
docker-compose down -v
```

## Local Development

### Backend Setup

1. **Start CockroachDB** (if not using Docker for the database):
   ```bash
   docker-compose up -d cockroachdb
   ```

2. **Navigate to the backend directory**:
   ```bash
   cd backend
   ```

3. **Install Go dependencies**:
   ```bash
   go mod download
   ```

4. **Set environment variables**:
   ```bash
   export DATABASE_URL="postgresql://root@localhost:26257/stockdb?sslmode=disable"
   export KARENAI_API_URL="https://api.karenai.click"
   export KARENAI_AUTH_TOKEN="your_token_here"
   export SERVER_PORT="8080"
   export GIN_MODE="debug"
   ```

5. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

   The backend will be available at http://localhost:8080

### Frontend Setup

1. **Navigate to the frontend directory**:
   ```bash
   cd frontend
   ```

2. **Install dependencies**:
   ```bash
   pnpm install
   ```

3. **Set the API configuration** by copying the example environment file:
   ```bash
   cp .env.example .env
   ```

4. **Run the development server**:
   ```bash
   pnpm dev
   ```

   The frontend will be available at http://localhost:5173

### Database Only

To run just CockroachDB for local development:

```bash
docker-compose up -d cockroachdb
```

Access the database:
- **SQL CLI**: `docker exec -it rekko-cockroachdb cockroach sql --insecure`
- **Admin UI**: http://localhost:8081

## Deploy to Railway

The project includes a unified Dockerfile at the repository root that produces a single container image with both the Go backend and the Vue frontend. This image is designed for platforms such as Railway, where each service runs in its own isolated container with a managed database alongside it.

### Architecture

The unified image compiles the Vue application into static files and embeds them into the Go binary's runtime. The backend serves the API on `/api/v1/*`, the Swagger documentation on `/swagger/*`, and all remaining routes fall through to the Vue single-page application. There is no need for Nginx, no cross-origin configuration between services, and no coordination of multiple containers — a single process handles everything.

```
Browser → https://your-app.up.railway.app
         │
         ├── /api/v1/*  → Go API handlers
         ├── /swagger/* → Swagger UI
         └── /*         → Vue SPA — static files and client-side routing
```

### Step-by-step

#### 1. Create the project

Sign in to [Railway](https://railway.app) and create a new project. Connect your GitHub repository when prompted.

#### 2. Add PostgreSQL

Click **"+ New"** → **"Database"** → **"PostgreSQL"**. Railway provisions the instance and generates the connection string automatically.

#### 3. Add the application service

Click **"+ New"** → **"GitHub Repo"** and select the Rekko repository. Configure the service under **Settings**:

| Setting | Value |
|---------|-------|
| Builder | Dockerfile |
| Dockerfile Path | `./Dockerfile` |
| Healthcheck Path | `/api/v1/health` |

#### 4. Configure environment variables

Navigate to the service's **Variables** tab and add the following entries:

| Variable | Value | Description |
|----------|-------|-------------|
| `DATABASE_URL` | `${{Postgres.DATABASE_URL}}` | Railway resolves this reference to the PostgreSQL connection string |
| `DB_DRIVER` | `postgres` | Switches the migration driver from CockroachDB to PostgreSQL |
| `KARENAI_AUTH_TOKEN` | Your token | Bearer token for the KarenAI external API |
| `FINNHUB_API_KEY` | Your key | Finnhub API key for real-time market data — optional |

The remaining variables — `PORT`, `STATIC_DIR`, `MIGRATIONS_PATH`, and `GIN_MODE` — are preconfigured in the Dockerfile and do not require manual overrides.

#### 5. Generate a public domain

Go to **Settings** → **Networking** → **"Generate Domain"**. Railway assigns a URL in the form `https://<project>.up.railway.app`.

#### 6. Populate the database

After the first successful deploy, the database will be empty. Trigger an initial data sync:

```bash
curl -X POST https://<your-domain>.up.railway.app/api/v1/sync
```

### How deployments work

Railway monitors the configured branch — `main` by default — and triggers an automatic rebuild on every push. The platform builds the Docker image, runs the healthcheck against `/api/v1/health`, and routes traffic to the new container only after it responds successfully. The previous container remains active during this transition, ensuring zero-downtime deployments.

### Unified image details

The root `Dockerfile` executes a three-stage build:

| Stage | Base Image | Output |
|-------|------------|--------|
| **Frontend build** | `node:20-alpine` | Compiled Vue assets in `/app/dist` |
| **Backend build** | `golang:1.24-alpine` | Statically linked Go binary with Swagger docs |
| **Runtime** | `alpine:latest` | Final image — approximately 83 MB — containing the binary, migrations, and static files |

The `STATIC_DIR` environment variable tells the Go server where to find the frontend assets. When this variable is empty — as it is by default in local development — the server operates in API-only mode and does not serve static files.

### Compatibility note

All changes introduced for Railway are fully backward-compatible with the existing Docker Compose workflow. The `DB_DRIVER` variable defaults to `cockroachdb`, the `STATIC_DIR` defaults to an empty string, and the `SERVER_PORT` variable takes precedence over the Railway-injected `PORT`. Running `docker-compose up` locally continues to work without modifications.

## API Documentation

The backend ships with auto-generated **Swagger/OpenAPI** documentation, built using [swag](https://github.com/swaggo/swag) annotations embedded in every handler. The Swagger specification is regenerated at build time via `swag init` during the Docker build, ensuring that the documentation always reflects the current state of the codebase.

Once the backend is running, the interactive Swagger UI is available at:

```
http://localhost:8080/swagger/
```

The raw OpenAPI specification is served as JSON at `/swagger/doc.json`, making it straightforward to import into tools such as Postman, Insomnia, or any OpenAPI-compatible client.

For local development without Docker, you can regenerate the documentation manually:

```bash
cd backend
go install github.com/swaggo/swag/cmd/swag@v1.16.6
swag init -g cmd/server/main.go -o docs --parseInternal
```

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Health Check

**GET** `/health`

Returns the current health status of the service. When the `Accept` header includes `text/html`, the endpoint renders an HTML status page instead of JSON.

```bash
curl http://localhost:8080/api/v1/health
```

Response:
```json
{
  "status": true,
  "message": "Service is running"
}
```

### Stock Endpoints

#### List Stocks (Paginated)

**GET** `/stocks`

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page (max: 100) |
| `search` | string | - | Search by ticker or company name |
| `ticker` | string | - | Filter by exact ticker symbol |
| `action` | string | - | Filter by action type |
| `sortBy` | string | created_at | Sort field: `ticker`, `company`, `action`, `targetTo`, `createdAt` |
| `sortOrder` | string | desc | Sort order: `asc`, `desc` |

```bash
# Get first page of stocks
curl "http://localhost:8080/api/v1/stocks"

# Search for Apple stocks
curl "http://localhost:8080/api/v1/stocks?search=AAPL"

# Get upgraded stocks sorted by target price
curl "http://localhost:8080/api/v1/stocks?action=upgraded&sortBy=targetTo&sortOrder=desc"

# Pagination example
curl "http://localhost:8080/api/v1/stocks?page=2&limit=50"
```

Response:
```json
{
  "status": true,
  "message": "Stocks retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "ticker": "AAPL",
      "company": "Apple Inc",
      "brokerage": "Morgan Stanley",
      "action": "upgraded",
      "ratingFrom": "Hold",
      "ratingTo": "Buy",
      "targetFrom": 180.00,
      "targetTo": 220.00,
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "current_page": 1,
      "per_page": 20,
      "total_items": 150,
      "total_pages": 8,
      "has_next": true
    }
  }
}
```

#### Get Stock by ID

**GET** `/stocks/:id`

```bash
curl http://localhost:8080/api/v1/stocks/550e8400-e29b-41d4-a716-446655440000
```

#### Get Stocks by Ticker

**GET** `/stocks/ticker/:ticker`

```bash
curl http://localhost:8080/api/v1/stocks/ticker/AAPL
```

Response:
```json
{
  "status": true,
  "message": "Stocks retrieved successfully",
  "data": [
    {
      "id": "...",
      "ticker": "AAPL",
      "company": "Apple Inc",
      "brokerage": "Morgan Stanley",
      "action": "upgraded",
      ...
    },
    {
      "id": "...",
      "ticker": "AAPL",
      "company": "Apple Inc",
      "brokerage": "Goldman Sachs",
      "action": "maintained",
      ...
    }
  ]
}
```

#### Get Available Actions

**GET** `/stocks/actions`

```bash
curl http://localhost:8080/api/v1/stocks/actions
```

Response:
```json
{
  "status": true,
  "message": "Actions retrieved successfully",
  "data": [
    "upgraded",
    "downgraded",
    "initiated",
    "maintained",
    "reiterated",
    "target raised",
    "target lowered"
  ]
}
```

### Recommendation Endpoints

#### Get Top Recommendations

**GET** `/recommendations`

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | int | 50 | Maximum number of recommendations to return |
| `search` | string | - | Filter by ticker or company name |

```bash
# Get top recommendations
curl http://localhost:8080/api/v1/recommendations

# Get top 5 recommendations
curl "http://localhost:8080/api/v1/recommendations?limit=5"

# Search within recommendations
curl "http://localhost:8080/api/v1/recommendations?search=NVDA"
```

Response:
```json
{
  "status": true,
  "message": "Recommendations retrieved successfully",
  "data": [
    {
      "stock": {
        "id": "...",
        "ticker": "NVDA",
        "company": "NVIDIA Corporation",
        "brokerage": "Bank of America",
        "action": "upgraded",
        "ratingFrom": "Neutral",
        "ratingTo": "Buy",
        "targetFrom": 450.00,
        "targetTo": 600.00
      },
      "score": 85.5,
      "reasons": [
        "Rating upgraded from Neutral to Buy",
        "Target price increased 33.3% to $600.00",
        "upgraded by Bank of America"
      ],
      "upsidePotential": 33.33,
      "analystCount": 4,
      "marketData": null
    }
  ]
}
```

#### Get Best Single Recommendation

**GET** `/recommendations/top`

```bash
curl http://localhost:8080/api/v1/recommendations/top
```

Response:
```json
{
  "status": true,
  "message": "Top recommendation retrieved successfully",
  "data": {
    "stock": {
      "id": "...",
      "ticker": "NVDA",
      "company": "NVIDIA Corporation",
      ...
    },
    "score": 85.5,
    "reasons": [...],
    "upsidePotential": 33.33,
    "analystCount": 4,
    "marketData": null
  }
}
```

### Dashboard Endpoint

#### Get Dashboard Statistics

**GET** `/dashboard/stats`

Returns aggregated statistics including total stocks, action distribution across the dataset, the top brokerages by volume, and daily activity over the last 30 days.

```bash
curl http://localhost:8080/api/v1/dashboard/stats
```

Response:
```json
{
  "status": true,
  "message": "Dashboard stats retrieved successfully",
  "data": {
    "totalStocks": 150,
    "actionDistribution": [
      { "action": "upgraded", "count": 45 },
      { "action": "maintained", "count": 38 },
      { "action": "initiated", "count": 27 }
    ],
    "brokerageDistribution": [
      { "brokerage": "Morgan Stanley", "count": 22 },
      { "brokerage": "Goldman Sachs", "count": 18 }
    ],
    "recentActivity": [
      { "date": "2024-01-15", "count": 12 },
      { "date": "2024-01-14", "count": 8 }
    ]
  }
}
```

### Sync Endpoint

#### Trigger Data Sync

**POST** `/sync`

Fetches the latest stock data from the external API and updates the database.

```bash
curl -X POST http://localhost:8080/api/v1/sync
```

Response:
```json
{
  "status": true,
  "message": "Sync completed successfully",
  "data": {
    "count": 150
  }
}
```

**Note**: The sync process:
1. Connects to the KarenAI API
2. Fetches all stock recommendations (handles pagination automatically)
3. Upserts records into the database (updates existing, inserts new)
4. Returns the count of processed records

## Recommendation Algorithm

The recommendation engine employs a weighted multi-factor scoring model that combines analyst sentiment with real-time market data. Each ticker receives a composite score on a 0–10 scale, derived from up to eight distinct factors when market data is available, or five analyst-based factors as a fallback.

### Architecture Overview

The system operates in two modes depending on the availability of external market data from **Finnhub**:

- **Enriched mode** — When the Finnhub API key is configured, the engine fetches live quotes and company profiles for every ticker under evaluation. This enables three additional scoring dimensions that ground analyst opinions in actual market conditions.
- **Fallback mode** — When market data is unavailable — either because the API key is not set or because a specific ticker lacks coverage on Finnhub — the engine relies exclusively on analyst-derived signals, preserving full functionality without external dependencies.

### Scoring Factors

#### With Market Data — 8 Factors

| Factor | Weight | Source | Description |
|--------|--------|--------|-------------|
| **Rating Upgrade** | 15% | Analyst data | Positive transitions in analyst ratings, e.g., Hold to Buy |
| **Target Price Increase** | 10% | Analyst data | Percentage increase in the analyst's price target |
| **Action Type** | 15% | Analyst data | Base score derived from the nature of the analyst action |
| **Consensus** | 20% | Computed | Proportion of brokerages with a bullish stance on the ticker |
| **Momentum** | 10% | Computed | Time-weighted accumulation of recent analyst signals with exponential decay |
| **Real Upside** | 15% | Finnhub | Gap between the average analyst target and the current market price |
| **Market Cap** | 10% | Finnhub | Company size tier — larger capitalizations reflect lower risk and more reliable coverage |
| **Price Trend** | 5% | Finnhub | Whether today's price movement confirms or contradicts the analyst consensus |

#### Without Market Data — 5 Factors, Fallback

| Factor | Weight | Description |
|--------|--------|-------------|
| **Rating Upgrade** | 20% | Positive transitions in analyst ratings |
| **Target Price Increase** | 20% | Percentage increase in analyst price targets |
| **Action Type** | 20% | Base score from the nature of the analyst action |
| **Consensus** | 25% | Proportion of bullish brokerages |
| **Momentum** | 15% | Time-weighted analyst signals |

This dual-weight architecture ensures that tickers with market data benefit from richer context, while those without it are scored fairly under the original model — no ticker is penalized for the absence of external data.

### Rating Values

The system converts analyst ratings into a numerical scale from 1 to 5, enabling precise measurement of rating transitions:

| Rating | Value |
|--------|-------|
| Strong Sell | 1 |
| Sell, Underweight | 2 |
| Hold, Neutral, Equal-Weight, Market Perform, Sector Perform | 3 |
| Buy, Overweight, Outperform | 4 |
| Strong Buy, Top Pick | 5 |

### Action Scores

Each analyst action maps to a base score that reflects the strength of the signal:

| Action | Score |
|--------|-------|
| Upgraded | 100 |
| Initiated | 80 |
| Target Raised | 70 |
| Reiterated | 60 |
| Maintained | 50 |
| Target Lowered | 30 |
| Downgraded | 20 |

### Factor Details

**Rating Upgrade** — Measures the delta between the previous and current rating on the 1–5 scale, normalized to a 0–100 range. A transition from Hold to Buy yields a higher score than a reiteration at Buy. Tickers that already hold a strong rating — Buy or above — receive a baseline score of 50 even without an upgrade.

**Target Price Increase** — Captures the percentage change in the analyst's price target, capped at 100 points. A 50% increase in the target translates directly to a score of 50.

**Action Type** — Assigns the base score from the action table above. This factor rewards decisive bullish actions — upgrades and initiations — over neutral or bearish ones.

**Consensus** — Calculates the percentage of distinct brokerages that have taken a bullish action on the ticker. When fewer than three brokerages cover the stock, the score is discounted proportionally to reflect lower statistical confidence.

**Momentum** — Applies an exponential time-decay function over a 30-day window to weight recent analyst signals more heavily than older ones. Bullish actions contribute positively, while bearish actions reduce the accumulated signal at half the rate. A saturation function prevents any single ticker from achieving a disproportionately high momentum score.

**Real Upside** — Computes the actual upside potential by comparing the average analyst target price against the live market price from Finnhub. This is the most impactful enrichment: it transforms abstract analyst targets into a concrete percentage of potential gain. The score normalizes linearly, with 50% or more of upside mapping to the maximum score of 100.

**Market Cap** — Assigns a tier-based score that reflects company size and the reliability of analyst coverage:

| Market Cap | Tier | Score |
|------------|------|-------|
| $10B and above | Large-cap | 100 |
| $2B – $10B | Mid-cap | 75 |
| $300M – $2B | Small-cap | 50 |
| Below $300M | Micro-cap | 25 |

**Price Trend** — Evaluates whether the current day's price movement aligns with the prevailing analyst sentiment. When the majority of analysts are bullish and the price is trending upward, the score rises above the neutral baseline of 50. Conversely, a bearish price movement against bullish analyst consensus results in a moderate reduction, signaling a potential divergence worth noting.

### Upside Potential Calculation

The upside potential displayed to the user varies based on data availability:

- **With market data** — Calculated as the gap between the average analyst target and the current Finnhub price: this represents the real upside from the present trading level.
- **Without market data** — Derived from the maximum percentage change between `targetFrom` and `targetTo` across the ticker's analyst reports.

### Market Data Integration

Real-time market data is sourced from the **Finnhub API**, which provides live quotes and company profiles. The integration is designed with the following considerations:

- **In-memory cache** with a 15-minute TTL per ticker, balancing data freshness against API rate limits.
- **Parallel fetching** with a concurrency limit of 10 simultaneous requests, ensuring fast batch processing without overwhelming the external service.
- **Graceful degradation** — if Finnhub returns no data for a specific ticker, or if the API is unreachable, the system falls back to analyst-only scoring for that ticker without affecting others.

The final composite score is the weighted sum of all applicable factors, divided by 10 to produce the 0–10 scale. Recommendations are ranked by descending score.

## Environment Variables

### Backend

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | No | `postgresql://root@localhost:26257/stockdb?sslmode=disable` | PostgreSQL/CockroachDB connection string |
| `KARENAI_API_URL` | No | `https://api.karenai.click` | External API base URL |
| `KARENAI_AUTH_TOKEN` | Yes | - | Bearer token for external API |
| `FINNHUB_API_KEY` | No | - | Finnhub API key for real-time market data enrichment |
| `SERVER_PORT` | No | `8080` | Backend server port — falls back to `PORT` if not set, for Railway compatibility |
| `GIN_MODE` | No | `debug` | Gin framework mode — `debug` or `release` |
| `MIGRATIONS_PATH` | No | `./migrations` | Path to the SQL migration files directory |
| `DB_DRIVER` | No | `cockroachdb` | Database migration driver — `cockroachdb` for local, `postgres` for Railway |
| `STATIC_DIR` | No | - | Path to the frontend static files directory — when set, the backend serves the Vue SPA |

### Frontend

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `VITE_SERVER_API_URL` | No | `""` | Base URL for the backend API — leave empty when the frontend is served from the same origin as the backend |
| `VITE_SERVER_API_PREFIX` | Yes | - | API path prefix appended to the base URL, e.g. `/api` |
| `VITE_SERVER_API_TIMEOUT` | No | `5000` | HTTP request timeout in milliseconds |

### Example `.env` file

```bash
# Database
DATABASE_URL=postgresql://root@localhost:26257/stockdb?sslmode=disable

# External API
KARENAI_API_URL=https://api.karenai.click
KARENAI_AUTH_TOKEN=your_auth_token_here

# Market Data Enrichment
FINNHUB_API_KEY=your_finnhub_api_key_here

# Server
SERVER_PORT=8080
GIN_MODE=debug
```

The frontend uses its own `.env` file inside the `frontend/` directory. Refer to `frontend/.env.example` for the template:

```bash
VITE_SERVER_API_URL=http://localhost:3000
VITE_SERVER_API_PREFIX=/api
VITE_SERVER_API_TIMEOUT=5000
```

## Testing

### Backend Testing

Run Go tests:

```bash
cd backend
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

Run with coverage:

```bash
go test -cover ./...
```

### Frontend Testing

Run Vue/TypeScript tests (if configured):

```bash
cd frontend
pnpm test
```

### Manual API Testing

Test the API endpoints using curl, or use the interactive [Swagger UI](http://localhost:8080/swagger/) for a browser-based experience:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Sync data
curl -X POST http://localhost:8080/api/v1/sync

# Get stocks
curl "http://localhost:8080/api/v1/stocks?limit=5"

# Get recommendations
curl http://localhost:8080/api/v1/recommendations

# Get best recommendation
curl http://localhost:8080/api/v1/recommendations/top

# Get dashboard statistics
curl http://localhost:8080/api/v1/dashboard/stats
```

## Troubleshooting

### Common Issues

#### 1. Docker containers won't start

**Symptom**: `docker-compose up` fails or containers exit immediately.

**Solutions**:
- Check if ports 3000, 8080, 8081, or 26257 are already in use:
  ```bash
  lsof -i :8080
  ```
- Ensure Docker has sufficient resources (memory, disk space)
- Check Docker logs:
  ```bash
  docker-compose logs backend
  docker-compose logs cockroachdb
  ```

#### 2. Backend can't connect to database

**Symptom**: Backend logs show "connection refused" errors.

**Solutions**:
- Wait for CockroachDB to be fully healthy:
  ```bash
  docker-compose logs cockroachdb | grep "health"
  ```
- Verify the DATABASE_URL is correct
- Check if the database exists:
  ```bash
  docker exec -it rekko-cockroachdb cockroach sql --insecure -e "SHOW DATABASES"
  ```

#### 3. Sync returns 401 Unauthorized

**Symptom**: `POST /api/v1/sync` returns authentication error.

**Solutions**:
- Verify `KARENAI_AUTH_TOKEN` is set in your `.env` file
- Ensure the token is valid and not expired
- Check if the token was loaded:
  ```bash
  docker-compose exec backend env | grep KARENAI
  ```

#### 4. Frontend shows "Network Error"

**Symptom**: UI can't fetch data, browser console shows CORS or network errors.

**Solutions**:
- Verify the backend is running:
  ```bash
  curl http://localhost:8080/api/v1/health
  ```
- Check `VITE_SERVER_API_URL` and `VITE_SERVER_API_PREFIX` are set correctly in `frontend/.env`
- Ensure CORS is properly configured in the backend
- If using Docker, ensure all services are on the same network

#### 5. Empty recommendations list

**Symptom**: Recommendations endpoint returns empty array.

**Solutions**:
- Ensure data has been synced:
  ```bash
  curl -X POST http://localhost:8080/api/v1/sync
  ```
- Check if stocks exist:
  ```bash
  curl "http://localhost:8080/api/v1/stocks?limit=1"
  ```

#### 6. Database data persists after `docker-compose down`

**Symptom**: Old data appears after restarting containers.

**Solution**:
- Remove volumes when stopping:
  ```bash
  docker-compose down -v
  ```

#### 7. pnpm install fails

**Symptom**: Frontend dependencies won't install.

**Solutions**:
- Ensure pnpm is installed globally:
  ```bash
  npm install -g pnpm
  ```
- Clear pnpm cache:
  ```bash
  pnpm store prune
  ```
- Delete `node_modules` and `pnpm-lock.yaml`, then retry

### Getting Help

If you encounter issues not covered here:

1. Check the Docker logs for specific error messages
2. Verify all environment variables are set correctly
3. Ensure all prerequisites are installed with the correct versions
4. Try rebuilding containers: `docker-compose build --no-cache`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
