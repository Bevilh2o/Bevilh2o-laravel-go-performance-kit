# Laravel + Go Performance Kit --- Project Context

## Purpose

Persistent technical context for the project. Use this as the source of
truth when extending the codebase in future sessions.

The project is a production-oriented Laravel SaaS starter kit with a Go
performance layer.

> Don't use Go because it is fast. Use Go when being fast matters.

The goal is to demonstrate how a Laravel application can selectively
introduce Go for workloads where throughput, concurrency, latency, or
resource efficiency provide a measurable benefit.

We are NOT trying to rewrite Laravel in Go.

## 1. Product Concept

Laravel remains responsible for normal application concerns:

-   authentication
-   authorization
-   business logic
-   CRUD
-   administration
-   billing
-   APIs
-   integrations
-   domain workflows

Go is introduced selectively for performance-sensitive workloads such
as:

-   high-throughput event ingestion
-   concurrent processing
-   background workers
-   CPU-intensive processing
-   data transformation
-   aggregation
-   streaming

The project should remain understandable and practical. A performance
improvement is not worthwhile if it introduces disproportionate
architectural complexity.

## 2. Main Architectural Principle

> Laravel owns the application and business domain. Go owns workloads
> that have a demonstrated performance reason to exist outside Laravel.

Before moving functionality to Go, ask:

1.  Is there an actual bottleneck?
2.  What is causing the bottleneck?
3.  Does Go address that specific bottleneck?
4.  How large is the measurable improvement?
5.  Is the improvement worth the additional operational complexity?

## 3. Planned Architecture

    Client
       |
       v
    Laravel
       |
       +---- Business logic
       +---- Authentication
       +---- Authorization
       +---- API
       +---- SaaS/domain functionality
       |
       v
    Redis
       |
       v
    Go
       |
       +---- High-throughput workloads
       +---- Concurrent processing
       +---- Workers
       +---- Data processing
       |
       v
    PostgreSQL

The Laravel ↔ Go communication mechanism must be chosen according to the
workload. Possible mechanisms include HTTP, Redis, asynchronous queues,
and event-driven communication.

Do not introduce gRPC, Kafka, or another technology merely for
architectural complexity.

## 4. Technology Stack

### Core

-   Laravel
-   Go
-   PostgreSQL
-   Redis
-   Docker

### Laravel

Keep Laravel relatively close to the standard architecture.

Potential technologies:

-   Blade
-   Alpine.js
-   Livewire when useful
-   Laravel Queue
-   Laravel Horizon
-   Laravel Sanctum or another appropriate authentication mechanism

Avoid unnecessary frontend complexity.

### Go

Prefer idiomatic, standard-library-first design where practical.

Potential areas:

-   goroutines
-   channels
-   worker pools
-   context.Context
-   net/http
-   pprof
-   benchmarks

Do not introduce frameworks or dependencies without a reason.

## 5. Frontend Strategy

The frontend is not the primary purpose of this project.

For the initial version:

-   prefer Blade + Alpine.js or another lightweight Laravel-native
    approach
-   Livewire is acceptable when it improves development speed
-   avoid making React/Next/etc. central to the project
-   avoid frontend complexity merely for visual appeal

## 6. Performance Philosophy

Performance work must be evidence-driven.

Never claim that Go is automatically faster or that Laravel cannot
handle high traffic.

Instead:

1.  implement a representative workload in Laravel
2.  establish a baseline
3.  measure it
4.  identify the bottleneck
5.  implement the appropriate Go version
6.  measure again
7.  compare results
8.  document trade-offs

The project should demonstrate optimization and architectural restraint.

## 7. Benchmarking

Primary metrics:

-   throughput
-   p50 latency
-   p95 latency
-   p99 latency
-   CPU usage
-   memory usage
-   concurrency behaviour
-   error rate
-   database load
-   processing time

Benchmarks should use controlled and repeatable workloads.

Example load levels:

-   10 concurrent requests
-   100
-   500
-   1,000
-   higher levels when appropriate

Do not compare numbers from materially different machines or
configurations without documenting the difference.

## 8. Candidate Workload

