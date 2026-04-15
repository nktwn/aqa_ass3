import type { Page } from "@playwright/test";
import { qaUsers, type QaUserKey } from "../data/ecommerce-data";

export async function seedAuthSession(page: Page, userKey: QaUserKey = "customer") {
  const user = qaUsers[userKey];
  const expiresAt = Date.now() + 6 * 60 * 60 * 1000;

  await page.addInitScript(
    ({ accessToken, refreshToken, expiresAt: nextExpiresAt }) => {
      window.localStorage.setItem("access_token", accessToken);
      window.localStorage.setItem("refresh_token", refreshToken);
      window.localStorage.setItem("expires_at", String(nextExpiresAt));
    },
    {
      accessToken: user.accessToken,
      refreshToken: user.refreshToken,
      expiresAt,
    },
  );
}
