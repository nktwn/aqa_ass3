import { test, expect } from "../../fixtures/test";
import { buildMockOrder, mockEcommerceApi } from "../../helpers/mock-backend";
import { seedAuthSession } from "../../helpers/session";
import { OrdersPage } from "../../pages/orders-page";

test.describe("Orders", () => {
  test("orders page shows created order details and supports cancellation @critical @smoke @orders", async ({
    page,
  }) => {
    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, {
      authenticatedAs: "customer",
      initialOrders: [buildMockOrder(5101)],
    });
    const ordersPage = new OrdersPage(page);

    await ordersPage.goto();

    await expect(page.getByRole("heading", { name: /Мои заказы/ })).toBeVisible();
    await expect(page.getByText("Заказ #5101")).toBeVisible();

    await ordersPage.cancelFirstOrder();
    await expect(page.getByText("Отменён")).toBeVisible();
  });

  test("unauthorized user is redirected from orders page @critical @regression @orders", async ({
    page,
  }) => {
    await mockEcommerceApi(page);

    await page.goto("/orders");

    await expect(page).toHaveURL("/login");
  });
});