The first likely workload is high-volume event ingestion.

Example:

    POST /events

Example payload:

    {
      "tenant": "tenant_123",
      "event": "page_view",
      "timestamp": 1724312345
    }

Potential pipeline:

    Client
      |
      v
    Laravel API
      |
      v
    Queue / Redis
      |
      v
    Go worker
      |
      +---- validate
      +---- transform
      +---- aggregate
      +---- batch
      |
      v
    PostgreSQL

This is a candidate, not a permanent decision. The actual workload
should be selected after profiling and architectural evaluation.

## 9. Repository Structure

Initial intended structure:

    laravel-go-performance-kit/
    |
    +-- laravel/
    +-- go/
    +-- benchmarks/
    +-- docs/
    |   +-- architecture.md
    |   +-- benchmarking.md
    |   +-- performance.md
    +-- docker/
    +-- README.md
    +-- PROJECT_CONTEXT.md

This structure may evolve. If it changes, update this document.

## 10. Code Organization

### Laravel

Keep responsibilities explicit.

Prefer:

    Controller
        |
        v
    Application / Domain Service
        |
        v
    Repository / Model / Infrastructure

Do not put large amounts of business logic directly inside controllers.

Use dependency injection.

### Go

Prefer small packages with clear responsibilities.

Avoid large `utils` packages.

Prefer:

-   explicit dependencies
-   small interfaces where they provide real value
-   context propagation
-   clear error handling
-   simple concurrency models

Do not introduce abstractions before they are needed.

## 11. Error Handling

Prefer wrapping errors with context:

    fmt.Errorf("process event: %w", err)

Use sentinel or typed errors when callers need to distinguish expected
conditions.

Do not use error strings as an API for programmatic branching.

HTTP/API error formatting belongs at the transport boundary.

## 12. Context and Cancellation

Long-running and I/O-bound operations should accept `context.Context`.

Propagate context through:

-   database operations
-   HTTP calls
-   queue operations
-   external services
-   Go worker operations where cancellation matters

Graceful shutdown is required for long-running Go services.

## 13. Concurrency

Concurrency should be used when the workload benefits from overlapping
independent work.

Do not equate:

    goroutines = faster

Goroutines provide concurrency; parallel speedup depends on CPU cores,
workload, synchronization, I/O, contention, scheduling, and memory
behaviour.

When concurrency is introduced, document why it exists and benchmark the
result.

Avoid shared mutable state where possible.

## 14. Database Principles

PostgreSQL is the initial relational database.

When investigating database performance:

-   inspect query plans
-   inspect indexes
-   measure query latency
-   consider connection pool behaviour
-   measure database load
-   distinguish application bottlenecks from database bottlenecks

Do not move work into Go merely because a query is slow. Fix the actual
bottleneck first.

## 15. Redis

Redis may be used for:

-   queues
-   caching
-   transient state
-   event coordination
-   rate limiting where appropriate

Redis should not become a generic dumping ground for state that belongs
in PostgreSQL.

## 16. Docker

The development environment should eventually provide services for:

-   Laravel
-   Go
-   PostgreSQL
-   Redis

Development should be reproducible.

## 17. Observability

The project should eventually provide:

-   health checks
-   structured logging
-   metrics
-   request timing
-   job timing
-   error tracking
-   CPU/memory measurements
-   Go profiling where appropriate

For Go, use standard profiling tools such as `pprof` where useful.

## 18. Testing

Both Laravel and Go components should have automated tests.

Tests should cover:

-   business rules
-   API behaviour
-   error conditions
-   concurrency-sensitive logic
-   integration boundaries
-   Go worker behaviour

Performance benchmarks are not substitutes for functional tests.

## 19. API and Service Boundaries

Laravel ↔ Go communication should use explicit contracts.

Contracts should define:

-   request shape
-   response shape
-   errors
-   timeouts
-   retry behaviour
-   idempotency where relevant
-   authentication/authorization
-   versioning strategy when required

The Go service should not directly depend on Laravel internals.

## 20. Reliability

For asynchronous or retried operations, consider:

-   idempotency keys
-   duplicate messages
-   retry policies
-   dead-letter handling
-   partial failure
-   graceful degradation

