# Benchmark Suite & Empirical Comparison

## Overview
This document records reproducible performance benchmarks comparing three distinct architectural strategies for high-throughput event ingestion within the Laravel Go Performance Kit.

---

## 1. Test Environment & Methodology
- **Load Generation Tool**: Grafana k6 (Containerized)
- **Database**: PostgreSQL 16 Alpine
- **Message Broker / Cache**: Redis 7 Alpine
- **Load Profile**: Ramping Virtual Users (0 -> 20 -> 50 VUs) over 40 seconds
- **Network**: Docker Bridge Network (`lgpk-network`)
- **Data Target**: `events` table (with composite index on `(tenant_id, occurred_at)` and JSONB payload)

---

## 2. Comparative Matrix

| Workload Strategy | Total Req (40s) | Throughput (RPS) | p50 Latency | p90 Latency | p95 Latency | Max Latency | Error Rate |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Laravel Sync** (Direct DB) | 73 | 1.05 req/s | 21.10s | 33.69s | 35.48s | 38.32s | 0.00% |
| **Laravel Async** (Redis Queue) | 61 | 0.87 req/s | 24.00s | 35.94s | 38.62s | 41.80s | 0.00% |
| **Go Ingestion Layer** | **68,959** | **1,718.60 req/s** | **4.31ms** | **5.81ms** | **6.21ms** | **129.80ms** | **0.00%** |

---

## 3. Key Architectural Takeaways

### The "Queue Fallacy"
Offloading synchronous database writes to a Redis queue in Laravel (`POST /api/events/async`) did not alleviate connection saturation under concurrent load. The root bottleneck was socket-level connection acceptance and HTTP runtime bootstrapping, which is serialized under single-process/worker-constrained environments.

### The Go Advantage
1. **Asynchronous Socket Multiplexing**: Go's native netpoll uses `epoll`/`kqueue` to accept thousands of concurrent sockets with negligible memory overhead.
2. **Goroutines**: Handlers execute concurrently across OS threads without thread starvation.
3. **Database Connection Pooling**: Direct integration with `pgxpool` allows concurrent requests to share persistent database connections, executing binary-protocol prepared inserts efficiently.

### Pragmatic Conclusion
Keep Laravel for SaaS application workflows (Authentication, RBAC, Billing, Admin, Business Logic). Selectively route high-concurrency ingestion endpoints (`/events`) through the Go performance layer.