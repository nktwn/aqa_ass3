import type { Page } from "@playwright/test";

export class CartPage {
  constructor(private readonly page: Page) {}

  async goto() {
    await this.page.goto("/cart");
  }

  async increaseQuantity() {
    await this.page.getByRole("button", { name: "+" }).first().click();
  }

  async decreaseQuantity() {
    await this.page.getByRole("button", { name: "−" }).first().click();
  }

  async clear() {
    await this.page.getByRole("button", { name: /Очистить корзину/ }).click();
  }

  async goToCheckout() {
    await this.page.getByRole("button", { name: /Перейти к оформлению/ }).click();
  }
}
