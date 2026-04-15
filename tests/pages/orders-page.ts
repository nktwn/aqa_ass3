import type { Page } from "@playwright/test";

export class OrdersPage {
  constructor(private readonly page: Page) {}

  async goto() {
    await this.page.goto("/orders");
  }

  async cancelFirstOrder() {
    await this.page.getByRole("button", { name: /Отменить заказ/ }).first().click();
  }
}
