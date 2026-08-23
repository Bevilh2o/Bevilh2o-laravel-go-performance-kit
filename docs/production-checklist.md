# Production Deployment Checklist

A production readiness guide for deploying the Laravel + Go Performance Kit in enterprise environments.

---

## 1. Network & Routing (Reverse Proxy)
- [ ] **Unified Ingress / API Gateway**: Configure Nginx, Envoy, or Cloudflare to route incoming traffic:
  - Route `/api/*` and web requests to Laravel (PHP-FPM / Octane).
  - Route `/events` (high-throughput ingestion) directly to the Go Ingestion Service (Port 8080).
- [ ] **TLS / SSL Termination**: Enforce TLS 1.3 at the reverse proxy layer.
- [ ] **HTTP Keep-Alive**: Enable keep-alive pooling between the reverse proxy and the Go service.

---

## 2. Go Service Hardening
- [ ] **Environment Variables**:
  - `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE` configured securely via secrets manager.
  - `HTTP_PORT` bound to internal network interface.
- [ ] **Database Connection Pool Tuning**:
  - Adjust `MaxConns` and `MinConns` in `repository/postgres.go` based on total database capacity and expected concurrency.
- [ ] **Profiling Endpoints**:
  - Restrict access to `/debug/pprof/` behind internal VPC / VPN authentication or disable it in public environments.
- [ ] **Process Management**:
  - Deploy Go container with restart policy `unless-stopped` or orchestrate via Kubernetes / ECS with liveness and readiness probes pointing to `GET /health`.

---

## 3. Database & Redis Optimization
- [ ] **PostgreSQL**:
  - Ensure connection pool max limits match `pgxpool.MaxConns` multiplied by Go service replica count.
  - Periodic `VACUUM ANALYZE` on the `events` table.
  - Monitor index usage on `(tenant_id, occurred_at)`.
- [ ] **Redis**:
  - Configure `maxmemory` and eviction policy (`volatile-lru` or `allkeys-lru`).
  - Enable Redis persistence (AOF or RDB snapshots) if queue durability is critical.

---

## 4. Security & Compliance
- [ ] **Service-to-Service Authentication**: Implement shared API secret or JWT verification between clients/Laravel and the Go ingestion service.
- [ ] **Rate Limiting**: Apply rate-limiting rules per tenant/IP at the reverse proxy or API gateway.
- [ ] **Input Sanitization**: Validate payload schema limits (max JSON body size enforced in Go).