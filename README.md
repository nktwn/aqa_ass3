# Ecommerce QA Midterm Repository

## System Overview

This repository contains the ecommerce application and the QA automation evolved from Assignment 1 and Assignment 2:

- `frontend/`: Next.js 15 storefront and account UI
- `backend/`: Go Gin API with PostgreSQL, Redis, Swagger, and Prometheus
- `tests/`: Playwright E2E and API automation, fixtures, mocks, and metadata helpers
- `scripts/`: metrics extraction, coverage summarization, and quality gate enforcement
- `qa-docs/`: risk analysis, comparative analysis, reproducibility notes, and the midterm report
- `logs/`: machine-readable execution, coverage, flaky, failure, and quality-gate outputs
- `evidence/charts/`: generated SVG charts used in the report

Core high-risk functionality covered by the midterm implementation:

- authentication and authorization
- product catalog and listing
- cart management
- checkout and payment initiation
- order lifecycle

## Midterm Additions

Midterm implementation extends the existing Assignment 2 baseline with:

- repaired and expanded Go unit/integration tests for auth, cart, checkout, product, and order services
- new high-risk Playwright cases with traceable metadata under `tests/e2e/**/**/*.midterm.spec.ts`
- real Go coverage capture and module-level coverage gap extraction
- failed-test, flaky-test, unexpected-behavior, and risk re-evaluation summaries generated from `logs/`
- upgraded CI workflow in `.github/workflows/qa-assignment2.yml`
- programmatic quality gates with evaluation output in `logs/quality-gates-summary.json`
- midterm report and comparative analysis under `qa-docs/`

## Setup

1. Install frontend dependencies:
   `cd frontend && npm ci`
2. Install root QA dependencies:
   `npm ci`
3. Install Playwright browser binaries:
   `npm run qa:install:browsers`
4. Prepare the QA environment:
   `cp .env.qa.example .env.qa`

`QA_RUN_API_SMOKE` and `QA_RUN_API_TESTS` default to `false` in the example env so local runs do not fail when the backend stack is not running.

## Common Commands

- Run backend unit/integration tests:
  `npm run qa:go:test`
- Run backend tests with coverage artifacts:
  `npm run qa:go:test:coverage`
- Run smoke suite JSON report:
  `npm run qa:test:smoke:json`
- Run critical suite JSON report:
  `npm run qa:test:critical:json`
- Run regression suite JSON report:
  `npm run qa:test:regression:json`
- Run full E2E JSON report:
  `npm run qa:test:e2e:json`
- Run API JSON report:
  `npm run qa:test:api:json`
- Generate midterm metrics artifacts:
  `npm run qa:metrics`
- Enforce quality gates:
  `npm run qa:quality-gates`
- Start the live backend stack for Assignment 3 experiments:
  `npm run qa:env:up`
- Run Assignment 3 performance experiments:
  `npm run qa:performance`
- Run Assignment 3 mutation experiments:
  `npm run qa:mutation`
- Run Assignment 3 chaos experiments:
  `npm run qa:chaos`
- Generate Assignment 3 tables, charts, and markdown summaries:
  `npm run qa:metrics:assignment3`
- Stop the live backend stack:
  `npm run qa:env:down`

## Reproducibility Notes

- Full local Playwright browser execution requires the frontend web server to bind to `127.0.0.1:3000`.
- In this sandboxed session, that bind failed with `listen EPERM: operation not permitted 127.0.0.1:3000`, so the newly added mock-backed E2E cases were implemented but not executed locally.
- Backend Go test and coverage artifacts were executed successfully in this session and are reflected in `logs/go-test.json` and `logs/go-coverage.out`.
- Live API suites are wired into CI with Docker Compose so they can run when backend infrastructure is available.

## Key Midterm Artifacts

- `qa-docs/midterm-report.md`
- `qa-docs/risk-reevaluation.md`
- `qa-docs/comparative-analysis.md`
- `qa-docs/reproducibility.md`
- `qa-docs/tables/risk-reevaluation.csv`
- `qa-docs/tables/planned-vs-actual.csv`
- `qa-docs/tables/midterm-test-cases.csv`
- `logs/midterm-metrics-summary.json`
- `logs/coverage-summary.json`
- `logs/failed-tests-summary.json`
- `logs/flaky-tests-summary.json`
- `logs/unexpected-behavior-summary.json`
- `logs/quality-gates-summary.json`

## Assignment 3 Artifacts

- `qa/performance/plan.md`
- `qa/mutation/plan.md`
- `qa/chaos/plan.md`
- `qa-docs/performance-results.md`
- `qa-docs/mutation-results.md`
- `qa-docs/chaos-results.md`
- `qa-docs/assignment3-experimental-engineering-report.md`
- `logs/performance-summary.json`
- `logs/mutation-results.json`
- `logs/chaos-summary.json`
- `qa-docs/tables/assignment3-*.csv`
- `evidence/charts/performance-response-time.svg`
- `evidence/charts/performance-throughput.svg`
- `evidence/charts/mutation-score.svg`
- `evidence/charts/chaos-availability.svg`
