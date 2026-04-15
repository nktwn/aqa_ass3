import type { Page } from "@playwright/test";

export class CatalogPage {
  constructor(private readonly page: Page) {}

  async goto() {
    await this.page.goto("/catalog");
  }

  async searchFor(value: string) {
    await this.page.getByPlaceholder("Поиск по названию").fill(value);
    await this.page.getByRole("button", { name: "Применить фильтры" }).click();
  }

  async openProduct(name: string) {
    await this.page.getByText(name).first().click();
  }
}
