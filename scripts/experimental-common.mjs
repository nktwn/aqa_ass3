import fs from "node:fs";
import path from "node:path";
import { spawnSync } from "node:child_process";
import http from "node:http";
import dotenv from "dotenv";

dotenv.config({ path: ".env.qa" });

export const rootDir = process.cwd();
export const logsDir = path.join(rootDir, "logs");
export const tablesDir = path.join(rootDir, "qa-docs", "tables");
export const chartsDir = path.join(rootDir, "evidence", "charts");

export const qaApiBaseUrl = String(process.env.QA_API_BASE_URL || "http://127.0.0.1:8080/api").replace(/\/+$/, "");
export const qaApiHealthUrl = String(process.env.QA_API_HEALTH_URL || "http://127.0.0.1:8080").replace(/\/+$/, "");
export const qaApiLoginPhoneNumber = process.env.QA_API_LOGIN_PHONE || "+12345671";
export const qaApiLoginPassword = process.env.QA_API_LOGIN_PASSWORD || "password1@";

export function ensureDir(directory) {
  fs.mkdirSync(directory, { recursive: true });
}

export function readJson(filePath) {
  return JSON.parse(fs.readFileSync(filePath, "utf8"));
}

export function readJsonIfExists(filePath, fallback = null) {
  if (!fs.existsSync(filePath)) {
    return fallback;
  }

  try {
    return readJson(filePath);
  } catch (error) {
    return fallback;
  }
}

export function writeJson(filePath, value) {
  ensureDir(path.dirname(filePath));
  fs.writeFileSync(filePath, `${JSON.stringify(value, null, 2)}\n`);
}

export function writeText(filePath, value) {
  ensureDir(path.dirname(filePath));
  fs.writeFileSync(filePath, value);
}

export function appendNdjson(filePath, value) {
  ensureDir(path.dirname(filePath));
  fs.appendFileSync(filePath, `${JSON.stringify(value)}\n`);
}

export function round(value, digits = 2) {
  if (!Number.isFinite(value)) {
    return 0;
  }
  return Number(value.toFixed(digits));
}

export function percentile(values, target) {
  if (!values.length) {
    return 0;
  }

  const sorted = [...values].sort((a, b) => a - b);
  const rank = Math.ceil((target / 100) * sorted.length) - 1;
  const index = Math.min(sorted.length - 1, Math.max(0, rank));
  return sorted[index];
}

export function median(values) {
  return percentile(values, 50);
}

export function average(values) {
  if (!values.length) {
    return 0;
  }
  return values.reduce((sum, value) => sum + value, 0) / values.length;
}

