import path from "node:path";
import {
  addCartItem,
  appendNdjson,
  average,
  catalogProbe,
  chartsDir,
  clearCart,
  dockerCompose,
  ensureDir,
  fetchMetricsText,
  getCart,
  login,
  logsDir,
  median,
  nowIso,
  parseDockerComposePsJson,
  parseDockerStats,
  parseMemoryUsageMiB,
  parsePercentString,
  percentile,
  qaApiBaseUrl,
  readJson,
  rootDir,
  round,
  runCommand,
  sleep,
  statusBucket,
  writeText,
  writeJson,
} from "./experimental-common.mjs";

const rawLogPath = path.join(logsDir, "performance-raw.ndjson");
const summaryPath = path.join(logsDir, "performance-summary.json");
const resourcePath = path.join(logsDir, "performance-resources.json");
const configPath = path.join(rootDir, "qa", "performance", "scenarios.json");

ensureDir(logsDir);
ensureDir(chartsDir);
writeText(rawLogPath, "");

const config = readJson(configPath);

async function ensureBackendReady() {
  const metrics = await fetchMetricsText();
  if (!metrics) {
    throw new Error(
      "Performance tests require the live backend stack. Start it with `npm run qa:env:up` and wait for /metrics.",
    );
  }
}

async function runFlow(flowId) {
  if (flowId === "auth_login") {
    return [await login()];
  }

  if (flowId === "catalog_browse") {
    return [await catalogProbe()];
  }

  if (flowId === "cart_session") {
    const auth = await login();
    if (!auth.ok || !auth.token) {
      return [auth];
    }

    const steps = [auth];
    steps.push(await clearCart(auth.token));
    steps.push(await addCartItem(auth.token));
    steps.push(await getCart(auth.token));
    steps.push(await clearCart(auth.token));
    return steps;
  }

  throw new Error(`Unsupported flow: ${flowId}`);
}

function pickFlow(flows) {
  const random = Math.random();
  let cursor = 0;
  for (const flow of flows) {
    cursor += flow.weight;
    if (random <= cursor) {
      return flow;
    }
  }
  return flows.at(-1);
}

function activeUsersForScenario(scenario, elapsedMs) {
  let elapsedSeconds = elapsedMs / 1000;
  for (const stage of scenario.stages) {
    if (elapsedSeconds <= stage.seconds) {
      return stage.users;
    }
    elapsedSeconds -= stage.seconds;
  }
  return scenario.stages.at(-1)?.users ?? 0;
}

function summarizeScenario(scenario, requestLogs) {
  const durations = requestLogs.map((item) => item.durationMs);
  const errors = requestLogs.filter((item) => !item.ok).length;
  const totalRequests = requestLogs.length;
  const throughputRps = totalRequests / Math.max(1, scenario.durationSeconds);
  const statusSummary = {};
  for (const log of requestLogs) {
    const bucket = statusBucket(log.status);
    statusSummary[bucket] = (statusSummary[bucket] || 0) + 1;
  }

  const metrics = {
    totalRequests,
    averageResponseTimeMs: round(average(durations)),
    medianResponseTimeMs: round(median(durations)),
    p95ResponseTimeMs: round(percentile(durations, 95)),
    throughputRps: round(throughputRps),
    errorRatePercent: round((errors / Math.max(1, totalRequests)) * 100),
    errors,
    statusSummary,
  };

  const thresholds = scenario.thresholds;
  const thresholdEvaluation = {
    averageResponseTimeMs: metrics.averageResponseTimeMs <= thresholds.averageResponseTimeMs,
    medianResponseTimeMs: metrics.medianResponseTimeMs <= thresholds.medianResponseTimeMs,
    p95ResponseTimeMs: metrics.p95ResponseTimeMs <= thresholds.p95ResponseTimeMs,
    throughputRps: metrics.throughputRps >= thresholds.throughputRps,
    errorRatePercent: metrics.errorRatePercent <= thresholds.errorRatePercent,
  };

  const passed = Object.values(thresholdEvaluation).every(Boolean);

  return {
    id: scenario.id,
    type: scenario.type,
    description: scenario.description,
    durationSeconds: scenario.durationSeconds,
    concurrencyProfile: scenario.stages,
    thresholds,
    thresholdEvaluation,
    passed,
    metrics,
  };
}

function detectBottlenecks(results) {
  const findings = [];
  for (const scenario of results) {
    if (!scenario.passed && !scenario.thresholdEvaluation.p95ResponseTimeMs) {
      findings.push(
        `${scenario.id}: tail latency exceeded threshold (${scenario.metrics.p95ResponseTimeMs}ms p95 vs ${scenario.thresholds.p95ResponseTimeMs}ms).`,
      );
    }
    if (!scenario.passed && !scenario.thresholdEvaluation.errorRatePercent) {
      findings.push(
        `${scenario.id}: error rate reached ${scenario.metrics.errorRatePercent}% with status mix ${JSON.stringify(scenario.metrics.statusSummary)}.`,
      );
    }
    if (!scenario.thresholdEvaluation.throughputRps) {
      findings.push(
        `${scenario.id}: achieved ${scenario.metrics.throughputRps} rps, below the ${scenario.thresholds.throughputRps} rps target.`,
      );
    }
  }

  return findings.length ? findings : ["No threshold breaches were observed in the executed scenarios."];
}

