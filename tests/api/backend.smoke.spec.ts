import { expect, test } from "@playwright/test";
import { qaApiHealthPath, qaApiHealthUrl, shouldRunApiSmoke } from "../helpers/env";

test.skip(!shouldRunApiSmoke, "Set QA_RUN_API_SMOKE=true when the backend stack is running.");

test("backend smoke endpoint responds successfully @smoke", async ({ request }) => {
  const response = await request.get(`${qaApiHealthUrl}${qaApiHealthPath}`);

  expect(response.ok()).toBeTruthy();
});
