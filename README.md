<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&height=240&text=NeverNet%20Edge%20API&fontAlign=50&fontAlignY=40&color=0:081226,35:123A73,70:2F6BFF,100:8BE9FF&fontColor=EAF8FF&desc=Blue%20Unit%20%7C%20Go%20Gateway%20%7C%20gRPC%20Bridge&descAlignY=62&animation=fadeIn" width="100%" />

# ❄️ NeverNet Edge API

<p>
  <img src="https://img.shields.io/badge/Go-1.23+-00BFFF?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/gRPC-Rust%20backend-2F6BFF?style=for-the-badge&logo=grpc&logoColor=white" />
  <img src="https://img.shields.io/badge/Swagger-OpenAPI-38BDF8?style=for-the-badge&logo=swagger&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-Linux-0A84FF?style=for-the-badge&logo=docker&logoColor=white" />
  <img src="https://img.shields.io/badge/Router-chi-5BC0FF?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Status-ONLINE-8BE9FF?style=for-the-badge" />
</p>

<img src="./assets/Ayanami.jpg" width="320" alt="Ayanami Rei" />

> **Cold blue gateway between frontend and Rust microservices.**
> Minimal. Fast. Dockerized. Swagger-ready.

</div>

---

## ✦ About

**NeverNet Edge API** is a Go-based HTTP gateway / BFF that communicates with Rust backend services over **gRPC**.

It is designed as a replacement for `backend/apps/edge-api` and exposes a stable REST interface for the frontend.

---

## ✦ Blue System State

```text
SYSTEM      :: ONLINE
INTERFACE   :: BLUE
TRANSPORT   :: HTTP <-> gRPC
STATUS      :: STABLE
MODE        :: EDGE API
SYNC        :: READY
```

---

## ✦ Stack

* **Go**
* **chi** — HTTP router
* **gRPC Go** — Rust service clients
* **Swagger / OpenAPI**
* **Docker**
* **Env-based config**

---

## ✦ Features

* REST API gateway
* gRPC bridge to Rust backend
* health endpoints
* Swagger / OpenAPI docs
* CSRF + cookie auth flow
* Docker-ready deploy
* simple hackathon-friendly structure

---

## ✦ Project Structure

```text
.
├── api/
├── assets/
│   └── Ayanami.jpg
├── cmd/
│   └── edge-api/
├── gen/
├── internal/
│   ├── app/
│   ├── config/
│   ├── docs/
│   ├── grpcclient/
│   ├── handlers/
│   ├── middleware/
│   └── response/
├── .dockerignore
├── .env.example
├── .gitignore
├── buf.gen.yaml
├── buf.yaml
├── docker-compose.local.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## ✦ Quick Start

### Clone

```bash
git clone https://github.com/Denbay0/TestApi.git
cd TestApi
```

### Create `.env`

#### Linux / macOS

```bash
cp .env.example .env
```

#### Windows PowerShell

```powershell
Copy-Item .env.example .env -Force
```

### Prepare dependencies

```bash
go mod tidy
go mod vendor
go build ./cmd/edge-api
```

### Run locally

```bash
go run ./cmd/edge-api
```

Service will be available at:

* `http://localhost:8080`
* `http://localhost:9100`

When you open `http://localhost:8080/` in a browser, the service now returns a small landing page with links to docs, OpenAPI artifacts, and health checks instead of a 404.

---

## ✦ Swagger / OpenAPI

After startup:

* **Swagger UI** → `http://localhost:8080/docs`
* **OpenAPI JSON** → `http://localhost:8080/openapi.json`
* **OpenAPI YAML** → `http://localhost:8080/openapi.yaml`

The `/docs` page is fully self-contained, uses only assets served by the edge-api process itself, and fetches the local `/openapi.json`, so it remains available without internet access. The page groups endpoints by tags, shows request/response schemas and examples, documents cookie + CSRF security, and includes a built-in Try it out flow for interactive requests.

---


Quick browser checks:

* `http://localhost:8080/`
* `http://localhost:8080/docs`
* `http://localhost:8080/openapi.json`
* `http://localhost:8080/openapi.yaml`
* `http://localhost:8080/health`
* `http://localhost:8080/healthz`

## ✦ Health Checks

```bash
curl http://localhost:8080/health
curl http://localhost:8080/healthz
curl http://localhost:9100/health
```

---

## ✦ Docker

### Pull base images

```bash
docker pull golang:1.23-alpine
docker pull alpine:3.20
```

### Build image

```bash
docker build -t edge-api:test .
```

### Run container

```bash
docker run --rm -p 8080:8080 -p 9100:9100 --env-file .env edge-api:test
```

After the container starts, open `/docs` in your browser or use the smoke-test commands below.

---

## ✦ Quick Smoke Test

```bash
curl -i http://localhost:8080/
curl -i http://localhost:8080/favicon.ico
curl -i http://localhost:8080/health
curl -i http://localhost:8080/healthz
curl -i http://localhost:8080/openapi.json
curl -i http://localhost:8080/openapi.yaml
```

## ✦ Docker Compose

```bash
docker compose -f docker-compose.local.yml up --build
```

If an external network is required:

```bash
docker network create rust-backend
```

---

## ✦ Environment Variables

```env
PORT=8080
METRICS_PORT=9100
REDIS_URL=redis://redis:6379

IDENTITY_SERVICE_URL=http://identity-svc:50051
EVENT_COMMAND_SERVICE_URL=http://event-command-svc:50052
EVENT_QUERY_SERVICE_URL=http://event-query-svc:50053
REPORT_SERVICE_URL=http://report-svc:50054

FRONTEND_ORIGINS=http://localhost:3000,http://localhost:5173
AUTH_COOKIE_SECURE=false
```

---

## ✦ Main Endpoints

### Auth

* `GET /api/auth/csrf`
* `POST /api/auth/register`
* `POST /api/auth/login`
* `POST /api/auth/logout`
* `GET /api/auth/me`

### Data

* `GET /api/categories`
* `GET /api/events`
* `POST /api/events`
* `GET /api/calendar`
* `GET /api/dashboard`
* `GET /api/reports/summary`
* `GET /api/reports/by-category`
* `GET /api/settings`
* `PUT /api/settings`
* `GET /api/exports`

---

## ✦ Git Policy

Tracked:

* `go.mod`
* `go.sum`
* source code
* docs
* Docker files

Ignored:

* `.env`
* `vendor/`
* binaries like `*.exe`
* local artifacts

---

## ✦ Development Notes

This project is optimized for:

* fast local iteration
* Linux Docker runtime
* simple deployment flow
* easy integration with existing Rust backend

---

## ✦ Roadmap

* [ ] connect real protobuf/gRPC contracts
* [ ] complete compatibility with Rust edge-api
* [ ] add integration tests
* [ ] harden auth/session flow
* [ ] production deployment profile

---

## ✦ Blue Interface Quote

> *Silence in the transport. Precision in the gateway. Stability in the build.*

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&section=footer&height=160&color=0:081226,35:123A73,70:2F6BFF,100:8BE9FF" width="100%" />

</div>
