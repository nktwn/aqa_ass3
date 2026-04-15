# Alerting and Failure Handling

## Intent

Assignment 2 must not only run automated checks but also define what happens when those checks fail. The workflow and supporting scripts now make the failure mode explicit.

## Failure-handling rules

| Scenario | Detection mechanism | CI behavior | Team action |
| --- | --- | --- | --- |
| Critical flow fails | `qa:test:critical:json` or `qa:quality-gates` returns non-zero | Workflow fails | Inspect `logs/critical-report.json`, fix regression before merge |
| Smoke flow fails | `qa:test:smoke:json` or smoke gate fails | Workflow fails early | Review latest smoke artifact and unblock basic app reachability |
| Coverage gate below target | `criticalModulesCovered < 5` in quality-gate summary | Workflow fails | Add or repair missing high-risk coverage before release |
| Smoke run too slow | `smokeDurationMs > 180000` | Workflow fails | Investigate performance regressions or unstable setup |
| Frontend lint fails | `frontend` lint step fails | Workflow fails | Fix static issues before merging |
| Backend Go baseline fails | Go test step exits non-zero | Workflow continues | Record as baseline signal, do not block Assignment 2 flow |
| Live API tests skipped | `QA_RUN_API_TESTS=false` | Workflow remains green | Document environment limitation; enable for full-stack validation |
| Artifact generation fails | `qa:metrics` step fails | Workflow fails | Repair parsing/generation logic because evidence is incomplete |
| CI config failure | Workflow YAML or setup step fails | Workflow fails | Fix pipeline configuration before relying on results |

## Latest observed status

The latest local Assignment 2 verification on 2026-04-10 produced a passing gate summary in `logs/quality-gates-summary.json`, so no blocking alert remained open after the final rerun.
