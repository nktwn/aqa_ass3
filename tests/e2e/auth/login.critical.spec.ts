import { test, expect } from "../../fixtures/test";
import { qaUsers } from "../../data/ecommerce-data";
import { mockEcommerceApi } from "../../helpers/mock-backend";
import { LoginPage } from "../../pages/login-page";

test.describe("Authentication", () => {
  test("valid user can log in and reach the storefront @critical @smoke @auth", async ({
    page,
  }) => {
    await mockEcommerceApi(page);
    const loginPage = new LoginPage(page);

    await loginPage.goto();
    await loginPage.login(qaUsers.customer.phoneNumber, qaUsers.customer.password);

    await expect(page).toHaveURL("/");
    await expect(
      page.getByRole("heading", { name: "Все нужные товары в одном аккуратном каталоге" }),
    ).toBeVisible();
  });

  test("invalid credentials are rejected with an error message @critical @regression @auth", async ({
    page,
  }) => {
    await mockEcommerceApi(page);
    const loginPage = new LoginPage(page);

    await loginPage.goto();
    await loginPage.login("+70000000000", "wrong-password");

    await expect(loginPage.error()).toBeVisible();
    await expect(page).toHaveURL("/login");
  });
});
