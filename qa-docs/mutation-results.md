# Mutation Results

Generated from `logs/mutation-results.json`.

## Module Scores

| Module | Total | Killed | Survived | Score (%) |
| --- | ---: | ---: | ---: | ---: |
| Authentication and authorization | 2 | 2 | 0 | 100 |
| Cart management | 1 | 0 | 1 | 0 |
| Product catalog and listing | 2 | 1 | 1 | 50 |
| Order lifecycle | 2 | 2 | 0 | 100 |
| Overall | 7 | 5 | 2 | 71.43 |

## Surviving Mutants

- CART-001 (Cart management): Require cart totals to strictly exceed the supplier minimum instead of meeting it.
- PRODUCT-002 (Product catalog and listing): Return nil rather than an explicit empty suggestion list for blank input.
