# Assignment 3 Mutation Testing Plan

## Scope

Mutation testing focuses on high-risk backend service logic that already has targeted Go tests from the Midterm:

| Module | Why Included |
| --- | --- |
| Authentication and authorization | Login and registration are critical control points with direct business and security impact. |
| Cart management | Cart and checkout validation remain the most empirically risky customer flow. |
| Product catalog and listing | Coverage is still thin, making it a strong candidate for survivable mutants. |
| Order lifecycle | Order status and authorization errors have high operational impact. |

## Mutant Strategy

Controlled source-level mutants are defined in `qa/mutation/mutants.json` and executed by a custom harness because no mutation framework is already established in this stack. The harness:

1. Applies one mutant at a time.
2. Runs the relevant existing Go test package.
3. Records whether the test suite kills or misses the mutant.
4. Restores the original source file before continuing.

## Mutant Categories Used

- logical operator changes
- constant changes
- modified return values

## Output Artifacts

- Raw execution logs: `logs/mutation/*.log`
- Machine-readable results: `logs/mutation-results.json`
- Generated tables: `qa-docs/tables/assignment3-mutation-summary.csv`, `qa-docs/tables/assignment3-mutation-scores.csv`
- Generated chart: `evidence/charts/mutation-score.svg`

## Review Standard

Surviving mutants are treated as evidence of test design gaps, not as successful product behavior. Recommendations in the Assignment 3 report are derived only from actual mutant outcomes.
