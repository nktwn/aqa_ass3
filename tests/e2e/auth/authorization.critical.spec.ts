import { test, expect } from "../../fixtures/test";
import { mockEcommerceApi } from "../../helpers/mock-backend";
import { seedAuthSession } from "../../helpers/session";

test.describe("Authorization", () => {
  test("unauthorized user is redirected from protected cart route @critical @smoke @auth", async ({
    page,
  }) => {
    await mockEcommerceApi(page);

    await page.goto("/cart");

    await expect(page).toHaveURL("/login");
    await expect(page.getByRole("heading", { name: /Вход в аккаунт/ })).toBeVisible();
  });

  test("customer is redirected away from admin route @critical @regression @auth", async ({
    page,
  }) => {
    await seedAuthSession(page, "customer");
    await mockEcommerceApi(page, { authenticatedAs: "customer" });

    await page.goto("/admin/products");

    await expect(page).toHaveURL("/");
    await expect(
      page.getByRole("heading", { name: "Все нужные товары в одном аккуратном каталоге" }),
    ).toBeVisible();
  });
});
