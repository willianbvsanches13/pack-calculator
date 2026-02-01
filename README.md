# Pack Calculator

Go API (Gin) + React frontend to calculate the optimal number of packs needed for an order.

## Tech Stack

- **Backend:** Go 1.24, Gin-gonic v1.10
- **Frontend:** React 19, Vite 7, Tailwind CSS 4, React Query 5
- **Web Server:** Nginx (serves static files + reverse proxy)
- **Container:** Docker with separate frontend/backend services

## Architecture

```
                    ┌─────────────────┐
                    │     Client      │
                    └────────┬────────┘
                             │ :80
                    ┌────────▼────────┐
                    │      Nginx      │
                    │   (frontend)    │
                    │  - Static files │
                    │  - Gzip         │
                    │  - Cache        │
                    └────────┬────────┘
                             │ /api, /health
                    ┌────────▼────────┐
                    │    Go (API)     │
                    │   (backend)     │
                    │    :8080        │
                    └─────────────────┘
```

## Problem

Given a set of available pack sizes, calculates the minimum number of packs to fulfill an order:

1. Only whole packs can be sent
2. Send the least amount of items possible (>= order)
3. Use the fewest number of packs

## Running

### Development

```bash
# Terminal 1: Backend (API on :8080)
make dev-backend

# Terminal 2: Frontend with hot reload (:5173 proxies to :8080)
make dev-frontend
```

### Production (Docker)

```bash
# Build and run both services
make docker-up

# Run in background
make docker-up-detached

# View logs
make docker-logs

# Stop
make docker-down
```

Application runs on **http://localhost:80**

### Local Production Build

```bash
make build
./pack-calculator  # API only, frontend needs to be served separately
```

## API

Base URL: `http://localhost/api` (production) or `http://localhost:8080/api` (dev)

### Health check
```
GET /health
```

### Pack sizes
```
GET  /api/pack-sizes
PUT  /api/pack-sizes
POST /api/pack-sizes/add
POST /api/pack-sizes/remove
```

### Calculate packs
```
POST /api/calculate
{
  "amount": 12001
}

GET /api/calculate?amount=12001
```

Response:
```json
{
  "order_amount": 12001,
  "total_items": 12250,
  "total_packs": 4,
  "packs": {
    "5000": 2,
    "2000": 1,
    "250": 1
  }
}
```

## Tests

```bash
make test
make test-coverage
make bench
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build frontend + backend locally |
| `make run` | Run backend locally |
| `make dev-backend` | Run Go API in dev mode |
| `make dev-frontend` | Run Vite dev server |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make bench` | Run benchmarks |
| `make lint` | Lint Go + JS code |
| `make docker-up` | Build and run with docker-compose |
| `make docker-up-detached` | Run in background |
| `make docker-down` | Stop containers |
| `make docker-logs` | View container logs |
| `make docker-build` | Build Docker images |
| `make clean` | Remove build artifacts |

## Algorithm

Uses dynamic programming (similar to coin change problem) to find the optimal combination:

1. Build DP table where dp[i] = min packs to get exactly i items
2. Find smallest total >= requested amount
3. Reconstruct solution

Complexity: O(amount * num_pack_sizes)

## Examples

| Order | Packs | Total |
|-------|-------|-------|
| 1 | 1x250 | 250 |
| 250 | 1x250 | 250 |
| 251 | 1x500 | 500 |
| 501 | 1x500 + 1x250 | 750 |
| 12001 | 2x5000 + 1x2000 + 1x250 | 12250 |

### Edge case

Pack sizes [23, 31, 53], amount 500000:
Result: {23: 2, 31: 7, 53: 9429} = 500000 items

## Project Structure

```
pack-calculator/
├── cmd/server/main.go        # API entry point
├── internal/
│   ├── calculator/           # Pack calculation logic (DP algorithm)
│   ├── handler/              # Gin HTTP handlers
│   └── storage/              # In-memory storage (thread-safe)
├── web/                      # React + Vite + Tailwind
│   ├── src/
│   │   ├── App.tsx           # Main component
│   │   ├── api.ts            # API client
│   │   └── types.ts          # TypeScript interfaces
│   └── package.json
├── nginx/
│   └── nginx.conf            # Nginx configuration
├── Dockerfile.frontend       # Frontend: Node build + Nginx
├── Dockerfile.backend        # Backend: Go build + Alpine
├── docker-compose.yml        # Orchestrates frontend + backend
├── Makefile
└── go.mod
```

## Docker Images

| Service | Base Image | Purpose |
|---------|------------|---------|
| frontend | nginx:alpine | Serve static files, proxy API |
| backend | alpine:3.21 | Run Go binary |

Both images use multi-stage builds for minimal size.
