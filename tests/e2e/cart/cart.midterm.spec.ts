import { test, expect } from "../../fixtures/test";
import { mockEcommerceApi } from "../../helpers/mock-backend";
import { seedAuthSession } from "../../helpers/session";
import { attachQaMetadata } from "../../helpers/test-metadata";

test.describe("Cart Midterm", () => {
  test("rapid double add keeps cart state consistent for the same product @critical @regression @cart", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-CART-CONC-01",
      module: "Cart management",
      scenarioType: "concurrency",
      inputData: "Authenticated customer double-clicks Add on the same supplier offer for product 101",
      expectedOutput: "Cart records two units of the same product without a UI crash or duplicate supplier groups",
    });

    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, { authenticatedAs: "customer" });

    await page.goto("/product/101");
    await expect(page.getByRole("button", { name: "Выбрать количество" }).first()).toBeVisible();

    await page.evaluate(async () => {
      const token = window.localStorage.getItem("access_token");
      const request = () =>
        fetch("/api/cart/add", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
          },
          body: JSON.stringify({
            product_id: 101,
            supplier_id: 201,
            quantity: 1,
          }),
        });

      await Promise.all([request(), request()]);
    });

    await page.goto("/cart");

    await expect(page.getByText("Arabica Coffee Beans 1kg")).toBeVisible();
    await expect(page.getByText("2 шт")).toBeVisible();
  });
});
