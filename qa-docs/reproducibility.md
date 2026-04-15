# Reproducibility Instructions

## Local Reproduction

1. Install dependencies:
   `cd frontend && npm ci`
   `cd .. && npm ci`
2. Install Playwright Chromium:
   `npm run qa:install:browsers`
3. Prepare QA env:
   `cp .env.qa.example .env.qa`
4. Run backend tests with coverage:
   `npm run qa:go:test:coverage`
5. Run Playwright reports as needed:
   `npm run qa:test:smoke:json`
   `npm run qa:test:critical:json`
   `npm run qa:test:regression:json`
   `npm run qa:test:e2e:json`
6. Generate artifacts and evaluate gates:
   `npm run qa:metrics`
   `npm run qa:quality-gates`

## Live API Reproduction

1. Start the backend stack:
   `cd backend && docker compose up -d pg redis migrator data-seeder app`
2. Wait for metrics readiness:
   `curl -fsS http://127.0.0.1:8080/metrics`
3. Enable live API toggles in `.env.qa` or use env overrides:
   `QA_RUN_API_SMOKE=true QA_RUN_API_TESTS=true npm run qa:test:api:json`
4. Stop the stack when finished:
   `cd backend && docker compose down -v`

## Known Local Limitation in This Session

- The sandbox blocked the frontend Playwright web server from binding to `127.0.0.1:3000` with `listen EPERM: operation not permitted 127.0.0.1:3000`.
- Because of that, the new mock-backed midterm E2E cases were implemented and wired into CI but not executed locally in this session.
- Backend Go tests, coverage generation, metrics extraction, and quality-gate evaluation were executed successfully here.

## Generated Outputs

- Machine-readable logs: `logs/*.json`, `logs/go-test.json`, `logs/go-coverage.out`
- Tables: `qa-docs/tables/*.csv`
- Charts: `evidence/charts/*.svg`
- Human-readable report: `qa-docs/midterm-report.md`
