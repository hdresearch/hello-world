# Hello World

A three-service example for Vers Repos:

```text
Frontend (nginx) -> Backend (FastAPI) -> Database (PostgreSQL)
```

The example displays a persisted visit counter.

## Run locally

```bash
docker compose up --build -d
./scripts/verify-compose.sh
```

Open <http://localhost:8080>.

The database's HTTP gateway runs as the unprivileged `postgres` user and listens
on port `8080`, configurable through `PORT`. PostgreSQL itself listens on
`127.0.0.1:5432`. Compose publishes the gateway at <http://localhost:8082> and
connects the backend to `http://database:8080`.

## Images

```text
ghcr.io/hdresearch/hello-world-frontend:v0.1.0
ghcr.io/hdresearch/hello-world-backend:v0.1.0
ghcr.io/hdresearch/hello-world-database:v0.1.0
```
