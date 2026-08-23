# Laravel Go Performance Kit

A production-oriented SaaS starter kit combining **Laravel** for application and business logic with **Go** for performance-sensitive workloads.

The goal is not to rewrite Laravel applications in Go.

The goal is to demonstrate how a Laravel application can selectively introduce Go where **throughput, concurrency, latency, or resource efficiency** make it technically justified.

---

## Goal

Most Laravel applications do not need to be rewritten in another language when they start experiencing performance problems.

In many cases, the application itself is perfectly suitable for handling:

- authentication
- authorization
- business logic
- CRUD operations
- billing
- administration
- APIs
- integrations

The problems tend to appear in specific workloads.

This project explores an alternative approach:

> **Keep Laravel where Laravel is effective. Move only the workloads that benefit from Go.**

The project will implement equivalent workloads in Laravel and Go and benchmark them under increasing levels of concurrency and load.

The results will determine whether introducing Go provides a meaningful advantage for each workload.

---

## Architecture

```text
                         Client
                           │
                           ▼
                    ┌──────────────┐
                    │   Laravel    │
                    │              │
                    │ Application  │
                    │ Business     │
                    │ Logic        │
                    │ Auth / RBAC  │
                    │ API          │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │    Redis     │
                    │              │
                    │ Queue /      │
                    │ Cache /      │
                    │ Messaging    │
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │      Go      │
                    │              │
                    │ High-        │
                    │ throughput   │
                    │ workloads    │
                    │ Processing   │
                    │ Workers      │
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ PostgreSQL   │
                    └──────────────┘
```

The exact communication mechanism between Laravel and Go will depend on the workload being evaluated.

Possible approaches include:

- HTTP
- asynchronous queues
- Redis
- event-driven processing

The project will prefer the simplest mechanism that makes sense for each workload.

---

## Why Laravel + Go?

Laravel provides a productive environment for building application-level functionality.

Go provides a different set of characteristics that can be useful for specific backend workloads:

- lightweight concurrency
- efficient memory usage
- predictable execution
- high-throughput network services
- inexpensive worker processes
- straightforward deployment as a single binary

However, these characteristics do not automatically make Go a better choice.

This project therefore does **not** assume that Go will be faster.

Instead, every performance-sensitive component should answer three questions:

1. Is there an actual bottleneck?
2. Does moving the workload to Go address that bottleneck?
3. Is the improvement significant enough to justify the additional architectural complexity?

---

## Workloads

The project will focus on workloads where performance characteristics can be measured objectively.

Potential workloads include:

### Event ingestion

High-volume API requests representing application events.

```text
POST /events
```

Example:

```json
{
  "tenant": "tenant_123",
  "event": "page_view",
  "timestamp": 1724312345
}
```

The system will evaluate how Laravel and Go behave when processing increasing numbers of concurrent events.

### Data processing

Processing, transforming, validating, or aggregating large numbers of events.

### Background workers

Comparing Laravel queue workers with Go workers for CPU- or throughput-sensitive tasks.

### Concurrent processing

Evaluating workloads that can be split into independent units of work and processed concurrently.

The final workload selection will be based on actual profiling and benchmarking rather than assumptions.

---

## Benchmarking

Performance comparisons will be performed using controlled workloads.

The benchmark suite will measure:

| Metric | Description |
|---|---|
| Throughput | Requests or jobs processed per second |
| p50 latency | Median response/processing latency |
| p95 latency | 95th percentile latency |
| p99 latency | 99th percentile latency |
| CPU usage | CPU consumed under load |
| Memory usage | Process memory consumption |
| Concurrency | Behaviour under increasing concurrent workloads |
| Error rate | Failed requests/jobs under load |
| Database load | Impact on the database where relevant |

The objective is not to produce a single "Go vs Laravel" number.

Different workloads may produce different results.

---

## Benchmark Methodology

The same workload should be implemented as similarly as possible in both environments.

For example:

```text
Laravel implementation
        │
        ▼
   Benchmark
        │
        ├── throughput
        ├── latency
        ├── CPU
        └── memory


Go implementation
        │
        ▼
   Benchmark
        │
        ├── throughput
        ├── latency
        ├── CPU
        └── memory
```

