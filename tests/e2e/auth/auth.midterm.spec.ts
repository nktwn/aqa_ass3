import { test, expect } from "../../fixtures/test";
import { mockEcommerceApi } from "../../helpers/mock-backend";
import { attachQaMetadata } from "../../helpers/test-metadata";

test.describe("Authentication Midterm", () => {
  test("expired token is cleared and protected cart route redirects to login @critical @regression @auth", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-AUTH-INVALID-01",
      module: "Authentication and authorization",
      scenarioType: "invalid-user-behavior",
      inputData: "Expired access token and refresh token stored in localStorage before visiting /cart",
      expectedOutput: "Client clears expired credentials and redirects the user to /login",
    });

    await page.addInitScript(() => {
      window.localStorage.setItem("access_token", "expired-access-token");
      window.localStorage.setItem("refresh_token", "expired-refresh-token");
      window.localStorage.setItem("expires_at", String(Date.now() - 60_000));
    });
    await mockEcommerceApi(page);

    await page.goto("/cart");

    await page.waitForURL(/\/login$/);
    await expect(page).toHaveURL(/\/login$/);
    await expect(page.getByRole("button", { name: "Войти" })).toBeVisible();
  });

  test("guest user attempting to favorite a product is redirected to login @critical @regression @auth", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-AUTH-FAIL-02",
      module: "Authentication and authorization",
      scenarioType: "failure",
      inputData: "Anonymous user opens product details and clicks the favorite button",
      expectedOutput: "Unauthorized favorite action is blocked and the user is redirected to /login",
    });

    await mockEcommerceApi(page);

    await page.goto("/product/101");
    await expect(page.getByRole("button", { name: /Добавить в избранное/i })).toBeVisible();
    await page.getByRole("button", { name: /Добавить в избранное/i }).click();

    await page.waitForURL(/\/login$/);
    await expect(page).toHaveURL(/\/login$/);
  });
});
