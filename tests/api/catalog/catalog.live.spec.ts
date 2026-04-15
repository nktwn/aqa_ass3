import { expect, test } from "@playwright/test";
import { qaApiBaseUrl, shouldRunLiveApiTests } from "../../helpers/env";

test.describe("Live API catalog", () => {
  test.skip(!shouldRunLiveApiTests, "Set QA_RUN_API_TESTS=true when backend services are running.");

  test("product list endpoint returns catalog payload @api @critical @catalog", async ({
    request,
  }) => {
    const response = await request.get(`${qaApiBaseUrl}/product/list?limit=5&offset=0`);

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(Array.isArray(body.product_list)).toBeTruthy();
    expect(typeof body.total).toBe("number");
  });
});
