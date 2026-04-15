import type { Page } from "@playwright/test";

export class LoginPage {
  constructor(private readonly page: Page) {}

  async goto() {
    await this.page.goto("/login");
  }

  async login(phoneNumber: string, password: string) {
    await this.page.getByPlaceholder("Номер телефона").fill(phoneNumber);
    await this.page.getByPlaceholder("Пароль").fill(password);
    await this.page.getByRole("button", { name: "Войти" }).click();
  }

  error() {
    return this.page.getByText("Неверный номер телефона или пароль");
  }
}
