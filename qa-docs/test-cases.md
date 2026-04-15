# Assignment 2 Test Cases

## Coverage notes

- Test cases map directly to the implemented Playwright suites under `tests/e2e/` and `tests/api/`.
- Positive and negative scenarios are both included.
- Live API-positive cases are environment-gated and clearly marked.

| Test Case ID | Module / Feature | Description | Input data | Expected result | Scenario type | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| A2-TC-001 | Authentication and authorization | Login with valid QA customer credentials | `+77000000001 / password1@` | Redirect to `/` and storefront hero visible | Positive | Implemented in UI |
| A2-TC-002 | Authentication and authorization | Login with invalid credentials | Wrong phone/password | Error text `Неверный номер телефона или пароль` | Negative | Implemented in UI |
| A2-TC-003 | Authentication and authorization | Open protected cart route while anonymous | No token | Redirect to `/login` | Negative | Implemented in UI |
| A2-TC-004 | Authentication and authorization | Customer opens admin route | Customer token | Redirect to `/` | Negative | Implemented in UI |
| A2-TC-005 | Product catalog and listing | Render catalog and apply search filter | Search `Green Tea` | Matching product remains, unrelated product disappears | Positive | Implemented in UI |
| A2-TC-006 | Product catalog and listing | Open product detail from catalog | Product card click | Product page and supplier section load | Positive | Implemented in UI |
| A2-TC-007 | Cart management | Add product to cart from product detail | Auth token, product `101`, quantity `2` | Cart shows product and updated total | Positive | Implemented in UI |
| A2-TC-008 | Cart management | Increase quantity then clear cart | Existing cart item | Total recalculates, empty-cart state shown after clear | Positive | Implemented in UI |
| A2-TC-009 | Checkout and payment initiation | Open checkout with valid cart but no address | Auth token, cart only | Payment-link button disabled and address warning shown | Negative | Implemented in UI |
| A2-TC-010 | Checkout and payment initiation | Save address and create payment link | Address `Дом / Almaty, Abay 10` | QA-safe payment link block appears | Positive | Implemented in UI |
| A2-TC-011 | Order lifecycle | View existing order list | Auth token, pending order | Order card visible on `/orders` | Positive | Implemented in UI |
| A2-TC-012 | Order lifecycle | Cancel pending order | Existing order `5101` | Status changes to `Отменён` | Positive | Implemented in UI |
| A2-TC-013 | API baseline | Invalid live login request | Wrong credentials | `401 Unauthorized` | Negative | Implemented, runs only with backend stack |
| A2-TC-014 | API baseline | Anonymous live cart access | No token | `401 Unauthorized` | Negative | Implemented, runs only with backend stack |
| A2-TC-015 | API baseline | Public product-list live request | `limit=5&offset=0` | `200 OK`, `product_list[]`, numeric `total` | Positive | Implemented, runs only with backend stack |
| A2-TC-016 | API baseline | Positive live seeded login | `QA_API_LOGIN_PHONE`, `QA_API_LOGIN_PASSWORD` | Access and refresh tokens returned | Positive | Environment-gated |
