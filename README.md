# Relay

**Reliable event ingestion and delivery infrastructure built with Go and MongoDB Atlas.**

Relay is an open-source project exploring how to accept external events, persist them reliably, and deliver them to downstream services.

## Current Status

- [x] Go HTTP server
- [x] MongoDB Atlas connectivity
- [x] Liveness and readiness endpoints
- [ ] Event ingestion API
- [ ] Idempotent event persistence
- [ ] Concurrent delivery workers
- [ ] Automatic retries and dead-letter handling
- [ ] Event replay and delivery observability

## Technology

- Go
- MongoDB Atlas
- MongoDB Go Driver v2
- HTTP APIs

## Local Development

1. Install Go.
2. Copy `.env.example` to `.env`.
3. Configure `MONGODB_URI` locally.
4. Run `go mod tidy`.
5. Start the service with `go run ./cmd/relay`.

## Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/healthz` | Liveness check |
| GET | `/readyz` | Database readiness check |

## Security

Never commit `.env` or expose database credentials.

## License

A license will be added before the first public release.
