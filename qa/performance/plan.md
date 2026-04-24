# Assignment 3 Performance Testing Plan

## Scope

Performance experiments reuse the Midterm risk re-evaluation in `qa-docs/risk-reevaluation.md` and focus on the three most practical high-risk areas for live stack execution:

| Module | Midterm Updated Risk Score | Reason for Assignment 3 Selection |
| --- | ---: | --- |
| Cart management | 25 | Highest empirical risk and weakest coverage among core purchase flows. |
| Authentication and authorization | 20 | System entry point with high impact and availability sensitivity. |
| Product catalog and listing | 20 | Heavy-read customer path with low direct backend coverage. |

## Target Endpoints and Flows

| Flow ID | Module | Endpoint(s) | Request Pattern | Success Criteria |
| --- | --- | --- | --- | --- |
| `auth_login` | Authentication and authorization | `POST /api/auth/login` | Seeded customer login | `200 OK` with access and refresh tokens |
| `catalog_browse` | Product catalog and listing | `GET /api/product/list?limit=12&offset=0` | Anonymous read | `200 OK` with `product_list` and `total` |
| `cart_session` | Cart management | `DELETE /api/cart/clear`, `POST /api/cart/add`, `GET /api/cart/`, `DELETE /api/cart/clear` | Authenticated short cart session | All steps complete without 5xx responses |

## Scenarios

| Scenario | Type | Concurrency | Duration | Main Objective |
| --- | --- | ---: | ---: | --- |
| `normal_load` | Normal | 4 concurrent users | 15s | Confirm baseline responsiveness under steady traffic |
| `peak_load` | Peak | 8 concurrent users | 15s | Observe saturation under sustained elevated demand |
| `spike_load` | Spike | 3 -> 12 -> 4 concurrent users | 12s | Evaluate burst absorption and stabilization |
| `endurance_load` | Endurance | 3 concurrent users | 30s | Detect degradation over a longer steady run |

## Metrics Collected

- Average response time
- Median response time
- p95 response time
- Throughput (requests per second)
- Error rate
- Best-effort container CPU, memory, and block I/O snapshots when Docker statistics are accessible
- Prometheus `/metrics` snapshots before and after each scenario

## Thresholds

Thresholds are intentionally incremental and reviewable rather than production-SLO claims. They are encoded in `qa/performance/scenarios.json` and evaluated by the runner for each scenario.

## Execution Artifacts

- Raw request log: `logs/performance-raw.ndjson`
- Scenario summary: `logs/performance-summary.json`
- Resource samples: `logs/performance-resources.json`
- Generated tables: `qa-docs/tables/assignment3-performance-summary.csv`
- Generated charts: `evidence/charts/performance-response-time.svg`, `evidence/charts/performance-throughput.svg`

## Limitations

- CPU, memory, and block I/O are best-effort Docker snapshots, not full APM traces.
- The current repository exposes Prometheus metrics but not a ready-made historical dashboard export pipeline for Assignment 3; the implemented process therefore preserves snapshots and derived summaries under `logs/`.
