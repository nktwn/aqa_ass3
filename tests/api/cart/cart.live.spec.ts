import { expect, test } from "@playwright/test";
import { qaApiBaseUrl, shouldRunLiveApiTests } from "../../helpers/env";

test.describe("Live API cart", () => {
  test.skip(!shouldRunLiveApiTests, "Set QA_RUN_API_TESTS=true when backend services are running.");

  test("unauthorized cart access is rejected @api @critical @cart", async ({ request }) => {
    const response = await request.get(`${qaApiBaseUrl}/cart/`);

    expect(response.status()).toBe(401);
  });
});
