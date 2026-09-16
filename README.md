# Hello World on Vers

A small, real three-service application used by the Vers Repos launcher. It
proves that independently imaged Vers computers can communicate over their
HTTPS endpoints and persist state.

```text
Browser -> frontend/nginx -> backend/FastAPI -> database gateway -> PostgreSQL
```

The frontend displays a persisted visit counter. Reading or incrementing it
travels through all three services.

## Services

| Service | Runtime | Public endpoint | Purpose |
| --- | --- | --- | --- |
| Frontend | nginx Alpine | `:80` | Static UI and authenticated reverse proxy |
| Backend | Python 3.12 / FastAPI | `:80` | Application API and database client |
| Database | PostgreSQL 16 + Go gateway | `:80` | Narrow authenticated data API and persistence |

PostgreSQL listens only on localhost inside the database container. The data
gateway is the only network-facing process in that image.

## Run locally

Docker with Compose v2, `curl`, and `jq` are the only prerequisites.

```bash
docker compose up --build -d
./scripts/verify-compose.sh
```

Open <http://localhost:8080>. To stop the stack while retaining its database
volume, run `docker compose down`. Add `--volumes` to also delete local data.

## Vers release images

Chelsea launches these immutable release tags:

```text
ghcr.io/hdresearch/hello-world-frontend:v0.1.0
ghcr.io/hdresearch/hello-world-backend:v0.1.0
ghcr.io/hdresearch/hello-world-database:v0.1.0
```

Pushing a `v*` Git tag runs the release workflow. It verifies the complete
Compose path, builds Linux/AMD64 images, and publishes both the version tag and
the full Git commit SHA. Version tags are immutable and must never be moved.

The Vers launcher injects these runtime values; none are built into an image:

| Service | Required runtime environment |
| --- | --- |
| Frontend | `BACKEND_URL`, `BACKEND_TOKEN` |
| Backend | `FRONTEND_TOKEN`, `DATA_API_URL`, `DATA_API_TOKEN` |
| Database | `BACKEND_TOKEN`, `POSTGRES_PASSWORD` |
