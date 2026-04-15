import { test, expect } from "../../fixtures/test";
import { buildMockCart, mockEcommerceApi } from "../../helpers/mock-backend";
import { seedAuthSession } from "../../helpers/session";
import { attachQaMetadata } from "../../helpers/test-metadata";
import { CheckoutPage } from "../../pages/checkout-page";

test.describe("Checkout Midterm", () => {
  test("empty cart cannot generate a payment link even when an address exists @critical @regression @checkout", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-CHECKOUT-FAIL-01",
      module: "Checkout and payment initiation",
      scenarioType: "failure",
      inputData: "Authenticated customer opens checkout with a saved address but an empty cart and clicks payment link creation",
      expectedOutput: "Checkout flow rejects the action with an invalid cart message",
    });

    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, {
      authenticatedAs: "customer",
      initialCart: buildMockCart([]),
      initialAddresses: [{ id: 1, description: "Home", street: "Almaty, Abay 10" }],
    });
    const checkoutPage = new CheckoutPage(page);

    await checkoutPage.goto();
    await checkoutPage.createPaymentLink();

    await expect(page.getByText("invalid cart")).toBeVisible();
  });

  test("empty address fields are blocked before checkout submission @critical @regression @checkout", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-CHECKOUT-EDGE-02",
      module: "Checkout and payment initiation",
      scenarioType: "edge",
      inputData: "Authenticated customer tries to save a new address with blank description and street",
      expectedOutput: "Client-side validation shows an address form error and does not persist a blank address",
    });

    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, {
      authenticatedAs: "customer",
      initialCart: buildMockCart([{ productId: 101, supplierId: 201, quantity: 2 }]),
    });
    const checkoutPage = new CheckoutPage(page);

    await checkoutPage.goto();
    await page.getByRole("button", { name: "Сохранить адрес" }).click();

    await expect(page.getByText("Заполните название адреса и полный адрес доставки.")).toBeVisible();
  });
});