Do not assume queues deliver exactly once.

## 21. Security

Security must not be sacrificed for performance.

Consider:

-   authentication
-   authorization
-   input validation
-   rate limiting
-   secret management
-   TLS
-   least privilege
-   database credentials
-   service-to-service authentication

Benchmarks should not silently disable security mechanisms.

## 22. Commercial Product Direction

Long-term, the project may become a developer starter kit.

Possible positioning:

> A production-ready Laravel SaaS foundation with an optional Go
> performance layer for workloads that need higher throughput or
> concurrency.

The product should not promise that Go makes every Laravel application
faster.

Potential future contents:

-   SaaS foundation
-   authentication
-   multi-tenancy
-   RBAC
-   billing
-   queues
-   Redis
-   Go worker/service
-   benchmark suite
-   observability
-   Docker setup
-   deployment documentation
-   production checklist

Commercial features should be added only after the technical foundation
is validated.

## 23. Portfolio Positioning

Target professional positioning:

**Backend Engineer --- Performance --- Go**

The project should demonstrate:

-   backend architecture
-   performance analysis
-   profiling
-   benchmarking
-   concurrency
-   Go
-   Laravel
-   databases
-   Redis
-   service architecture
-   practical engineering judgement

The important story is:

> I identified a workload, measured the baseline, evaluated the
> bottleneck, introduced Go where it made sense, and measured the
> result.

## 24. Documentation Standards

For significant architectural decisions, record:

-   problem
-   alternatives considered
-   chosen approach
-   reason
-   trade-offs
-   benchmark evidence where applicable

Do not rewrite this document for minor implementation details.

## 25. Current Status

**Early development --- architecture and benchmarking methodology.**

Current state:

-   repository concept defined
-   README defined
-   architecture direction defined
-   no performance claims yet
-   first workload not permanently selected
-   implementation has not started

Do not assume future components already exist.

## 26. Rules for Future Code Generation

When generating code for this project:

1.  Preserve the architecture described here.
2.  Prefer simple solutions.
3.  Do not add dependencies without a clear reason.
4.  Do not introduce Go where there is no demonstrated need.
5.  Do not optimize without identifying what is being optimized.
6.  Keep Laravel responsible for normal business/application concerns.
7.  Keep Go components focused on clearly defined workloads.
8.  Keep service boundaries explicit.
9.  Propagate context and handle cancellation correctly.
10. Write tests alongside significant functionality.
11. Include benchmarks for performance-sensitive components.
12. Document important architectural decisions.
13. Do not invent components that have not yet been implemented.
14. If a proposed change conflicts with this document, flag the conflict
    before changing the architecture.

## 27. Decision Log

### D001 --- Laravel remains the primary application

**Decision:** Laravel is the main SaaS/application layer.

**Reason:** It provides the productivity and ecosystem needed for normal
application development.

**Status:** Accepted.

### D002 --- Go is introduced selectively

**Decision:** Go is used only for workloads where performance
characteristics justify it.

**Reason:** Avoid unnecessary complexity and avoid treating Go as a
universal replacement for Laravel.

**Status:** Accepted.

### D003 --- Benchmark before making performance claims

**Decision:** Performance claims must be supported by measurements.

**Reason:** The project is explicitly about performance engineering and
should demonstrate evidence-based decisions.

**Status:** Accepted.

### D004 --- Keep the first frontend simple

**Decision:** Do not make a heavy frontend framework central to V1.

**Reason:** The project is primarily a backend/performance
demonstration.

**Status:** Accepted.

## 28. Future Decision Log Template

    ## D00X — Title

    **Decision:**

    **Reason:**

    **Alternatives considered:**

    **Trade-offs:**

    **Benchmark evidence:**

    **Status:** Proposed / Accepted / Rejected / Superseded

## 29. Source of Truth

This file describes the intended architecture and engineering
conventions.

If implementation and this document diverge, do not silently assume
which one is correct.

First determine whether:

1.  the implementation is wrong,
2.  the documentation is outdated, or
3.  the architecture has intentionally changed.

Then update the appropriate source.
