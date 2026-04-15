# Midterm Risk Re-Evaluation

## Evidence Base

This re-evaluation uses empirical evidence from:

- `logs/smoke-report.json`
- `logs/critical-report.json`
- `logs/regression-report.json`
- `logs/go-test.json`
- `logs/go-coverage.out`
- `logs/coverage-summary.json`
- `logs/failed-tests-summary.json`
- `logs/unexpected-behavior-summary.json`

## Updated High-Risk Assessment

| Module | Original Risk Score | Observed Issues from Automation | Likelihood | Impact | Detectability | Updated Risk Score | Justification |
| --- | ---: | --- | ---: | ---: | ---: | ---: | --- |
| Authentication and authorization | 20 | 1 runtime-dependent live API auth failure, 11 uncovered backend files, 6 performance anomalies in auth-adjacent flows | 4 | 5 | 4 | 20 | Auth has the strongest backend coverage of the five high-risk modules at 43%, but live API availability and several slower auth-path tests still keep it high risk. |
| Cart management | 20 | 1 runtime-dependent live API cart failure, 19 uncovered backend files, 1 slow cart-path test | 5 | 5 | 2 | 25 | Cart coverage is only 14.19%, and the current evidence still leaves many repository/client branches unexecuted. |
| Checkout and payment initiation | 20 | 4 uncovered checkout-related backend files, no direct non-runtime failures in the current parsed runs | 3 | 5 | 3 | 15 | Checkout service coverage improved to 33.56%, so detectability is better than cart/orders, but handler/payment client gaps remain. |
| Order lifecycle | 20 | 17 uncovered backend files, low backend coverage despite new order-service tests | 4 | 5 | 2 | 20 | New tests moved orders above zero, but 9.02% coverage still leaves many order retrieval and detail branches lightly instrumented. |
| Product catalog and listing | 16 | 1 runtime-dependent live API catalog failure, 17 uncovered backend files, 1 slow catalog flow | 5 | 4 | 2 | 20 | Catalog remains highly exposed because backend coverage is only 5.58% and most handler/repository paths still rely on E2E visibility instead of direct instrumentation. |

## Key Re-Evaluation Takeaways

1. The risk model from Assignment 1 still points at the right top five modules.
2. Empirical data shifted the highest current concern toward `Cart management`, where coverage remains thin and live API runtime evidence also showed failures.
3. `Checkout and payment initiation` improved most from the new backend tests and now has the best midterm balance between impact and detectability after auth.
4. `Order lifecycle` and `Product catalog and listing` remain under-instrumented despite being revenue-critical.

## Source Tables

- `qa-docs/tables/risk-reevaluation.csv`
- `qa-docs/tables/module-coverage.csv`
- `qa-docs/tables/coverage-gaps.csv`
