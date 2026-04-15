import { test, expect } from "../../fixtures/test";
import { mockEcommerceApi } from "../../helpers/mock-backend";
import { attachQaMetadata } from "../../helpers/test-metadata";
import { CatalogPage } from "../../pages/catalog-page";

test.describe("Catalog Midterm", () => {
  test("special characters search returns an empty-state instead of crashing @critical @regression @catalog", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-CATALOG-EDGE-01",
      module: "Product catalog and listing",
      scenarioType: "edge",
      inputData: "Search query with special characters only: <script>alert(1)</script>",
      expectedOutput: "Catalog stays responsive and shows the empty-state message for no matching products",
    });

    await mockEcommerceApi(page);
    const catalogPage = new CatalogPage(page);

    await catalogPage.goto();
    await catalogPage.searchFor("<script>alert(1)</script>");

    await expect(page.getByText("По текущим фильтрам товары не найдены.")).toBeVisible();
  });

  test("extreme price filter yields no results and keeps the filter state visible @critical @regression @catalog", async ({
    page,
  }) => {
    await attachQaMetadata(test.info(), {
      testId: "TC-CATALOG-EDGE-02",
      module: "Product catalog and listing",
      scenarioType: "edge",
      inputData: "Minimum price filter set to 999999 while using the real catalog page filters",
      expectedOutput: "No products are shown and the high min-price value remains applied in the UI",
    });

    await mockEcommerceApi(page);

    await page.goto("/catalog");
    await page.getByPlaceholder("Цена от").fill("999999");
    await page.getByRole("button", { name: "Применить фильтры" }).click();

    await expect(page.getByText("По текущим фильтрам товары не найдены.")).toBeVisible();
    await expect(page.getByPlaceholder("Цена от")).toHaveValue("999999");
  });
});
