# Midterm Metrics Report

## Source of Metrics

Metrics were generated from current repository evidence with:

- `logs/smoke-report.json`
- `logs/critical-report.json`
- `logs/regression-report.json`
- `logs/go-test.json`
- `logs/go-coverage.out`
- `npm run qa:metrics`

## Current Summary

- Detected Playwright defects: `4`
- Flaky rate: `0%`
- Critical non-runtime pass rate: `100%`
- High-risk average backend coverage: `21.07%`
- Parsed Playwright execution time: `110.18s`
- Estimated pipeline runtime: `134.7s`

## High-Risk Backend Coverage

| Module | Coverage |
| --- | ---: |
| Authentication and authorization | 43% |
| Cart management | 14.19% |
| Checkout and payment initiation | 33.56% |
| Order lifecycle | 9.02% |
| Product catalog and listing | 5.58% |

## Interpretation

1. Coverage is now real and reproducible, but still well below 70% for every high-risk backend module.
2. Remaining observed failures came from runtime-dependent live API checks where the backend service was unavailable in the parsed runs.
3. The weakest remaining detectability is in `Order lifecycle` and `Product catalog and listing`.

## Related Artifacts

- `logs/midterm-metrics-summary.json`
- `logs/coverage-summary.json`
- `qa-docs/tables/module-coverage.csv`
- `qa-docs/tables/coverage-gaps.csv`
- `qa-docs/tables/midterm-metrics-summary.csv`
