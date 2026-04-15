# Midterm Quality Gates

## Implemented Gates

| Gate ID | Criterion | Threshold | Actual | Status |
| --- | --- | --- | --- | --- |
| QG-MID-001 | Backend Go test pass rate | `>= 95%` | `100%` | Pass |
| QG-MID-002 | Critical non-runtime failures | `0` | `0` | Pass |
| QG-MID-003 | Average high-risk backend coverage | `>= 20%` | `21.07%` | Pass |
| QG-MID-004 | Observed flaky rate | `<= 5%` | `0%` | Pass |
| QG-MID-005 | Non-runtime smoke pass rate | `>= 85%` | `88.89%` | Pass |

## Evaluation

- Threshold fit: The thresholds are balanced. They are strict enough to detect regression pressure while excluding purely runtime-dependent API outages from blocking the branch.
- Failure cause: The remaining failures in current parsed reports were infrastructure/runtime related, not non-runtime regressions in mock-backed critical flows.
- Test sufficiency: Still insufficient in some areas because at least one high-risk backend module remains below `10%` coverage.

## Enforcement

- `scripts/check-quality-gates.mjs` reads generated midterm summaries from `logs/`.
- The script writes `logs/quality-gates-summary.json` and `qa-docs/tables/quality-gates-evaluation.csv`.
- CI fails when any blocking gate fails.
