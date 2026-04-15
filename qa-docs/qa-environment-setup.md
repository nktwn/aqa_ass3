# QA Environment Setup Report

## Overview

Assignment 1 QA assets were added at the repository root to avoid disrupting existing frontend and backend conventions. The setup is intentionally small, reproducible, and focused on smoke coverage plus documentation.

## Configured Tools

| Tool | Location | Purpose |
| --- | --- | --- |
| Playwright | root `package.json`, `playwright.config.ts` | Browser smoke tests and optional API smoke |
| Dotenv | root `package.json` | Loads `.env.qa` values for QA runs |
| GitHub Actions | `.github/workflows/qa-assignment2.yml` | Midterm CI/CD pipeline with coverage, metrics, and quality gates |
| Existing Go test runner | `go test ./...` in `backend/` | Reuses already-present backend tests as baseline evidence |

## QA Repository Structure

| Path | Purpose |
| --- | --- |
| `qa-docs/` | Assignment 1 documentation and metrics |
| `tests/e2e/` | UI smoke tests |
| `tests/api/` | API smoke tests |
| `tests/helpers/` | Shared env helpers |
| `tests/fixtures/` | Shared Playwright fixtures |
| `evidence/` | Reserved location for screenshots and exported evidence |

## Test Commands

Run from repository root unless noted otherwise.

| Command | Purpose |
| --- | --- |
| `npm install` | Install root QA dependencies |
| `npm run qa:install:browsers` | Install Playwright browser binaries |
| `npm run qa:test` | Run all Assignment 1 QA tests |
| `npm run qa:test:smoke` | Run smoke-tagged QA tests |
| `npm run qa:test:e2e` | Run only UI smoke tests |
| `npm run qa:test:api` | Run only API smoke tests |
| `GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...` | Run backend Go tests from `backend/` without host-cache permission issues |

## Environment and Configuration Requirements

### Frontend

- The Playwright config starts the frontend automatically with `npm --prefix frontend run dev -- --hostname 127.0.0.1`.
- QA base URL defaults to `http://127.0.0.1:3000`.

### Backend

- The backend expects PostgreSQL, Redis, JWT, HTTP, and Fondy-related environment variables.
- `backend/.env` already documents the backend runtime contract in the repository.
- API smoke is disabled by default because backend startup also requires its supporting services.

### QA env file

Copy `.env.qa.example` to `.env.qa` and adjust values if needed:

- `QA_BASE_URL`
- `QA_API_BASE_URL`
- `QA_FRONTEND_START_COMMAND`
- `QA_API_HEALTH_PATH`
- `QA_RUN_API_SMOKE`
- `NEXT_PUBLIC_API_BASE_URL`

## CI/CD Summary

The Assignment 1 workflow:

1. Checks out the repository
2. Sets up Node.js
3. Installs frontend and root QA dependencies
4. Installs Playwright Chromium
5. Runs frontend lint if available
6. Runs backend Go tests with writable local `GOCACHE` and `GOMODCACHE` as informational baseline evidence
7. Runs smoke QA tests
8. Uploads Playwright artifacts when present

## Known Limitations

- API smoke remains opt-in until backend infrastructure is started in the same environment.
- Existing backend Go tests are not fully green in the current repository state. Verification showed failures in `modules/auth/jwt`, `modules/auth/service/auth`, `modules/cart/service`, and `modules/product/service`.
- No seeded end-to-end login credentials were introduced because that would expand scope and touch business data assumptions.
- Payment and callback paths are documented as high risk but are not deeply automated yet.

## Next Steps for Assignment 2

- Add authenticated E2E coverage for cart and order flows
- Add API tests for catalog, auth, cart, and order endpoints
- Add test data strategy for deterministic user accounts and cart state
- Capture execution metrics such as smoke duration and pass trend over multiple runs
