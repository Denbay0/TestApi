FROM golang:1.23-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY gen ./gen

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -trimpath -ldflags='-s -w' -o /out/edge-api ./cmd/edge-api

FROM alpine:3.20 AS runtime
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/edge-api /app/edge-api

EXPOSE 8080 9100
ENTRYPOINT ["/app/edge-api"]