export function toCsv(rows) {
  return `${rows
    .map((row) => row.map((value) => `"${String(value ?? "").replaceAll('"', '""')}"`).join(","))
    .join("\n")}\n`;
}

export function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function nowIso() {
  return new Date().toISOString();
}

export function statusBucket(status) {
  if (typeof status !== "number") {
    return "error";
  }
  return `${Math.floor(status / 100)}xx`;
}

export async function timedJsonRequest({
  label,
  method = "GET",
  url,
  headers = {},
  body,
  expectedStatus,
  timeoutMs = 10000,
}) {
  const startedAt = Date.now();
  const startedAtIso = nowIso();
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const response = await fetch(url, {
      method,
      headers: {
        "content-type": "application/json",
        ...headers,
      },
      body: body === undefined ? undefined : JSON.stringify(body),
      signal: controller.signal,
    });

    const durationMs = Date.now() - startedAt;
    const text = await response.text();
    let json = null;
    try {
      json = text ? JSON.parse(text) : null;
    } catch (error) {
      json = null;
    }

    return {
      label,
      ok: expectedStatus === undefined ? response.ok : response.status === expectedStatus,
      expectedStatus: expectedStatus ?? null,
      status: response.status,
      durationMs,
      startedAt: startedAtIso,
      body: json,
      bodyPreview: text.slice(0, 300),
      error: null,
    };
  } catch (error) {
    return {
      label,
      ok: false,
      expectedStatus: expectedStatus ?? null,
      status: null,
      durationMs: Date.now() - startedAt,
      startedAt: startedAtIso,
      body: null,
      bodyPreview: "",
      error: error instanceof Error ? error.message : String(error),
    };
  } finally {
    clearTimeout(timeout);
  }
}

export async function login(baseUrl = qaApiBaseUrl) {
  const response = await timedJsonRequest({
    label: "auth_login",
    method: "POST",
    url: `${baseUrl}/auth/login`,
    body: {
      phone_number: qaApiLoginPhoneNumber,
      password: qaApiLoginPassword,
    },
    expectedStatus: 200,
  });

  return {
    ...response,
    token: response.body?.access_token || "",
  };
}

export async function ensureAddress(token, baseUrl = qaApiBaseUrl) {
  const listResponse = await timedJsonRequest({
    label: "get_address",
    method: "GET",
    url: `${baseUrl}/user/address`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    expectedStatus: 200,
  });

  if (!listResponse.ok) {
    return {
      ok: false,
      addressId: null,
      steps: [listResponse],
    };
  }

  const currentAddress = listResponse.body?.address_list?.[0];
  if (currentAddress?.id) {
    return {
      ok: true,
      addressId: currentAddress.id,
      steps: [listResponse],
    };
  }

  const createResponse = await timedJsonRequest({
    label: "set_address",
    method: "POST",
    url: `${baseUrl}/user/address`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: {
      address: {
        street: "Assignment 3 Test Street 10",
        description: "Experimental engineering checkout address",
      },
    },
    expectedStatus: 200,
  });

  if (!createResponse.ok) {
    return {
      ok: false,
      addressId: null,
      steps: [listResponse, createResponse],
    };
  }

  const refreshed = await timedJsonRequest({
    label: "get_address_after_set",
    method: "GET",
    url: `${baseUrl}/user/address`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    expectedStatus: 200,
  });

  return {
    ok: refreshed.ok && Boolean(refreshed.body?.address_list?.[0]?.id),
    addressId: refreshed.body?.address_list?.[0]?.id ?? null,
    steps: [listResponse, createResponse, refreshed],
  };
}

export async function clearCart(token, baseUrl = qaApiBaseUrl) {
  return timedJsonRequest({
    label: "clear_cart",
    method: "DELETE",
    url: `${baseUrl}/cart/clear`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    expectedStatus: 200,
  });
}

export async function addCartItem(token, baseUrl = qaApiBaseUrl, quantity = 1) {
  return timedJsonRequest({
    label: "cart_add",
    method: "POST",
    url: `${baseUrl}/cart/add`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: {
      product_id: 1,
      supplier_id: 1,
      quantity,
    },
    expectedStatus: 200,
  });
}

export async function getCart(token, baseUrl = qaApiBaseUrl) {
  return timedJsonRequest({
    label: "cart_get",
    method: "GET",
    url: `${baseUrl}/cart/`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    expectedStatus: 200,
  });
}

export async function checkout(token, addressId, baseUrl = qaApiBaseUrl) {
  return timedJsonRequest({
    label: "cart_checkout",
    method: "POST",
    url: `${baseUrl}/cart/checkout`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: {
      address_id: addressId,
    },
    expectedStatus: 200,
  });
}

export async function catalogProbe(baseUrl = qaApiBaseUrl) {
  return timedJsonRequest({
    label: "catalog_browse",
    method: "GET",
    url: `${baseUrl}/product/list?limit=12&offset=0`,
    expectedStatus: 200,
  });
}

export function runCommand(command, args, options = {}) {
  const result = spawnSync(command, args, {
    cwd: options.cwd ?? rootDir,
    encoding: "utf8",
    env: {
      ...process.env,
      ...options.env,
    },
  });

  return {
    command,
    args,
    exitCode: result.status ?? 1,
    stdout: result.stdout || "",
    stderr: result.stderr || "",
  };
}

export function dockerCompose(args, options = {}) {
  return runCommand("docker", ["compose", "-f", "backend/docker-compose.yaml", ...args], options);
}

export async function waitForHealthy({
  url = `${qaApiHealthUrl}/metrics`,
  timeoutMs = 120000,
  intervalMs = 2000,
}) {
  const startedAt = Date.now();
  while (Date.now() - startedAt < timeoutMs) {
    const probe = await timedJsonRequest({
      label: "health_probe",
      method: "GET",
      url,
      timeoutMs: intervalMs,
    });

    if (probe.status && probe.status < 500) {
      return true;
    }

    await sleep(intervalMs);
  }

  return false;
}

export async function fetchMetricsText() {
  try {
    const response = await fetch(`${qaApiHealthUrl}/metrics`);
    if (!response.ok) {
      return "";
    }
    return await response.text();
  } catch (error) {
    return "";
  }
}

export function parseDockerComposePsJson(rawOutput) {
  const trimmed = rawOutput.trim();
  if (!trimmed) {
    return [];
  }

  try {
    const parsed = JSON.parse(trimmed);
    if (Array.isArray(parsed)) {
      return parsed;
    }
  } catch (error) {
    // Fall back to line-delimited objects.
  }

  return trimmed
    .split("\n")
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      try {
        return JSON.parse(line);
      } catch (error) {
        return null;
      }
    })
    .filter(Boolean);
}

export function parseDockerStats(rawOutput) {
  return rawOutput
    .split("\n")
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      try {
        return JSON.parse(line);
      } catch (error) {
        return null;
      }
    })
    .filter(Boolean);
}

export function parsePercentString(value) {
  return Number(String(value || "0").replace("%", "").trim()) || 0;
}

export function parseMemoryUsageMiB(value) {
  const usage = String(value || "").split("/")[0].trim();
  const match = usage.match(/^([\d.]+)\s*([A-Za-z]+)$/);
  if (!match) {
    return 0;
  }

  const amount = Number(match[1]);
  const unit = match[2].toLowerCase();
  const multipliers = {
    b: 1 / (1024 * 1024),
    kib: 1 / 1024,
    mib: 1,
    gib: 1024,
  };

  return amount * (multipliers[unit] ?? 1);
}

export function createDelayProxy({ listenPort, targetOrigin, delayMs }) {
  const server = http.createServer(async (req, res) => {
    const bodyChunks = [];
    for await (const chunk of req) {
      bodyChunks.push(chunk);
    }

    const bodyBuffer = Buffer.concat(bodyChunks);
    await sleep(delayMs);

    try {
      const response = await fetch(`${targetOrigin}${req.url}`, {
        method: req.method,
        headers: req.headers,
        body: bodyBuffer.length ? bodyBuffer : undefined,
      });

      res.writeHead(response.status, Object.fromEntries(response.headers.entries()));
      const arrayBuffer = await response.arrayBuffer();
      res.end(Buffer.from(arrayBuffer));
    } catch (error) {
      res.writeHead(502, { "content-type": "application/json" });
      res.end(JSON.stringify({ error: error instanceof Error ? error.message : String(error) }));
    }
  });

  return {
    async start() {
      await new Promise((resolve) => server.listen(listenPort, "127.0.0.1", resolve));
    },
    async stop() {
      await new Promise((resolve) => server.close(resolve));
    },
  };
}
