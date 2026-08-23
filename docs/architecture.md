# Architectural Decision Records (ADR)

## ADR-001: Selective Go Offloading for High-Volume Event Ingestion

### Status
**Accepted**

### Context
High-concurrency event ingestion (`POST /events`) saturated Laravel's HTTP worker processes under 50 concurrent virtual users, resulting in request queuing latencies exceeding 20 seconds. Moving the database write to an asynchronous Redis queue in Laravel did not resolve the socket-level saturation.

### Decision
Introduce a dedicated, lightweight Go HTTP ingestion microservice that writes directly to PostgreSQL using an optimized connection pool (`pgxpool`).

### Consequences
- **Positive**:
  - Throughput increased from 1.05 req/s to 1,718.60 req/s (>1,600x gain).
  - Median latency dropped from 21.1s to 4.31ms.
  - Zero dropped connections or errors under load.
- **Trade-offs / Complexity**:
  - Introduced a second container runtime and deployment artifact.
  - Requires maintaining database schema consistency between Laravel migrations and Go entity definitions.

### Alternatives Considered
1. *Laravel Octane (Swoole / RoadRunner)*: Increases PHP throughput, but still retains framework bootstrap memory overhead and multi-tenancy state-leak risks compared to a pure Go microservice.
2. *Kafka / RabbitMQ intermediary*: Adds substantial operational complexity without necessity, given that PostgreSQL with connection pooling easily sustains ~1.7k+ RPS directly.