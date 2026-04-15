# Midterm CI/CD Integration

## Workflow

The main QA workflow is `.github/workflows/qa-assignment2.yml`.

On every push, pull request, or manual dispatch it now:

1. checks out the repository
2. installs Node and Go toolchains
3. installs frontend and root dependencies
4. installs Playwright Chromium
5. prepares `.env.qa`
6. runs frontend lint
7. runs backend Go tests with coverage
8. runs smoke, critical, regression, and full E2E Playwright suites with JSON reports
9. starts the backend stack with Docker Compose
10. waits for `http://127.0.0.1:8080/metrics`
11. runs live API Playwright coverage
12. generates metrics artifacts
13. enforces quality gates
14. uploads artifacts
15. dumps backend logs on failure and tears the stack down

## Uploaded Artifacts

- `playwright-report/`
- `test-results/`
- `logs/`
- `evidence/charts/`
- `qa-docs/tables/`

## Why the Pipeline Was Changed

Assignment 2 already had useful smoke/critical/regression automation, but the midterm required:

- backend/unit/integration execution
- real coverage outputs
- metrics summaries based on logs instead of placeholders
- artifact publication for failure, flaky, and coverage analysis
- programmatic quality-gate enforcement