Tests should be repeated under multiple load levels rather than relying on a single benchmark.

Example:

```text
10 concurrent requests
100 concurrent requests
500 concurrent requests
1,000 concurrent requests
5,000 concurrent requests
...
```

The actual levels will depend on the workload and available hardware.

---

## Performance First, Complexity Second

A faster implementation is not automatically a better implementation.

Moving functionality from Laravel to Go introduces additional complexity:

- another runtime
- another service
- another deployment artifact
- additional communication
- additional monitoring
- additional failure modes

Therefore the project will consider **performance per unit of architectural complexity**, rather than performance alone.

The preferred architecture is:

> **The simplest architecture that provides the required performance characteristics.**

---

## Observability

Performance measurements should be reproducible and observable.

The project will eventually include:

- structured logging
- health checks
- application metrics
- request/job timing
- resource usage measurements
- benchmark results
- profiling where appropriate

Where useful, Go profiling will use standard tooling such as `pprof`.

---

## Project Structure

The initial structure is intentionally minimal.

```text
laravel-go-performance-kit/
│
├── laravel/
│   └── ...
│
├── go/
│   └── ...
│
├── benchmarks/
│   └── ...
│
├── docs/
│   ├── architecture.md
│   ├── benchmarking.md
│   └── performance.md
│
├── docker/
│   └── ...
│
└── README.md
```

The structure may evolve as the project develops.

---

## Design Principles

### 1. Do not rewrite everything

Laravel remains the primary application.

Go is introduced selectively.

### 2. Measure before optimizing

No component should be moved to Go solely because Go is expected to be faster.

### 3. Keep boundaries explicit

Laravel and Go should communicate through clearly defined interfaces.

### 4. Prefer asynchronous processing where appropriate

Not every workload needs to block the HTTP request.

### 5. Keep the Go components independently deployable

Go components should be capable of being built, tested, benchmarked, and deployed independently from Laravel.

### 6. Keep the system understandable

Performance improvements should not come at the cost of unnecessary architectural complexity.

---

## Roadmap

### Phase 1 — Foundation
- [x] Create Laravel application
- [x] Define initial domain
- [x] Add PostgreSQL
- [x] Add Redis
- [x] Add Docker development environment
- [x] Establish baseline Laravel implementation

### Phase 2 — Baseline Performance
- [x] Define representative workloads
- [x] Build load tests
- [x] Establish Laravel baseline
- [x] Measure CPU and memory usage
- [x] Record latency distributions

### Phase 3 — Go Performance Layer
- [x] Identify the first performance-sensitive workload
- [x] Implement the workload in Go
- [x] Define Laravel ↔ Go communication
- [x] Add concurrent processing where appropriate
- [x] Benchmark the Go implementation

### Phase 4 — Comparison
- [x] Run equivalent workloads
- [x] Compare throughput
- [x] Compare p50/p95/p99 latency
- [x] Compare CPU usage
- [x] Compare memory usage
- [x] Document trade-offs

### Phase 5 — Production Hardening
- [x] Health checks
- [x] Structured logging (`slog`)
- [x] Metrics & pprof profiling
- [x] Failure handling & panic recovery
- [x] Graceful shutdown
- [x] Containerized deployment
- [x] Load testing under realistic conditions

### Phase 6 — Starter Kit
- [x] Clean installation process
- [x] Configuration documentation
- [x] Example workload
- [x] Example Go component
- [x] Benchmark suite (k6)
- [x] Deployment documentation
- [x] Production checklist

---

## Current Status

**Production-Ready Starter Kit (v1.0.0)**

All foundational components, benchmarks, and production-hardening layers have been implemented and validated empirically. The benchmark suite demonstrates a **>1,600x throughput increase** and **~4,900x lower median latency** by offloading high-volume event ingestion from Laravel to the Go performance layer.

---

## Philosophy

This project is built around a simple idea:

> **Don't use Go because it is fast. Use Go when being fast matters.**