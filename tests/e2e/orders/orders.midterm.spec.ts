import { test, expect } from "../../fixtures/test";
import { buildMockOrder, mockEcommerceApi } from "../../helpers/mock-backend";
import { seedAuthSession } from "../../helpers/session";
import { attachQaMetadata } from "../../helpers/test-metadata";

test.describe("Orders Midterm", () => {
  test("double cancellation leaves the order in a stable cancelled state @critical @regression @orders", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-ORDER-CONC-01",
      module: "Order lifecycle",
      scenarioType: "concurrency",
      inputData: "Authenticated customer rapidly double-clicks order cancellation for a pending order",
      expectedOutput: "Order ends in Cancelled state without duplicate rows or a broken detail page",
    });

    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, {
      authenticatedAs: "customer",
      initialOrders: [buildMockOrder(5001)],
      initialContracts: [],
    });

    await page.goto("/orders/5001");
    await expect(page.getByRole("button", { name: /Отменить заказ/i })).toBeVisible();

    await page.evaluate(async () => {
      const token = window.localStorage.getItem("access_token");
      const request = () =>
        fetch("/api/order/cancel", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
          },
          body: JSON.stringify({ order_id: 5001 }),
        });

      await Promise.all([request(), request()]);
    });

    await page.reload({ waitUntil: "networkidle" });

    await expect(page.getByText("Заказ #5001")).toBeVisible();
    await expect(page.getByText(/Отмен(е|ё)н/i)).toBeVisible();
    await expect(page.getByRole("button", { name: /Отменить заказ/i })).toHaveCount(0);
  });
});
