# Assignment 3 Experimental Engineering Report

## 1. System and Target High-Risk Modules

This Assignment 3 implementation extends the existing Midterm QA repository rather than replacing it. The baseline already contained:

- Assignment 1 risk analysis
- Assignment 2 automation
- Midterm empirical risk re-evaluation
- Playwright E2E and API automation
- Go backend tests
- CI/CD integration
- machine-readable logs, tables, and evidence charts

The Assignment 3 experiments reused the Midterm risk evidence in `qa-docs/risk-reevaluation.md` and selected the following high-risk modules for experimental testing:

| Module | Midterm Updated Risk Score | Assignment 3 Use |
| --- | ---: | --- |
| Cart management | 25 | Primary performance target and dependency-failure target |
| Authentication and authorization | 20 | Performance target and database-failure target |
| Product catalog and listing | 20 | Performance target and downtime/latency target |
| Checkout and payment initiation | 15 | Chaos target through Redis-backed checkout persistence |
| Order lifecycle | 20 | Mutation target via supplier-status logic |

The selection stays aligned with the Midterm rather than inventing a new priority model.

## 2. Environment and Tooling

The live Assignment 3 execution in this session used:

| Item | Value |
| --- | --- |
| Execution date | April 24, 2026 |
| OS | macOS Darwin 25.4.0 arm64 |
| Node.js | `v20.19.0` |
| Go | `go1.26.2 darwin/arm64` |
| Docker Compose | `v2.31.0-desktop.2` |
| Backend stack | `pg`, `redis`, `migrator`, `data-seeder`, `app` from `backend/docker-compose.yaml` |
| Experimental scripts | `scripts/run-performance-tests.mjs`, `scripts/run-mutation-tests.mjs`, `scripts/run-chaos-tests.mjs`, `scripts/generate-assignment3-artifacts.mjs` |

The new repo-level npm commands are:

- `npm run qa:env:up`
- `npm run qa:performance`
- `npm run qa:mutation`
- `npm run qa:chaos`
- `npm run qa:metrics:assignment3`
- `npm run qa:env:down`

## 3. Performance Testing Methodology and Results

### 3.1 Method

Performance scenarios are defined in `qa/performance/scenarios.json` and documented in `qa/performance/plan.md`. The selected live flows were:

- `POST /api/auth/login`
- `GET /api/product/list?limit=12&offset=0`
- authenticated cart clear/add/get/clear sequence

Four scenarios were executed:

| Scenario | Concurrency | Duration |
| --- | ---: | ---: |
| `normal_load` | 4 users | 15s |
| `peak_load` | 8 users | 15s |
| `spike_load` | 3 -> 12 -> 4 users | 12s |
| `endurance_load` | 3 users | 30s |

Metrics collected:

- average response time
- median response time
- p95 response time
- throughput
- error rate
- best-effort Docker CPU, memory, and block I/O snapshots

### 3.2 Results

| Scenario | Avg (ms) | Median (ms) | p95 (ms) | Throughput (rps) | Error Rate (%) | Threshold Result |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| `normal_load` | 19.85 | 2 | 71 | 159.67 | 0.38 | Pass |
| `peak_load` | 24.97 | 5 | 95 | 277.73 | 1.25 | Pass |
| `spike_load` | 34.57 | 8 | 114 | 167.83 | 2.63 | Pass |
| `endurance_load` | 21.16 | 2 | 80 | 130.97 | 0.28 | Pass |

All four scenarios met the encoded Assignment 3 thresholds. The detailed outputs are stored in:

- `logs/performance-summary.json`
- `logs/performance-raw.ndjson`
- `logs/performance-resources.json`
- `qa-docs/tables/assignment3-performance-summary.csv`
- `evidence/charts/performance-response-time.svg`
- `evidence/charts/performance-throughput.svg`

### 3.3 Bottleneck Analysis

Even though all thresholds passed, the raw logs exposed a real concurrency defect in the highest-risk module:

- 125 failed requests occurred during performance execution.
- All 125 failures came from `cart_session -> cart_add`.
- The observed backend error was `duplicate key value violates unique constraint "cart_items_pkey" (SQLSTATE 23505)`.

This means the system remains fast under the chosen load levels, but cart writes are not safely handling concurrent insert races. In practical terms, the repo now has traceable evidence of a stability weakness in the same module that Midterm already ranked highest risk.

### 3.4 Resource Observations

The Docker snapshot collection is best-effort and intentionally honest:

- `backend-app-1` memory grew from `19.37 MiB` before `normal_load` to `46.3 MiB` during `endurance_load`.
- `backend-pg-1` memory grew from `41.32 MiB` to `59.52 MiB`.
- PostgreSQL block I/O increased from `26MB / 13.1MB` to `26.3MB / 56MB` across the run.
- App CPU snapshots stayed at `0%` in the collected no-stream samples, which is a limitation of coarse point-in-time Docker statistics rather than proof of zero CPU consumption.

### 3.5 Performance Recommendations

1. Harden `cart_add` against concurrent duplicate-key insertion with an upsert or retry-safe repository strategy.
2. Add a targeted backend concurrency test around cart item creation, because the issue only surfaced under live parallel request pressure.
3. If future assignments need stronger infrastructure evidence, replace no-stream Docker stats with a sampled Prometheus or cAdvisor pipeline.

## 4. Mutation Testing Methodology and Results

### 4.1 Method

Mutation definitions live in `qa/mutation/mutants.json`. A custom harness was used because no mutation framework was already established in the current stack. For each mutant, the harness:

1. modified a real backend source file
2. executed the relevant existing Go test package
3. recorded whether the mutant was killed or survived
4. restored the original source file

