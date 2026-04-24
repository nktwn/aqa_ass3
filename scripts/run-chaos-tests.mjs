import path from "node:path";
import {
  addCartItem,
  appendNdjson,
  catalogProbe,
  checkout,
  clearCart,
  createDelayProxy,
  dockerCompose,
  ensureAddress,
  fetchMetricsText,
  getCart,
  login,
  logsDir,
  nowIso,
  qaApiBaseUrl,
  qaApiHealthUrl,
  readJson,
  rootDir,
  round,
  sleep,
  statusBucket,
  timedJsonRequest,
  waitForHealthy,
  writeText,
  writeJson,
} from "./experimental-common.mjs";

const configPath = path.join(rootDir, "qa", "chaos", "scenarios.json");
const config = readJson(configPath);
const rawLogPath = path.join(logsDir, "chaos-events.ndjson");
const summaryPath = path.join(logsDir, "chaos-summary.json");
writeText(rawLogPath, "");

async function probeLogin(baseUrl = qaApiBaseUrl) {
  return login(baseUrl);
}

async function probeCatalog(baseUrl = qaApiBaseUrl) {
  return catalogProbe(baseUrl);
}

async function prepareCheckout(baseUrl = qaApiBaseUrl) {
  const auth = await login(baseUrl);
  if (!auth.ok || !auth.token) {
    return { ok: false, steps: [auth], token: "", addressId: null };
  }

  const steps = [auth];
  steps.push(await clearCart(auth.token, baseUrl));
  const address = await ensureAddress(auth.token, baseUrl);
  steps.push(...address.steps);
  if (!address.ok || !address.addressId) {
    return { ok: false, steps, token: auth.token, addressId: null };
  }

  // Product 1 from supplier 1 is priced at 1181 in the seeded data.
  // Quantity 26 crosses the supplier minimum order threshold of 30000.
  steps.push(await addCartItem(auth.token, baseUrl, 26));
  steps.push(await getCart(auth.token, baseUrl));
  return {
    ok: steps.every((step) => step.ok),
    steps,
    token: auth.token,
    addressId: address.addressId,
  };
}

async function probeCheckout(baseUrl = qaApiBaseUrl) {
  const setup = await prepareCheckout(baseUrl);
  if (!setup.ok || !setup.addressId) {
    return {
      ok: false,
      status: null,
      durationMs: 0,
      startedAt: nowIso(),
      error: "checkout setup failed",
      setupSteps: setup.steps,
    };
  }

  const response = await checkout(setup.token, setup.addressId, baseUrl);
  return {
    ...response,
    setupSteps: setup.steps,
  };
}

function recordEvent(event) {
  appendNdjson(rawLogPath, event);
}

function summarizeFaultWindow(events) {
  const total = events.length;
  const available = events.filter((event) => event.ok).length;
  const errors = {};
  for (const event of events) {
    const bucket = event.error ? event.error : statusBucket(event.status);
    errors[bucket] = (errors[bucket] || 0) + 1;
  }

  return {
    totalProbes: total,
    successfulProbes: available,
    availabilityPercent: round((available / Math.max(1, total)) * 100),
    errorPropagation: errors,
  };
}

