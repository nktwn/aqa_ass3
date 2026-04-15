# Risk Assessment

## Scope

This Assignment 1 risk assessment is based on the current repository state as of 2026-04-10. The project is a separated ecommerce web application with:

- `frontend/`: Next.js 15 + React 19 storefront and account UI.
- `backend/`: Go 1.23/1.24 Gin API with PostgreSQL, Redis, Docker Compose, Prometheus, and Swagger.

## Scoring Model

- Probability: 1 (rare) to 5 (very likely)
- Impact: 1 (minor) to 5 (critical)
- Risk score = Probability x Impact
- Priority bands:
  - High: 15-25
  - Medium: 8-14
  - Low: 1-7

## Assumptions

- Product availability, supplier logic, payment callback behavior, and seeded users are inferred from code structure because a full seeded runtime was not available inside this session.
- Fondy is the intended payment integration because the backend requires `FONDY_*` variables.
- Admin product management is present but appears secondary to core customer purchase flow.

## System Overview

The main customer flow is: register or log in, browse catalog/product pages, add items to cart grouped by supplier, manage delivery address, create checkout/payment link, then monitor orders and contracts. Supporting flows include favorites, profile management, notifications, analytics, and admin product CRUD.

## Risk Table

| Module | Description | Failure scenario | Business impact | Probability | Impact | Score | Priority | Rationale |
| --- | --- | --- | --- | ---: | ---: | ---: | --- | --- |
| Authentication and authorization | Login, registration, JWT handling, route protection, role access | Users cannot sign in, sessions expire incorrectly, or protected/admin routes leak access | Revenue loss, account lockout, security exposure | 4 | 5 | 20 | High | All protected ecommerce journeys depend on valid auth; backend and frontend both enforce it |
| Cart management | Add/remove items, grouped supplier cart totals, quantity changes, cart retrieval | Wrong quantities, broken totals, or cart cannot be loaded | Direct conversion loss and poor trust | 4 | 5 | 20 | High | Cart is central to purchase conversion and spans UI, API, and persistence |
| Checkout and payment initiation | Address selection, minimum order validation, payment link creation, payment callback handoff | Orders cannot be paid, wrong address used, or payment link is not generated | Immediate revenue loss and abandoned orders | 4 | 5 | 20 | High | Multi-step flow with external payment integration and supplier rules |
| Order lifecycle | Order creation, retrieval, status updates, cancel flow, history pages | Orders are missing, duplicated, or stuck in wrong status | Revenue reconciliation issues and support load | 4 | 5 | 20 | High | Post-purchase flow affects fulfillment, customer trust, and supplier operations |
| Product catalog and listing | Homepage featured products, catalog listing, pagination, sorting, filtering | Users cannot discover products or see incorrect list data | Reduced conversion and discoverability | 4 | 4 | 16 | High | Core browsing entry point used by homepage and catalog page |
| Product detail and favorites | Product page, favorite status, add/remove favorite, reviews | Product information is missing or personalization fails | Medium conversion impact and weaker retention | 3 | 4 | 12 | Medium | Important for shopper confidence but not always blocking initial browse |
| User profile and address book | Profile retrieval/update, address create/list, role lookup | Users cannot maintain delivery details or profile data | Checkout failures and account frustration | 3 | 4 | 12 | Medium | Address data directly feeds checkout; profile issues are less catastrophic than auth/cart |
| Search and filters | Query, price filter, sort controls, query-string state | Users cannot narrow catalog efficiently | Lower findability and degraded UX | 3 | 3 | 9 | Medium | High user-facing value but less critical than checkout flow |
| Notifications | Notification list, unread counts, mark-as-read | Operational notifications are stale or unread counts wrong | Reduced transparency, limited operational impact | 2 | 3 | 6 | Low | Useful supporting flow, not core to initial revenue path |
| Contracts and signatures | Contract retrieval and signature actions in order flow | Required signatures block fulfillment or status progression | Fulfillment delay and support escalation | 3 | 4 | 12 | Medium | Important post-order control, especially for supplier/customer coordination |
| Admin product management | Admin-only create/update/delete product endpoints and page | Admin cannot maintain catalog or unauthorized users gain access | Catalog maintenance delays or security issue | 2 | 4 | 8 | Medium | Needed for operations, but less frequent than customer journeys |
| Observability and metrics | Prometheus metrics, Swagger docs, pipeline artifacts | Failures are harder to diagnose and regressions linger | Slower recovery and weaker QA evidence | 2 | 3 | 6 | Low | Operationally valuable but indirect customer impact |

## Highest-Priority Areas

The highest-priority ecommerce areas for Assignment 1 are:

1. Authentication and authorization
2. Cart management
3. Checkout and payment initiation
4. Order lifecycle
5. Product catalog and listing

These modules represent the minimum viable revenue path and should receive first-pass smoke coverage plus the earliest deeper tests in Assignment 2.

