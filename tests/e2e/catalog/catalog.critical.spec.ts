import { test, expect } from "../../fixtures/test";
import { mockEcommerceApi } from "../../helpers/mock-backend";
import { CatalogPage } from "../../pages/catalog-page";

test.describe("Catalog", () => {
  test("catalog renders products and applies search filters @critical @smoke @catalog", async ({
    page,
  }) => {
    await mockEcommerceApi(page);
    const catalogPage = new CatalogPage(page);

    await catalogPage.goto();

    await expect(page.getByRole("heading", { name: "Каталог товаров" })).toBeVisible();
    await expect(page.getByText("Arabica Coffee Beans 1kg")).toBeVisible();

    await catalogPage.searchFor("Green Tea");

    await expect(page).toHaveURL(/q=Green\+Tea|q=Green%20Tea/);
    await expect(page.getByText("Premium Green Tea Set")).toBeVisible();
    await expect(page.getByText("Arabica Coffee Beans 1kg")).toHaveCount(0);
  });

  test("user can open product details from the catalog @critical @regression @catalog", async ({
    page,
  }) => {
    await mockEcommerceApi(page);
    const catalogPage = new CatalogPage(page);

    await catalogPage.goto();
    await catalogPage.openProduct("Arabica Coffee Beans 1kg");

    await expect(page).toHaveURL("/product/101");
    await expect(page.getByRole("heading", { name: "Arabica Coffee Beans 1kg" })).toBeVisible();
    await expect(page.getByRole("heading", { name: "Поставщики" })).toBeVisible();
  });
});