function collectResourceSample() {
  const psResult = dockerCompose(["ps", "--format", "json"]);
  if (psResult.exitCode !== 0) {
    return {
      collectedAt: nowIso(),
      supported: false,
      error: psResult.stderr.trim() || psResult.stdout.trim() || "docker compose ps failed",
      samples: [],
    };
  }

  const containers = parseDockerComposePsJson(psResult.stdout)
    .map((item) => item.Name || item.name)
    .filter(Boolean);

  if (!containers.length) {
    return {
      collectedAt: nowIso(),
      supported: false,
      error: "No running containers found for backend compose project.",
      samples: [],
    };
  }

  const statsResult = runCommand("docker", ["stats", "--no-stream", "--format", "{{json .}}", ...containers]);
  if (statsResult.exitCode !== 0) {
    return {
      collectedAt: nowIso(),
      supported: false,
      error: statsResult.stderr.trim() || statsResult.stdout.trim() || "docker stats failed",
      samples: [],
    };
  }

  const parsed = parseDockerStats(statsResult.stdout).map((item) => ({
    container: item.Name,
    cpuPercent: round(parsePercentString(item.CPUPerc)),
    memoryUsageMiB: round(parseMemoryUsageMiB(item.MemUsage)),
    blockIo: item.BlockIO,
  }));

  return {
    collectedAt: nowIso(),
    supported: true,
    error: null,
    samples: parsed,
  };
}

async function main() {
  await ensureBackendReady();

  const resourceSamples = [];
  const scenarioResults = [];
  const metricsSnapshots = [];

  for (const scenario of config.scenarios) {
    const scenarioStart = Date.now();
    const scenarioRequestLogs = [];
    const maxUsers = Math.max(...scenario.stages.map((stage) => stage.users));
    const beforeMetrics = await fetchMetricsText();
    metricsSnapshots.push({
      scenarioId: scenario.id,
      phase: "before",
      collectedAt: nowIso(),
      metricsLength: beforeMetrics.length,
    });
    resourceSamples.push({
      scenarioId: scenario.id,
      phase: "before",
      ...collectResourceSample(),
    });

    const workers = Array.from({ length: maxUsers }, (_, workerIndex) =>
      (async () => {
        while (Date.now() - scenarioStart < scenario.durationSeconds * 1000) {
          const elapsedMs = Date.now() - scenarioStart;
          const activeUsers = activeUsersForScenario(scenario, elapsedMs);
          if (workerIndex >= activeUsers) {
            await sleep(150);
            continue;
          }

          const selectedFlow = pickFlow(config.flows);
          const steps = await runFlow(selectedFlow.id);
          for (const step of steps) {
            const record = {
              scenarioId: scenario.id,
              scenarioType: scenario.type,
              workerIndex,
              flowId: selectedFlow.id,
              module: selectedFlow.module,
              ...step,
            };
            scenarioRequestLogs.push(record);
            appendNdjson(rawLogPath, record);
          }
        }
      })(),
    );

    await Promise.all(workers);
    resourceSamples.push({
      scenarioId: scenario.id,
      phase: "after",
      ...collectResourceSample(),
    });
    const afterMetrics = await fetchMetricsText();
    metricsSnapshots.push({
      scenarioId: scenario.id,
      phase: "after",
      collectedAt: nowIso(),
      metricsLength: afterMetrics.length,
    });
    scenarioResults.push(summarizeScenario(scenario, scenarioRequestLogs));
  }

  writeJson(resourcePath, resourceSamples);

  const overallDurations = scenarioResults.flatMap((result) =>
    result.metrics.totalRequests ? [result.metrics.averageResponseTimeMs] : [],
  );
  const summary = {
    generatedAt: nowIso(),
    selectedHighRiskModules: config.selectedHighRiskModules,
    baseUrl: qaApiBaseUrl,
    scenarios: scenarioResults,
    totals: {
      executedScenarios: scenarioResults.length,
      allThresholdsPassed: scenarioResults.every((item) => item.passed),
      averageOfScenarioAveragesMs: round(average(overallDurations)),
    },
    metricsSnapshots,
    bottlenecks: detectBottlenecks(scenarioResults),
  };

  writeJson(summaryPath, summary);

  console.log(`Performance summary written to ${summaryPath}`);
}

main().catch((error) => {
  console.error(error instanceof Error ? error.stack || error.message : String(error));
  process.exit(1);
});
