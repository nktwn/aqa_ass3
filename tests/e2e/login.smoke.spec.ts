import { expect } from "@playwright/test";
import { test } from "../fixtures/test";

test("login page smoke path is reachable @smoke", async ({ page }) => {
  await page.goto("/login");

  await expect(
    page.getByRole("heading", { name: "Вход в аккаунт", exact: false }),
  ).toBeVisible();
  await expect(page.getByPlaceholder("Номер телефона")).toBeVisible();
  await expect(page.getByPlaceholder("Пароль")).toBeVisible();
});

