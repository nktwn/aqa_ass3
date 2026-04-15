# QA Test Strategy

## 1. Project Scope and Objectives

This strategy covers Assignment 1 for the existing ecommerce repository. The immediate objective is to establish a reproducible QA foundation, prioritize the most business-critical risks, and create a small but runnable smoke suite that can be expanded in later assignments.

Primary objectives:

- Validate that the storefront can render core entry pages.
- Create a lightweight automation baseline around the highest-value user paths.
- Align documentation, metrics, and CI with the current codebase.

## 2. System Overview

### Architecture

- Frontend: Next.js 15 app-router application in `frontend/`
- Backend: Go Gin API in `backend/`
- Data stores: PostgreSQL and Redis via `backend/docker-compose.yaml`
- Supporting tooling: Swagger docs, Prometheus metrics, Dockerfiles, existing Go unit tests
- Repository shape: separated frontend/backend, not a monolith

### Core Ecommerce Features Found

- Authentication and registration
- Product catalog, featured products, product details
- Search, filters, sort, pagination
- Favorites and reviews
- Cart with supplier grouping and quantity changes
- Checkout with address selection and payment-link creation
- Order history, cancel flow, status progression
- User profile and address management
- Notifications
- Contract/signature flow
- Admin product management

## 3. Risk Assessment Results

High-priority modules for early QA focus:

| Priority order | Module | Score | Why it comes first |
| --- | --- | ---: | --- |
| 1 | Authentication and authorization | 20 | Required for all protected user and admin actions |
| 2 | Cart management | 20 | Directly affects conversion and order accuracy |
| 3 | Checkout and payment initiation | 20 | Highest revenue sensitivity; external dependency involved |
| 4 | Order lifecycle | 20 | Critical to fulfillment and customer trust |
| 5 | Product catalog and listing | 16 | Main discovery path for buyers |

## 4. Test Approach

### Manual vs automated split

- Automated in Assignment 1:
  - UI smoke checks for homepage and login availability
  - Optional backend smoke endpoint check when API infrastructure is running
- Manual in Assignment 1:
  - Deep checkout validation
  - Payment callback handling
  - Contract/signature behavior
  - Role-specific admin and supplier workflows

### Why high-risk modules come first

The chosen priority order reflects the minimum revenue path: discover product, authenticate, manage cart, check out, and verify order handling. Failures in these modules either block purchases entirely or create immediate financial/support impact.

### Test levels emphasized

- Smoke: confirm the application is reachable and major entry points load
- Integration: leverage existing backend unit/integration-style Go tests where present
- E2E foundation: establish Playwright structure and conventions for future deeper journeys

Assignment 1 intentionally avoids broad automation breadth; the goal is a stable base rather than large test volume.

## 5. Tool Selection and Configuration

| Tool | Purpose | Selection rationale |
| --- | --- | --- |
| Playwright | UI/e2e and simple API smoke | Best fit for modern web app smoke coverage and CI artifacts |
| GitHub Actions | Basic CI pipeline | Repository already uses GitHub workflows in backend |
| Existing Go test tooling | Backend regression signal | Already present in repo and low-cost to keep using |
| Markdown + CSV deliverables | Assignment evidence | Easy to review, version, and reuse in later reports |

### Configuration summary

- Root-level `playwright.config.ts` for assignment-specific QA setup
- Dedicated folders: `tests/e2e`, `tests/api`, `tests/helpers`, `tests/fixtures`
- Root `.env.qa.example` for QA runtime variables
- Root `package.json` for consistent QA scripts without disturbing frontend package conventions

## 6. Planned Metrics

| Metric | Baseline |
| --- | --- |
| High-risk modules | 5 |
| Total modules assessed | 12 |
| Assignment 1 automated smoke coverage | 2 UI smoke tests, 1 optional API smoke test |
| Existing automated coverage already in repo | Go unit tests in auth, cart, product modules |
| Estimated Assignment 1 QA effort | 12-16 hours for setup, analysis, docs, and initial smoke coverage |

### Initial coverage plan by risk band

- High risk: smoke coverage now, deeper E2E and API scenarios in Assignment 2
- Medium risk: document and queue for targeted manual/integration testing
- Low risk: cover after core revenue path stabilizes

### Baseline effectiveness measures for later assignments

- Smoke pass rate
- Number of prioritized modules covered
- Number of defects found in high-risk flows
- Time to execute smoke suite
- CI pass/fail trend

## 7. Research-Paper Alignment

The following Assignment 1 outputs can be reused directly in later academic sections:

- Methodology inputs:
  - system overview
  - risk model
  - tool-selection rationale
  - coverage approach
- Results/baseline inputs:
  - risk matrix
  - baseline metrics
  - initial smoke scope
  - known gaps and planned expansion

This creates a clean bridge into later methodology, execution, and comparative-results sections.

