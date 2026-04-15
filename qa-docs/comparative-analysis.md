# Planned vs Actual Comparison

## Comparison Table

| Aspect | Planned (A1) | Actual (A2/Midterm) | Gap |
| --- | --- | --- | --- |
| High-risk focus | Auth, cart, checkout, orders, catalog should receive earliest automation | All five remain the focus; midterm added backend service tests and 8 new high-risk E2E cases with metadata | Coverage is still uneven across those five modules |
| Risk evidence | Risk scores based mostly on code inspection and feature analysis | Risk scores now updated from Playwright failures, runtime-dependent API gaps, coverage summaries, and slow-path anomalies | Empirical signal is stronger, but still limited by local browser sandbox constraints |
| Automation depth | Smoke-first approach with early critical-path coverage | Repository now has backend unit/integration tests plus Playwright E2E/API layers | Frontend component-level unit tests are still absent |
| Coverage visibility | Planned to estimate scope coverage qualitatively | Midterm now generates real Go coverage and module-level coverage gap tables | Coverage remains below 70% for all five high-risk backend modules |
| CI/CD | Basic CI integration expected after Assignment 2 | Workflow now installs deps, runs backend coverage, runs smoke/critical/regression/e2e/api, generates metrics, enforces gates, and uploads artifacts | Live API success still depends on Docker Compose environment readiness |
| Quality gates | Intended to use pass/fail thresholds and module representation | Midterm gates now evaluate backend pass rate, critical non-runtime failures, average coverage, flaky rate, and smoke pass rate | Thresholds are intentionally incremental, not release-grade hardening thresholds |

## Incorrect Assignment 1 Assumptions

1. Runtime availability was assumed to be a secondary concern; the current live API failures show that environment readiness itself is now a measurable quality risk.
2. Catalog risk was slightly understated at 16 because the backend coverage picture is much weaker than the Assignment 1 inspection implied.
3. Order observability was assumed to become easier once E2E coverage existed, but direct backend instrumentation remained near-zero until the midterm additions.

## Missing Test Scenarios Found During Midterm

1. Repeated rapid actions such as double add-to-cart and double cancel order.
2. Invalid-user behavior such as expired token reuse and guest favorite attempts.
3. Edge data handling such as special-character catalog searches and blank checkout address submissions.
4. Direct order-service transition rules and cancellation constraints at the backend layer.

## Inefficient Automation Decisions Discovered

1. Assignment 2 artifact generation used placeholder coverage values instead of real coverage instrumentation.
2. Live API toggles were enabled in the checked-in `.env.qa`, which created false-negative failures when backend services were not running.
3. Existing backend service tests had drifted from the current interfaces, so coverage looked broader on paper than it really was.

## Lessons Learned

1. Coverage numbers must come from instrumentation, not from intended scope.
2. Runtime-dependent suites need explicit tagging and gating so they generate evidence without polluting blocking results.
3. High-risk backend modules need direct tests even when E2E flows already exist.
4. Risk review improves substantially when failures, skips, coverage gaps, and execution time are analyzed together rather than separately.
