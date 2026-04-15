import type { Page } from "@playwright/test";

export class CheckoutPage {
  constructor(private readonly page: Page) {}

  async goto() {
    await this.page.goto("/checkout");
  }

  async saveAddress(description: string, street: string) {
    await this.page.getByPlaceholder("Например: Дом, офис, склад").fill(description);
    await this.page
      .getByPlaceholder("Город, улица, дом, квартира, ориентир")
      .fill(street);
    await this.page.getByRole("button", { name: /Сохранить адрес/ }).click();
  }

  async createPaymentLink() {
    await this.page.getByRole("button", { name: "Создать платёжную ссылку" }).click();
  }
}
