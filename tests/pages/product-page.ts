import type { Page } from "@playwright/test";

export class ProductPage {
  constructor(private readonly page: Page) {}

  async goto(productId: number) {
    await this.page.goto(`/product/${productId}`);
  }

  async addFirstOfferToCart(quantityAdjustments = 0) {
    const supplierCard = this.page.locator("li").filter({
      has: this.page.getByText("Roastery One"),
    }).first();

    await supplierCard.getByRole("button", { name: "Выбрать количество" }).click();
    for (let index = 0; index < quantityAdjustments; index += 1) {
      await supplierCard.getByRole("button", { name: "+" }).click();
    }
    await supplierCard.getByRole("button", { name: "Добавить" }).click();
  }
}
