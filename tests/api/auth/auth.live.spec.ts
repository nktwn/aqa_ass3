import { expect, test } from "@playwright/test";
import {
  qaApiBaseUrl,
  qaApiLoginPassword,
  qaApiLoginPhoneNumber,
  shouldRunLiveApiTests,
} from "../../helpers/env";

test.describe("Live API authentication", () => {
  test.skip(!shouldRunLiveApiTests, "Set QA_RUN_API_TESTS=true when backend services are running.");

  test("invalid login returns 401 @api @critical @auth", async ({ request }) => {
    const response = await request.post(`${qaApiBaseUrl}/auth/login`, {
      data: {
        phone_number: "+70000000000",
        password: "wrong-password",
      },
    });

    expect(response.status()).toBe(401);
  });

  test.skip(
    !qaApiLoginPhoneNumber || !qaApiLoginPassword,
    "Set QA_API_LOGIN_PHONE and QA_API_LOGIN_PASSWORD for positive live auth coverage.",
  );

  test("valid seeded login returns tokens @api @regression @auth", async ({ request }) => {
    const response = await request.post(`${qaApiBaseUrl}/auth/login`, {
      data: {
        phone_number: qaApiLoginPhoneNumber,
        password: qaApiLoginPassword,
      },
    });

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.access_token).toBeTruthy();
    expect(body.refresh_token).toBeTruthy();
  });
});