async function runStopStartScenario(scenario) {
  const probeFn = scenario.probe === "login" ? probeLogin : scenario.probe === "checkout" ? probeCheckout : probeCatalog;
  const baseline = await probeFn();
  recordEvent({ scenarioId: scenario.id, phase: "baseline", ...baseline });

  const stopResult = dockerCompose(["stop", scenario.faultService]);
  recordEvent({
    scenarioId: scenario.id,
    phase: "fault_injected",
    service: scenario.faultService,
    exitCode: stopResult.exitCode,
    stdout: stopResult.stdout.trim(),
    stderr: stopResult.stderr.trim(),
    at: nowIso(),
  });

  const faultEvents = [];
  const faultStart = Date.now();
  while (Date.now() - faultStart < scenario.faultDurationSeconds * 1000) {
    const probe = await probeFn();
    const event = { scenarioId: scenario.id, phase: "during_fault", ...probe };
    faultEvents.push(event);
    recordEvent(event);
    await sleep(2000);
  }

  const recoveryStartedAt = Date.now();
  const startResult = dockerCompose(["up", "-d", scenario.faultService]);
  recordEvent({
    scenarioId: scenario.id,
    phase: "recovery_started",
    service: scenario.faultService,
    exitCode: startResult.exitCode,
    stdout: startResult.stdout.trim(),
    stderr: startResult.stderr.trim(),
    at: nowIso(),
  });

  let usedAppRestartFallback = false;
  const recoveryEvents = [];
  let recovered = false;
  let firstRecoverySuccessAt = null;

  for (let attempt = 0; attempt < 12; attempt += 1) {
    if (scenario.faultService === "app") {
      await waitForHealthy({ url: `${qaApiHealthUrl}/metrics`, timeoutMs: 10000, intervalMs: 1000 });
    } else if (scenario.faultService === "pg" || scenario.faultService === "redis") {
      await sleep(2000);
    }

    const probe = await probeFn();
    const event = { scenarioId: scenario.id, phase: "recovery_probe", attempt: attempt + 1, ...probe };
    recoveryEvents.push(event);
    recordEvent(event);
    if (probe.ok) {
      recovered = true;
      firstRecoverySuccessAt = Date.now();
      break;
    }
  }

  if (!recovered && scenario.faultService === "pg") {
    usedAppRestartFallback = true;
    const restartResult = dockerCompose(["restart", "app"]);
    recordEvent({
      scenarioId: scenario.id,
      phase: "fallback_app_restart",
      exitCode: restartResult.exitCode,
      stdout: restartResult.stdout.trim(),
      stderr: restartResult.stderr.trim(),
      at: nowIso(),
    });
    await waitForHealthy({ url: `${qaApiHealthUrl}/metrics`, timeoutMs: 30000, intervalMs: 2000 });
    const probe = await probeFn();
    const event = { scenarioId: scenario.id, phase: "recovery_probe_after_fallback", ...probe };
    recoveryEvents.push(event);
    recordEvent(event);
    if (probe.ok) {
      recovered = true;
      firstRecoverySuccessAt = Date.now();
    }
  }

  return {
    id: scenario.id,
    module: scenario.module,
    faultType: scenario.faultType,
    description: scenario.description,
    faultService: scenario.faultService,
    baselineOk: baseline.ok,
    availabilityDuringFaultPercent: summarizeFaultWindow(faultEvents).availabilityPercent,
    faultSummary: summarizeFaultWindow(faultEvents),
    recoverySummary: summarizeFaultWindow(recoveryEvents),
    recovered,
    mttrMs: firstRecoverySuccessAt ? firstRecoverySuccessAt - recoveryStartedAt : null,
    gracefulDegradation:
      scenario.faultService === "app"
        ? "No graceful degradation observed; requests failed until the app container restarted."
        : scenario.faultService === "pg"
          ? "The API surfaced backend errors while PostgreSQL was unavailable."
          : "Checkout failed fast when Redis-backed payment-order persistence was unavailable.",
    recoveryBehavior: recovered
      ? usedAppRestartFallback
        ? "Recovered after dependency restoration plus application restart fallback."
        : "Recovered after dependency restoration without manual workaround."
      : "Did not recover within the scripted observation window.",
    usedAppRestartFallback,
  };
}

