# Benchmark Suite & Baseline Results

## Environment
- **Load Tool**: Grafana k6 (Containerized)
- **Database**: PostgreSQL 16 (Alpine)
- **Broker / Cache**: Redis 7 (Alpine)
- **Workload**: High-Volume Event Ingestion (`POST /api/events` vs `POST /api/events/async`)
- **Concurrency**: Ramping 0 -> 20 -> 50 Virtual Users (VUs) over 40s

---

## 1. Baseline Measurements (Laravel)

### Summary Table

| Workload Strategy | Throughput | p50 Latency | p90 Latency | p95 Latency | Max Latency | Error Rate |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Laravel Sync (PostgreSQL)** | 1.05 req/s | 21.10s | 33.69s | 35.48s | 38.32s | 0.00% |
| **Laravel Async (Redis Queue)** | 0.87 req/s | 24.00s | 35.94s | 38.62s | 41.80s | 0.00% |

### Architectural Root-Cause Analysis
- **Socket Serialization**: When concurrency increases beyond available worker threads, client connections are queued at the OS socket level.
- **Framework Lifecycle Overhead**: Offloading DB persistence to a Redis queue (`POST /api/events/async`) does not alleviate the HTTP entry saturation because every request still incurs the full framework bootstrap and dependency resolution cost per connection.
- **Conclusion**: The bottleneck is the HTTP connection handling layer under concurrent socket load, making this workload an ideal candidate for a lightweight Go ingestion layer.