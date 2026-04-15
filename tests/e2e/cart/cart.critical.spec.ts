import { test, expect } from "../../fixtures/test";
import { buildMockCart, mockEcommerceApi } from "../../helpers/mock-backend";
import { seedAuthSession } from "../../helpers/session";
import { ProductPage } from "../../pages/product-page";
import { CartPage } from "../../pages/cart-page";

test.describe("Cart", () => {
  test("authenticated user can add a product to cart from product details @critical @smoke @cart", async ({
    page,
  }) => {
    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, { authenticatedAs: "customer" });
    const productPage = new ProductPage(page);

    await productPage.goto(101);
    await productPage.addFirstOfferToCart(1);
    await page.goto("/cart");

    await expect(page.getByRole("heading", { name: /Ваша корзина/ })).toBeVisible();
    await expect(page.getByText("Arabica Coffee Beans 1kg")).toBeVisible();
    await expect(page.getByText("Общая сумма заказа: 36,000 ₸")).toBeVisible();
  });

  test("user can update quantity and clear cart @critical @regression @cart", async ({ page }) => {
    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, {
      authenticatedAs: "customer",
      initialCart: buildMockCart([{ productId: 101, supplierId: 201, quantity: 2 }]),
    });
    const cartPage = new CartPage(page);

    await cartPage.goto();
    await cartPage.increaseQuantity();
    await expect(page.getByText("Общая сумма заказа: 54,000 ₸")).toBeVisible();

    await cartPage.clear();
    await expect(page.getByRole("heading", { name: "Корзина пуста" })).toBeVisible();
  });
});
