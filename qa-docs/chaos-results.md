# Chaos Results

Generated from `logs/chaos-summary.json`.

## Scenario Summary

| Scenario | Fault | Availability During Fault (%) | Recovered | MTTR (ms) |
| --- | --- | ---: | --- | ---: |
| app_downtime_catalog | API downtime | 0 | Yes | 1579 |
| database_unavailable_auth | database unavailability | 0 | Yes | 2666 |
| redis_checkout_failure | dependency failure | 0 | Yes | 2803 |
| injected_latency_proxy | injected network latency | 100 | Yes | 0 |

## Resilience Recommendations

- app_downtime_catalog: consider graceful maintenance responses or health-gated routing in front of the app container.
- redis_checkout_failure: checkout should expose clearer dependency-failure diagnostics and possibly queue/retry payment-order persistence.
- injected_latency_proxy: client-visible latency rises directly; consider timeout budgets and user-facing retry guidance.
