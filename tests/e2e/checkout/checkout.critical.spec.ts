import { test, expect } from "../../fixtures/test";
import { buildMockCart, mockEcommerceApi } from "../../helpers/mock-backend";
import { seedAuthSession } from "../../helpers/session";
import { CheckoutPage } from "../../pages/checkout-page";

test.describe("Checkout", () => {
  test("checkout requires a delivery address before payment link creation @critical @regression @checkout", async ({
    page,
  }) => {
    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, {
      authenticatedAs: "customer",
      initialCart: buildMockCart([{ productId: 101, supplierId: 201, quantity: 2 }]),
    });
    const checkoutPage = new CheckoutPage(page);

    await checkoutPage.goto();

    await expect(page.getByText("Адрес пока не выбран")).toBeVisible();
    await expect(page.getByRole("button", { name: "Создать платёжную ссылку" })).toBeDisabled();
  });

  test("user can save address and create payment link for valid cart @critical @smoke @checkout", async ({
    page,
  }) => {
    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, {
      authenticatedAs: "customer",
      initialCart: buildMockCart([{ productId: 101, supplierId: 201, quantity: 2 }]),
    });
    const checkoutPage = new CheckoutPage(page);

    await checkoutPage.goto();
    await checkoutPage.saveAddress("Дом", "Almaty, Abay 10");
    await checkoutPage.createPaymentLink();

    await expect(page.getByText("Платёж готов к подтверждению")).toBeVisible();
    await expect(page.getByRole("link", { name: "Подтвердить тестовую оплату" })).toBeVisible();
  });
});
