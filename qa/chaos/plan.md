# Assignment 3 Chaos / Fault Injection Plan

## Scope

Chaos scenarios reuse the Midterm high-risk prioritization and focus on failures that can be injected safely and reproducibly in local Docker Compose or CI:

| Scenario | Fault | Target Capability | Injection Technique |
| --- | --- | --- | --- |
| `app_downtime_catalog` | API downtime | Catalog availability | `docker compose stop app` / `up -d app` |
| `database_unavailable_auth` | Database unavailability | Login path | `docker compose stop pg` / `up -d pg` |
| `redis_checkout_failure` | Dependency failure | Checkout persistence | `docker compose stop redis` / `up -d redis` |
| `injected_latency_proxy` | Network latency | Login and catalog flows | Local delay proxy forwarding to the healthy API |

## Observations to Capture

- Availability during the fault window
- Recovery behavior after service restoration
- Error propagation and visible status codes
- MTTR where measurable
- Graceful degradation or lack thereof

## Safety Controls

- One isolated fault at a time
- Automatic restart for stopped services
- No destructive data reset inside chaos steps
- Repeated probes use a seeded customer account and idempotent cart cleanup

## Output Artifacts

- Raw probe log: `logs/chaos-events.ndjson`
- Scenario summary: `logs/chaos-summary.json`
- Generated tables: `qa-docs/tables/assignment3-chaos-summary.csv`
- Generated chart: `evidence/charts/chaos-availability.svg`

## Known Boundaries

- Resource exhaustion is not force-injected in this repository because it would require broader host-level controls than the current local/CI setup safely exposes.
- The implemented latency scenario is a client-to-API delay proxy, which is honest and reproducible but not equivalent to full network emulation with `tc` or service-mesh tooling.
