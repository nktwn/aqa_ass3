import dotenv from "dotenv";

dotenv.config({ path: ".env.qa" });

const trimTrailingSlash = (value: string) => value.replace(/\/+$/, "");

export const qaBaseUrl = trimTrailingSlash(
  process.env.QA_BASE_URL || "http://127.0.0.1:3000",
);

export const qaApiBaseUrl = trimTrailingSlash(
  process.env.QA_API_BASE_URL || "http://127.0.0.1:8080",
);

export const qaApiHealthUrl = trimTrailingSlash(
  process.env.QA_API_HEALTH_URL || "http://127.0.0.1:8080",
);

export const qaApiHealthPath = process.env.QA_API_HEALTH_PATH || "/metrics";

export const shouldRunApiSmoke = process.env.QA_RUN_API_SMOKE === "true";
export const shouldRunLiveApiTests = process.env.QA_RUN_API_TESTS === "true";
export const qaApiLoginPhoneNumber = process.env.QA_API_LOGIN_PHONE || "";
export const qaApiLoginPassword = process.env.QA_API_LOGIN_PASSWORD || "";