Mutant categories used:

- logical operator changes
- constant changes
- modified return values

### 4.2 Results

| Module | Total Mutants | Killed | Survived | Mutation Score |
| --- | ---: | ---: | ---: | ---: |
| Authentication and authorization | 2 | 2 | 0 | 100% |
| Cart management | 1 | 1 | 0 | 100% |
| Product catalog and listing | 2 | 2 | 0 | 100% |
| Order lifecycle | 2 | 2 | 0 | 100% |
| Overall | 7 | 7 | 0 | 100% |

Artifacts:

- `logs/mutation-results.json`
- `logs/mutation/*.log`
- `qa-docs/tables/assignment3-mutation-summary.csv`
- `qa-docs/tables/assignment3-mutation-scores.csv`
- `evidence/charts/mutation-score.svg`

### 4.3 Interpretation

The current Midterm-strengthened Go test packages successfully killed all seven controlled mutants. That is a strong sign that the targeted service-layer assertions are meaningful, especially around:

- admin-role registration rejection
- invalid password handling
- checkout cart validation
- suggestion-limit defaults
- order supplier authorization
- pending-to-in-progress status transitions

There were no surviving mutants in this initial set, so the recommendation is not to claim “mutation testing is complete,” but to widen the set in future iterations to include:

- cart repository and Redis persistence branches
- handler-level error mapping
- order retrieval and product-list retrieval paths with lower Midterm coverage

## 5. Chaos / Fault Injection Methodology and Results

### 5.1 Method

Chaos scenarios are defined in `qa/chaos/scenarios.json` and documented in `qa/chaos/plan.md`. Safe, reproducible injections were used:

- `docker compose stop/up` for service outages
- a local delay proxy for injected latency

### 5.2 Results

| Scenario | Fault Type | Availability During Fault | Recovery | MTTR |
| --- | --- | ---: | --- | ---: |
| `app_downtime_catalog` | API downtime | 0% | Recovered | 1670 ms |
| `database_unavailable_auth` | database unavailability | 0% | Recovered | 2640 ms |
| `redis_checkout_failure` | dependency failure | 0% | Recovered | 2996 ms |
| `injected_latency_proxy` | injected network latency | 100% | Immediate | 0 ms |

Artifacts:

- `logs/chaos-events.ndjson`
- `logs/chaos-summary.json`
- `qa-docs/tables/assignment3-chaos-summary.csv`
- `evidence/charts/chaos-availability.svg`

### 5.3 Resilience Observations

`app_downtime_catalog`

- Catalog availability dropped completely while the app container was stopped.
- Recovery after restart was fast, but there was no graceful degradation layer.

`database_unavailable_auth`

- Login availability dropped to 0%.
- The API returned `401` with the raw PostgreSQL termination message:
  `FATAL: terminating connection due to administrator command (SQLSTATE 57P01)`.
- Recovery was automatic after PostgreSQL came back.

`redis_checkout_failure`

- Baseline checkout succeeded once the cart setup crossed the supplier minimum order amount.
- During Redis outage, checkout returned `500` with a raw dependency error:
  `dial tcp: lookup redis on 127.0.0.11:53: no such host`.
- Recovery succeeded automatically after Redis restart.

`injected_latency_proxy`

- Availability stayed at 100%.
- Average latency increased to `1271.92 ms`; p95 rose to `1341 ms`.
- The system remained functionally available, but the user-facing delay became obvious.

### 5.4 Chaos Recommendations

1. Prevent raw infrastructure errors from leaking through auth and checkout responses.
2. Add retry, fallback, or clearer failure handling around Redis-backed checkout persistence.
3. Consider routing or maintenance behavior that fails more gracefully during total app downtime.
4. Set explicit client and server timeout budgets informed by the latency-proxy results.

## 6. Expected vs Actual Behavior

| Area | Expected | Actual | Assessment |
| --- | --- | --- | --- |
| Performance | Stable latency and throughput under moderate Assignment 3 load | All threshold checks passed | Meets expectation |
| Cart under concurrency | Cart writes should remain stable | `cart_add` produced duplicate-key 500s under parallel pressure | Worse than expected |
| Mutation resistance | Critical service tests should detect seeded logic regressions | 7 of 7 mutants killed | Better than expected |
| API downtime recovery | Service should recover after restart | Recovered in ~1.67s | Meets expectation |
| DB outage handling | Failure should be visible but controlled | Recovered, but raw DB error leaked via 401 | Partially meets expectation |
| Redis outage handling | Checkout should fail and recover | Failed fast and recovered, but raw dependency error leaked | Partially meets expectation |
| Network latency | Service should remain available | 100% availability, but ~1.27s average latency | Meets expectation with UX concern |

## 7. Lessons Learned and Final Recommendations

Key lessons from the Assignment 3 extension:

1. Midterm risk prioritization was directionally correct: the most meaningful live issue again surfaced in cart management.
2. Fast average response times can hide correctness problems; the concurrent cart insert race would be invisible in a pure latency-only reading.
3. Service-layer tests are now strong enough to kill a focused set of realistic mutants.
4. The system recovers from stopped services quickly, but user-facing error hygiene during dependency failure still needs work.

Recommended next actions:

1. Fix the cart item insertion race and add a regression test specifically for concurrent `cart_add`.
2. Normalize infrastructure-failure responses so login and checkout do not leak raw backend error text.
3. Expand mutation scope into repository and handler layers where Midterm coverage remains thin.
4. If the repository continues to evolve experimentally, promote the optional Assignment 3 workflow into a repeatable scheduled or manual benchmark job.
