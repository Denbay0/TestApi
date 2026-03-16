# Edge API (Go)

Production-oriented Go API gateway/BFF that is intended to become a drop-in HTTP replacement for `backend/apps/edge-api` from the existing Rust monorepo. It talks to existing backend services over gRPC and exposes a stable HTTP/JSON contract.

> Assumption: exact Rust `.proto` contracts are not yet available in this repository, so this version includes scaffold proto + TODO wiring points for quick hackathon iteration.

## Project tree

```text
.
├── .env.example
├── Dockerfile
├── Makefile
├── README.md
├── api
│   └── proto
│       └── backend.proto
├── buf.gen.yaml
├── buf.yaml
├── cmd
│   └── edge-api
│       └── main.go
├── docker-compose.local.yml
├── gen
│   └── go
│       └── README.md
├── go.mod
└── internal
    ├── app
    │   └── app.go
    ├── auth
    │   └── cookies.go
    ├── config
    │   └── config.go
    ├── docs
    │   └── docs.go
    ├── grpcclient
    │   ├── client.go
    │   ├── event_command.go
    │   ├── event_query.go
    │   ├── identity.go
    │   ├── report.go
    │   └── target.go
    ├── handlers
    │   └── handlers.go
    ├── middleware
    │   ├── authcookie.go
    │   ├── context.go
    │   ├── csrf.go
    │   ├── logging.go
    │   └── requestid.go
    └── response
        └── response.go
```

## What this service includes

- HTTP server on `:8080` (configurable via env).
- Chi router + middleware stack:
  - request id
  - recoverer
  - real ip
  - structured logging (`slog`)
  - CORS
  - CSRF validation for mutating methods
  - auth cookie parsing
- Health endpoints:
  - `GET /health`
  - `GET /healthz`
- API contract envelope:
  - success: `{ "data": ... }`
  - error: `{ "error": { "code", "message", "request_id" } }`
- Swagger/OpenAPI endpoints from the service itself:
  - `GET /docs`
  - `GET /openapi.json`
  - `GET /openapi.yaml`
- gRPC connection layer with target normalization (`host:port`, `http://host:port`, `https://host:port`).
- Docker multi-stage build with static binary (`CGO_ENABLED=0`) and minimal Linux runtime image.

## Requirements

- Go 1.23+
- Docker / Docker Compose
- Buf + protoc plugins (for protobuf generation when real proto contracts are ready)

## Install dependencies

```bash
go mod tidy
```

If you need local protobuf tooling:

```bash
go install github.com/bufbuild/buf/cmd/buf@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Protobuf generation

```bash
make proto
```

This runs `buf generate` using `buf.gen.yaml`, outputs to `gen/go`.

## Local run

```bash
cp .env.example .env
make run
```

Service endpoints:
- API: `http://localhost:8080`
- Metrics health: `http://localhost:9100/health`
- Swagger UI: `http://localhost:8080/docs`

## Docker run

Build image:

```bash
make docker
```

Run compose integration:

```bash
cp .env.example .env
docker compose -f docker-compose.local.yml up --build
```

## Environment variables

- `PORT=8080`
- `METRICS_PORT=9100`
- `REDIS_URL=redis://redis:6379`
- `IDENTITY_SERVICE_URL=http://identity-svc:50051`
- `EVENT_COMMAND_SERVICE_URL=http://event-command-svc:50052`
- `EVENT_QUERY_SERVICE_URL=http://event-query-svc:50053`
- `REPORT_SERVICE_URL=http://report-svc:50054`
- `FRONTEND_ORIGINS=http://localhost:3000,http://localhost:5173`
- `AUTH_COOKIE_SECURE=false`
- `OPENAPI_SERVER_URL=http://localhost:8080`

## HTTP endpoints scaffolded

- `GET /api/auth/csrf`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/me`
- `GET /api/categories`
- `GET /api/events`
- `POST /api/events`
- `GET /api/calendar`
- `GET /api/dashboard`
- `GET /api/reports/summary`
- `GET /api/reports/by-category`
- `GET /api/settings`
- `PUT /api/settings`
- `GET /api/exports`

## CSRF/Auth behavior

- `GET /api/auth/csrf`: creates CSRF token cookie and returns token in JSON envelope.
- Mutating methods (`POST`, `PUT`, `PATCH`, `DELETE`) require `X-CSRF-Token` equal to CSRF cookie value.
- `register`/`login`: set auth cookie.
- `logout`: clears auth and CSRF cookies.
- Cookie secure flag comes from `AUTH_COOKIE_SECURE`.

## How to connect to existing Rust backend

1. Ensure Rust services are running and reachable on Docker network.
2. Configure `*_SERVICE_URL` env vars to those service aliases/ports.
3. Keep this Go service on same network (`rust-backend` in compose example).
4. Replace scaffold `.proto` in `api/proto` with real contracts from Rust monorepo.
5. Run `make proto` and wire grpcclient TODO methods to generated clients.

## Linux server deployment switch plan

1. Build Linux image in CI:
   - `docker build -t registry/edge-api:<tag> .`
2. Push image to registry.
3. On Linux server update compose/k8s manifest image tag.
4. Set production env (`AUTH_COOKIE_SECURE=true`, prod origins, internal service URLs).
5. Roll out and verify:
   - `/health`
   - `/docs`
   - key API smoke checks.

## TODO for real integration

- Replace scaffold proto with real Rust backend `.proto` files.
- Generate `gen/go` from real contracts.
- Implement real grpc calls in `internal/grpcclient/*` (currently placeholder responses).
- Map Rust domain errors into custom error envelope code/message catalog.
- Add integration tests against running Rust services.
- Add metrics/tracing (Prometheus + OpenTelemetry).
- Add cookie/session hardening strategy and secret rotation.
- Add CI pipeline for lint/test/build/docker/proto-breaking checks.
