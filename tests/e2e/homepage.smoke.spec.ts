import { expect } from "@playwright/test";
import { test } from "../fixtures/test";

test("homepage smoke path renders core entry points @smoke", async ({ page }) => {
  await page.goto("/");

  await expect(
    page.getByRole("heading", { name: "Все нужные товары в одном аккуратном каталоге" }),
  ).toBeVisible();
  await expect(page.getByRole("link", { name: "Открыть каталог" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Перейти к checkout" })).toBeVisible();
});

