import { test as base, expect } from "@playwright/test";

export { expect };

export const test = base.extend({
  appBaseUrl: async ({ baseURL }, use) => {
    await use(baseURL || "http://127.0.0.1:3000");
  },
});

