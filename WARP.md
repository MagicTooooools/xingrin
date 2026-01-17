# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Common commands

### Repo scripts (Docker-based deployment)
- `./install.sh` (prod) / `./install.sh --dev` (build local dev images) / `./install.sh --no-frontend`
- `./start.sh` / `./start.sh --dev`
- `./stop.sh`
- `./restart.sh`
- `./update.sh`
- `./uninstall.sh`

### Docker Compose (full stack with Django backend)
- Dev: `docker compose -f docker/docker-compose.dev.yml up -d`
- Dev + local Postgres: `docker compose -f docker/docker-compose.dev.yml --profile local-db up -d`
- Down: `docker compose -f docker/docker-compose.dev.yml down`

### Frontend (Next.js in `frontend/`)
- `cd frontend`
- `pnpm install`
- `pnpm dev` (or `pnpm dev:mock`, `pnpm dev:noauth`)
- `pnpm build`
- `pnpm start`
- `pnpm lint`

### Go server rewrite (Gin/GORM in `server/`)
- `cd server`
- `make run` / `make build` / `make test` / `make lint`
- Single test: `go test ./internal/... -run TestName`
- Dev deps only (Postgres/Redis): `docker compose -f docker-compose.dev.yml up -d`

### Go worker (scan executor in `worker/`)
- `cd worker`
- `make run` / `make build` / `make test`
- Single test: `go test ./internal/... -run TestName`

### Django backend (production server in `backend/`)
- `cd backend`
- `pytest`
- Single test: `pytest apps/<app>/... -k "TestName or test_name"`

### Seed data generator (API-based, `tools/seed-api/`)
- `cd tools/seed-api`
- `pip install -r requirements.txt`
- `python seed_generator.py` (see `tools/seed-api/README.md` for options)
- Tests: `pytest` (integration requires a running backend)

## Architecture overview (big picture)
- **Monorepo services**: 
  - **Django backend** in `backend/` (current production server, runs via `docker/server/start.sh` with migrations + `uvicorn`).
  - **Go backend rewrite** in `server/` (Gin/GORM; incomplete per `server/README.md`).
  - **Next.js frontend** in `frontend/` (App Router; containerized via `docker/frontend/Dockerfile`).
  - **Go worker** in `worker/` plus **agent** (heartbeat/monitor) in `docker/agent/Dockerfile`.
- **Deployment topology**: `docker/docker-compose.yml` (prod) and `docker/docker-compose.dev.yml` (dev) orchestrate Postgres, Redis, Django server, agent, frontend, and nginx. Nginx terminates HTTPS on `8083` and proxies to backend `8888`.
- **Versioning**: `VERSION` is the single source of release version. `IMAGE_TAG` in `docker/.env` pins all images to the same tag; `./update.sh` refreshes it (see `docs/version-management.md`).
- **Scan pipeline**: stages and toolchain are documented in `docs/scan-flow-architecture.md`. Stage ordering is defined in `backend/apps/scan/configs/command_templates.py`.
- **Templates & wordlists**:
  - Nuclei templates: server-side sync and worker-side checkout are documented in `docs/nuclei-template-architecture.md`.
  - Wordlists: upload on server, hash-based cache + download on worker (see `docs/wordlist-architecture.md`).
- **Backend domain layout**: Django apps under `backend/apps/` (e.g., `scan`, `engine`, `asset`, `targets`, `common`). Worker/agent deployment helpers live in `backend/scripts/worker-deploy/`.
- **Go server layout**: `server/internal/` is layered (`handler` → `service` → `repository` → `model`, plus `dto`, `middleware`, `config`, `database`).
- **Go worker layout**: `worker/internal/` splits workflow orchestration (`workflow`, `activity`) from runtime/server glue (`server`, `config`, `pkg`).
- **Config files**: `.env` templates live in `docker/.env.example` (Django stack), `server/.env.example` (Go server), and `worker/.env.example` (Go worker).