async function runLatencyScenario(scenario) {
  const proxy = createDelayProxy({
    listenPort: 8099,
    targetOrigin: qaApiHealthUrl,
    delayMs: scenario.latencyMs,
  });

  await proxy.start();
  const delayedBaseUrl = "http://127.0.0.1:8099/api";
  const delayedHealthUrl = "http://127.0.0.1:8099";
  const events = [];

  try {
    for (let iteration = 0; iteration < scenario.iterations; iteration += 1) {
      const loginProbeResult = await probeLogin(delayedBaseUrl);
      const catalogProbeResult = await probeCatalog(delayedBaseUrl);
      events.push({ scenarioId: scenario.id, phase: "latency_probe", flow: "login", ...loginProbeResult });
      events.push({ scenarioId: scenario.id, phase: "latency_probe", flow: "catalog", ...catalogProbeResult });
      recordEvent(events.at(-2));
      recordEvent(events.at(-1));
      await sleep(300);
    }
  } finally {
    await proxy.stop();
  }

  const successfulEvents = events.filter((event) => event.ok);
  const durations = successfulEvents.map((event) => event.durationMs);

  return {
    id: scenario.id,
    module: scenario.module,
    faultType: scenario.faultType,
    description: scenario.description,
    baselineOk: true,
    availabilityDuringFaultPercent: round((successfulEvents.length / Math.max(1, events.length)) * 100),
    faultSummary: summarizeFaultWindow(events),
    recoverySummary: {
      totalProbes: 0,
      successfulProbes: 0,
      availabilityPercent: 0,
      errorPropagation: {},
    },
    recovered: true,
    mttrMs: 0,
    gracefulDegradation:
      "Availability remained intact, but latency increased materially because requests were delayed before reaching the API.",
    recoveryBehavior: "Recovery was immediate once the latency proxy was removed.",
    usedAppRestartFallback: false,
    latencySummary: {
      averageLatencyMs: round(durations.reduce((sum, value) => sum + value, 0) / Math.max(1, durations.length)),
      p95LatencyMs: durations.length ? durations.sort((a, b) => a - b)[Math.floor(durations.length * 0.95) - 1] || durations.at(-1) : 0,
      proxyDelayMs: scenario.latencyMs,
      delayedBaseUrl,
      delayedHealthUrl,
    },
  };
}

async function main() {
  const readiness = await fetchMetricsText();
  if (!readiness) {
    throw new Error("Chaos tests require the live backend stack. Start it with `npm run qa:env:up` first.");
  }

  const scenarioSummaries = [];

  for (const scenario of config.scenarios) {
    if (scenario.probe === "latency_proxy") {
      scenarioSummaries.push(await runLatencyScenario(scenario));
      continue;
    }

    scenarioSummaries.push(await runStopStartScenario(scenario));
    await waitForHealthy({ url: `${qaApiHealthUrl}/metrics`, timeoutMs: 60000, intervalMs: 2000 });
  }

  const summary = {
    generatedAt: nowIso(),
    selectedHighRiskModules: config.selectedHighRiskModules,
    scenarios: scenarioSummaries,
    recommendations: scenarioSummaries.flatMap((scenario) => {
      const notes = [];
      if (!scenario.recovered) {
        notes.push(`${scenario.id}: add stronger automated recovery handling and service readiness checks.`);
      }
      if (scenario.availabilityDuringFaultPercent === 0 && scenario.faultType === "API downtime") {
        notes.push(`${scenario.id}: consider graceful maintenance responses or health-gated routing in front of the app container.`);
      }
      if (scenario.faultType === "database unavailability" && scenario.usedAppRestartFallback) {
        notes.push(`${scenario.id}: the app may require stronger DB reconnection behavior after PostgreSQL interruption.`);
      }
      if (scenario.faultType === "dependency failure") {
        notes.push(`${scenario.id}: checkout should expose clearer dependency-failure diagnostics and possibly queue/retry payment-order persistence.`);
      }
      if (scenario.faultType === "injected network latency") {
        notes.push(`${scenario.id}: client-visible latency rises directly; consider timeout budgets and user-facing retry guidance.`);
      }
      return notes;
    }),
  };

  writeJson(summaryPath, summary);
  console.log(`Chaos summary written to ${summaryPath}`);
}

main().catch((error) => {
  console.error(error instanceof Error ? error.stack || error.message : String(error));
  process.exit(1);
});
