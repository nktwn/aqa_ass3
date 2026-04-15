import fs from "node:fs";
import path from "node:path";

const rootDir = process.cwd();
const logsDir = path.join(rootDir, "logs");
const qaDocsDir = path.join(rootDir, "qa-docs");
const tablesDir = path.join(qaDocsDir, "tables");
const chartsDir = path.join(rootDir, "evidence", "charts");

const HIGH_RISK_MODULES = [
  "Authentication and authorization",
  "Cart management",
  "Checkout and payment initiation",
  "Order lifecycle",
  "Product catalog and listing",
];

const ORIGINAL_RISK_SCORES = {
  "Authentication and authorization": 20,
  "Cart management": 20,
  "Checkout and payment initiation": 20,
  "Order lifecycle": 20,
  "Product catalog and listing": 16,
};

const reportDefinitions = [
  { key: "smoke", label: "Smoke", file: path.join(logsDir, "smoke-report.json") },
  { key: "critical", label: "Critical", file: path.join(logsDir, "critical-report.json") },
  { key: "regression", label: "Regression", file: path.join(logsDir, "regression-report.json") },
  { key: "e2e", label: "E2E", file: path.join(logsDir, "e2e-report.json") },
  { key: "api", label: "API", file: path.join(logsDir, "api-report.json") },
];

fs.mkdirSync(logsDir, { recursive: true });
fs.mkdirSync(tablesDir, { recursive: true });
fs.mkdirSync(chartsDir, { recursive: true });

function readJsonIfExists(filePath, fallback) {
  if (!fs.existsSync(filePath)) {
    return fallback;
  }

  try {
    return JSON.parse(fs.readFileSync(filePath, "utf8"));
  } catch (error) {
    return fallback;
  }
}

function toCsv(rows) {
  return `${rows
    .map((row) =>
      row.map((value) => `"${String(value ?? "").replaceAll('"', '""')}"`).join(","),
    )
    .join("\n")}\n`;
}

function round(value, digits = 2) {
  return Number(value.toFixed(digits));
}

function percent(numerator, denominator) {
  if (!denominator) {
    return 0;
  }
  return round((numerator / denominator) * 100);
}

function formatMs(durationMs) {
  return `${round(durationMs / 1000, 2)}s`;
}

function parseAttachmentContent(attachment) {
  if (!attachment?.path || !fs.existsSync(attachment.path)) {
    return null;
  }

  if (attachment.contentType !== "application/json") {
    return null;
  }

  try {
    return JSON.parse(fs.readFileSync(attachment.path, "utf8"));
  } catch (error) {
    return null;
  }
}

function buildModuleFromHint({ title = "", file = "", tags = [], annotations = [] }) {
  const metadataModule = annotations.find((annotation) => annotation.type === "qa.module")?.description;
  if (metadataModule) {
    return metadataModule;
  }

  const normalizedTitle = title.toLowerCase();
  const normalizedFile = file.toLowerCase();
  const normalizedTags = tags.map((tag) => String(tag).toLowerCase());

  if (
    normalizedTags.includes("auth") ||
    /\/auth\//.test(normalizedFile) ||
    /(authentication|authorization|jwt|login)/.test(normalizedTitle)
  ) {
    return "Authentication and authorization";
  }

  if (
    normalizedTags.includes("cart") ||
    /\/cart\//.test(normalizedFile) ||
    /\bcart\b/.test(normalizedTitle)
  ) {
    return "Cart management";
  }

  if (
    normalizedTags.includes("checkout") ||
    /checkout/.test(normalizedFile) ||
    /(checkout|payment)/.test(normalizedTitle)
  ) {
    return "Checkout and payment initiation";
  }

  if (
    normalizedTags.includes("orders") ||
    /\/orders?\//.test(normalizedFile) ||
    /\border\b|\borders\b/.test(normalizedTitle)
  ) {
    return "Order lifecycle";
  }

  if (
    normalizedTags.includes("catalog") ||
    /catalog|product/.test(normalizedFile) ||
    /(catalog|product list|product details)/.test(normalizedTitle)
  ) {
    return "Product catalog and listing";
  }

  return "Cross-module";
}

function isRuntimeDependentTest(testCase) {
  const annotation = testCase.annotations.find((item) => item.type === "qa.runtimeDependent")?.description;
  if (annotation === "true") {
    return true;
  }

  return (
    testCase.file.includes("tests/api/") ||
    testCase.file.includes(".live.") ||
    /backend smoke|live api/i.test(testCase.title)
  );
}

function classifyFailureType(message) {
  const normalized = (message || "").toLowerCase();
  if (!normalized) return "unknown";
  if (normalized.includes("econnrefused")) return "environment-connection-refused";
  if (normalized.includes("timed out") || normalized.includes("timeout")) return "timeout";
  if (normalized.includes("tohaveurl") || normalized.includes("to be visible") || normalized.includes("expect(")) return "assertion-failure";
  if (normalized.includes("invalid cart")) return "business-rule-invalid-cart";
  if (normalized.includes("unauthorized")) return "authorization-failure";
  if (normalized.includes("address")) return "validation-failure";
  return "unexpected-error";
}

function collectReportTests(node, suiteMeta, parentTitle = "") {
  const currentTitle = [parentTitle, node.title].filter(Boolean).join(" > ");
  const testCases = [];

  for (const suite of node.suites ?? []) {
    testCases.push(...collectReportTests(suite, suiteMeta, currentTitle));
  }

  for (const spec of node.specs ?? []) {
    const specTitle = [currentTitle, spec.title].filter(Boolean).join(" > ");
    for (const test of spec.tests ?? []) {
      const results = test.results ?? [];
      const finalResult = results.at(-1) ?? {};
      const annotations = [
        ...(spec.annotations ?? []),
        ...(test.annotations ?? []),
        ...(finalResult.annotations ?? []),
      ];
      const attachmentMetadata = results
        .flatMap((result) => result.attachments ?? [])
        .map((attachment) => parseAttachmentContent(attachment))
        .find(Boolean);

      testCases.push({
        suiteKey: suiteMeta.key,
        suiteLabel: suiteMeta.label,
        reportFile: suiteMeta.file,
        file: `tests/${node.file ?? spec.file ?? ""}`.replace("tests/tests/", "tests/"),
        title: specTitle,
        status: finalResult.status ?? "unknown",
        durationMs: finalResult.duration ?? 0,
        projectName: test.projectName ?? "chromium",
        tags: spec.tags ?? [],
        annotations,
        metadata: attachmentMetadata,
        errors: results.flatMap((result) => result.errors ?? []),
        retryStatuses: results.map((result) => result.status).filter(Boolean),
      });
    }
  }

  return testCases;
}

function loadPlaywrightTests() {
  const suites = [];
  const allTests = [];

  for (const definition of reportDefinitions) {
    if (!fs.existsSync(definition.file)) {
      continue;
    }

    const report = readJsonIfExists(definition.file, null);
    if (!report) {
      continue;
    }

    const suiteTests = (report.suites ?? []).flatMap((suite) => collectReportTests(suite, definition));
    suites.push({
      key: definition.key,
      label: definition.label,
      tests: suiteTests,
    });
    allTests.push(...suiteTests);
  }

  return { suites, allTests };
}

function walkDir(directory, predicate = () => true) {
  if (!fs.existsSync(directory)) {
    return [];
  }

  const entries = fs.readdirSync(directory, { withFileTypes: true });
  return entries.flatMap((entry) => {
    const nextPath = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      return walkDir(nextPath, predicate);
    }
    return predicate(nextPath) ? [nextPath] : [];
  });
}

function importPathToFile(importPath) {
  if (importPath.startsWith("diploma/")) {
    return path.join(rootDir, "backend", importPath.replace(/^diploma\//, ""));
  }
  return path.join(rootDir, importPath);
}

function mapCoverageModule(importPath) {
  if (/^diploma\/modules\/auth\//.test(importPath)) return "Authentication and authorization";
  if (/^diploma\/modules\/order\//.test(importPath)) return "Order lifecycle";
  if (/^diploma\/modules\/product\//.test(importPath)) return "Product catalog and listing";
  if (/^diploma\/modules\/cart\/(service\/checkout\.go|handler\/checkout\.go|client\/payment\/)/.test(importPath)) {
    return "Checkout and payment initiation";
  }
  if (/^diploma\/modules\/cart\//.test(importPath)) return "Cart management";
  return "Other";
}

function extractGoFunctionNames(filePath) {
  if (!fs.existsSync(filePath)) {
    return [];
  }
  const source = fs.readFileSync(filePath, "utf8");
  const matches = [...source.matchAll(/^func\s+(?:\([^)]*\)\s*)?([A-Za-z0-9_]+)/gm)];
  return matches.map((match) => match[1]);
}

function parseGoCoverage() {
  const coverageFile = path.join(logsDir, "go-coverage.out");
  const profile = fs.existsSync(coverageFile) ? fs.readFileSync(coverageFile, "utf8").trim().split("\n").slice(1) : [];
  const fileCoverage = new Map();
  const moduleCoverage = new Map(
    HIGH_RISK_MODULES.map((moduleName) => [moduleName, { coveredStatements: 0, totalStatements: 0 }]),
  );

  for (const line of profile) {
    const [location, statementsRaw, countRaw] = line.split(" ");
    const importPath = location.split(":")[0];
    const statements = Number(statementsRaw);
    const count = Number(countRaw);
    const moduleName = mapCoverageModule(importPath);
    const fileStats = fileCoverage.get(importPath) ?? { coveredStatements: 0, totalStatements: 0 };

    fileStats.totalStatements += statements;
    if (count > 0) {
      fileStats.coveredStatements += statements;
    }
    fileCoverage.set(importPath, fileStats);

    if (moduleCoverage.has(moduleName)) {
      const moduleStats = moduleCoverage.get(moduleName);
      moduleStats.totalStatements += statements;
      if (count > 0) {
        moduleStats.coveredStatements += statements;
      }
    }
  }

  const allRiskFiles = [
    ...walkDir(path.join(rootDir, "backend", "modules", "auth"), (filePath) => filePath.endsWith(".go") && !filePath.endsWith("_test.go")),
    ...walkDir(path.join(rootDir, "backend", "modules", "cart"), (filePath) => filePath.endsWith(".go") && !filePath.endsWith("_test.go")),
    ...walkDir(path.join(rootDir, "backend", "modules", "order"), (filePath) => filePath.endsWith(".go") && !filePath.endsWith("_test.go")),
    ...walkDir(path.join(rootDir, "backend", "modules", "product"), (filePath) => filePath.endsWith(".go") && !filePath.endsWith("_test.go")),
  ];

  const gapEntries = allRiskFiles.map((absoluteFilePath) => {
    const importPath = `diploma/${path.relative(path.join(rootDir, "backend"), absoluteFilePath).replaceAll(path.sep, "/")}`;
    const stats = fileCoverage.get(importPath) ?? { coveredStatements: 0, totalStatements: 0 };
    const uncoveredFunctions = stats.coveredStatements === 0 ? extractGoFunctionNames(absoluteFilePath) : [];
    return {
      module: mapCoverageModule(importPath),
      file: path.relative(rootDir, absoluteFilePath),
      coveragePercent: percent(stats.coveredStatements, stats.totalStatements),
      uncoveredFunctions,
    };
  }).filter((entry) => HIGH_RISK_MODULES.includes(entry.module));

  const moduleSummary = HIGH_RISK_MODULES.map((moduleName) => {
    const stats = moduleCoverage.get(moduleName) ?? { coveredStatements: 0, totalStatements: 0 };
    return {
      module: moduleName,
      coveredStatements: stats.coveredStatements,
      totalStatements: stats.totalStatements,
      coveragePercent: percent(stats.coveredStatements, stats.totalStatements),
      belowThreshold: percent(stats.coveredStatements, stats.totalStatements) < 70,
    };
  });

  return {
    moduleSummary,
    gapEntries,
    averageCoveragePercent: round(
      moduleSummary.reduce((sum, entry) => sum + entry.coveragePercent, 0) / (moduleSummary.length || 1),
    ),
  };
}

function parseGoTestSummary() {
  const filePath = path.join(logsDir, "go-test.json");
  const lines = fs.existsSync(filePath)
    ? fs.readFileSync(filePath, "utf8").split("\n").filter(Boolean)
    : [];

  const packages = new Map();
  const tests = new Map();

  for (const line of lines) {
    const record = JSON.parse(line);
    if (record.Package) {
      const pkg = packages.get(record.Package) ?? { package: record.Package, status: "unknown", elapsed: 0 };
      if (record.Action === "pass" || record.Action === "fail" || record.Action === "skip") {
        pkg.status = record.Action;
        pkg.elapsed = record.Elapsed ?? pkg.elapsed;
      }
      packages.set(record.Package, pkg);
    }

    if (record.Test) {
      const key = `${record.Package}::${record.Test}`;
      const testRecord = tests.get(key) ?? { package: record.Package, test: record.Test, status: "unknown", elapsed: 0 };
      if (record.Action === "pass" || record.Action === "fail" || record.Action === "skip") {
        testRecord.status = record.Action;
        testRecord.elapsed = record.Elapsed ?? testRecord.elapsed;
      }
      tests.set(key, testRecord);
    }
  }

  const packageList = [...packages.values()];
  const testList = [...tests.values()];
  return {
    packages: packageList,
    tests: testList,
    packagePassRate: percent(packageList.filter((item) => item.status === "pass").length, packageList.length),
    testPassRate: percent(testList.filter((item) => item.status === "pass").length, testList.length),
    failedTests: testList.filter((item) => item.status === "fail"),
    totalElapsedSeconds: round(packageList.reduce((sum, item) => sum + (item.elapsed ?? 0), 0)),
  };
}

function summarizePlaywright({ suites, allTests }) {
  const normalizedTests = allTests.map((testCase) => {
    const moduleName = testCase.metadata?.module ?? buildModuleFromHint(testCase);
    const scenarioType = testCase.metadata?.scenarioType ?? testCase.annotations.find((item) => item.type === "qa.scenario")?.description ?? "baseline";
    const runtimeDependent = testCase.metadata?.runtimeDependent ?? isRuntimeDependentTest(testCase);
    const message = testCase.errors.map((error) => error.message).join("\n");

    return {
      ...testCase,
      module: moduleName,
      scenarioType,
      runtimeDependent,
      testId:
        testCase.metadata?.testId ??
        testCase.annotations.find((item) => item.type === "qa.test.id")?.description ??
        "unmapped",
      inputData:
        testCase.metadata?.inputData ??
        testCase.annotations.find((item) => item.type === "qa.input")?.description ??
        "n/a",
      expectedOutput:
        testCase.metadata?.expectedOutput ??
        testCase.annotations.find((item) => item.type === "qa.expected")?.description ??
        "n/a",
      failureType: classifyFailureType(message),
      failureMessage: message,
    };
  });

  const failureMap = new Map();
  for (const testCase of normalizedTests.filter((item) => item.status === "failed" || item.status === "timedOut")) {
    const key = `${testCase.title}::${testCase.failureType}`;
    const existing = failureMap.get(key) ?? {
      testId: testCase.testId,
      testName: testCase.title,
      module: testCase.module,
      failureType: testCase.failureType,
      failureMessage: testCase.failureMessage.split("\n")[0],
      frequencyAcrossRuns: 0,
      runtimeDependent: testCase.runtimeDependent,
    };
    existing.frequencyAcrossRuns += 1;
    failureMap.set(key, existing);
  }

  const flakyMap = new Map();
  for (const testCase of normalizedTests) {
    const key = testCase.testId !== "unmapped" ? testCase.testId : testCase.title;
    const existing = flakyMap.get(key) ?? {
      testId: testCase.testId,
      testName: testCase.title,
      module: testCase.module,
      passes: 0,
      failures: 0,
      statuses: new Set(),
      suspectedCause: "",
    };

    for (const status of testCase.retryStatuses) {
      existing.statuses.add(status);
      if (status === "passed") existing.passes += 1;
      if (status === "failed" || status === "timedOut") existing.failures += 1;
    }

    if (!testCase.retryStatuses.length) {
      existing.statuses.add(testCase.status);
      if (testCase.status === "passed") existing.passes += 1;
      if (testCase.status === "failed" || testCase.status === "timedOut") existing.failures += 1;
    }

    if (testCase.runtimeDependent && !existing.suspectedCause) {
      existing.suspectedCause = "Runtime-dependent backend availability or environment setup mismatch";
    } else if (testCase.failureType === "timeout") {
      existing.suspectedCause = "Slow response or synchronization issue";
    }

    flakyMap.set(key, existing);
  }

  const flakyTests = [...flakyMap.values()]
    .filter((item) => item.statuses.has("passed") && (item.statuses.has("failed") || item.statuses.has("timedOut")))
    .map((item) => ({
      testId: item.testId,
      testName: item.testName,
      module: item.module,
      passes: item.passes,
      failures: item.failures,
      suspectedCause: item.suspectedCause || "Inconsistent result across runs",
    }));

  const suiteSummary = suites.map((suite) => {
    const executableTests = suite.tests.filter((item) => item.status !== "skipped");
    return {
      key: suite.key,
      label: suite.label,
      total: suite.tests.length,
      passed: suite.tests.filter((item) => item.status === "passed").length,
      failed: suite.tests.filter((item) => item.status === "failed" || item.status === "timedOut").length,
      skipped: suite.tests.filter((item) => item.status === "skipped").length,
      durationMs: suite.tests.reduce((sum, item) => sum + item.durationMs, 0),
      passRate: percent(
        suite.tests.filter((item) => item.status === "passed").length,
        executableTests.length,
      ),
    };
  });

  const performanceAnomalies = normalizedTests
    .filter((item) => item.durationMs >= 5000)
    .map((item) => ({
      testId: item.testId,
      testName: item.title,
      module: item.module,
      durationMs: item.durationMs,
      note: "Execution duration exceeded 5 seconds",
    }));

  return {
    tests: normalizedTests,
    failures: [...failureMap.values()],
    flakyTests,
    suiteSummary,
    performanceAnomalies,
    runtimeConstraintFailures: normalizedTests.filter(
      (item) => item.runtimeDependent && (item.status === "failed" || item.status === "timedOut"),
    ),
  };
}

function buildUnexpectedBehavior(playwrightSummary, coverageSummary) {
  const unexpected = [];

  for (const failure of playwrightSummary.failures) {
    if (failure.failureType === "environment-connection-refused") {
      unexpected.push({
        category: "Runtime constraint",
        module: failure.module,
        description: `${failure.testName} failed because the required backend service was unavailable (ECONNREFUSED).`,
        source: failure.testId,
      });
    }
  }

  for (const anomaly of playwrightSummary.performanceAnomalies) {
    unexpected.push({
      category: "Performance anomaly",
      module: anomaly.module,
      description: `${anomaly.testName} took ${formatMs(anomaly.durationMs)}, exceeding the 5s anomaly threshold.`,
      source: anomaly.testId,
    });
  }

  for (const gap of coverageSummary.moduleSummary.filter((entry) => entry.coveragePercent < 10)) {
    unexpected.push({
      category: "Coverage anomaly",
      module: gap.module,
      description: `${gap.module} backend coverage is only ${gap.coveragePercent}%, leaving critical logic weakly instrumented.`,
      source: "coverage-summary",
    });
  }

  return unexpected;
}

function buildRiskReevaluation(playwrightSummary, coverageSummary) {
  return HIGH_RISK_MODULES.map((moduleName) => {
    const failures = playwrightSummary.failures.filter((item) => item.module === moduleName);
    const runtimeFailures = playwrightSummary.runtimeConstraintFailures.filter((item) => item.module === moduleName);
    const moduleCoverage = coverageSummary.moduleSummary.find((entry) => entry.module === moduleName);
    const coverageGaps = coverageSummary.gapEntries.filter((entry) => entry.module === moduleName && entry.coveragePercent === 0);
    const performanceAnomalies = playwrightSummary.performanceAnomalies.filter((item) => item.module === moduleName);

    let likelihood = 3;
    if (failures.length > 0) likelihood += 1;
    if ((moduleCoverage?.coveragePercent ?? 0) < 15) likelihood += 1;
    likelihood = Math.min(likelihood, 5);

    let impact = ORIGINAL_RISK_SCORES[moduleName] >= 20 ? 5 : 4;
    if (moduleName === "Product catalog and listing") impact = 4;

    let detectability = 4;
    if ((moduleCoverage?.coveragePercent ?? 0) < 15) detectability = 2;
    else if ((moduleCoverage?.coveragePercent ?? 0) < 35) detectability = 3;

    const updatedRiskScore = likelihood * impact;
    const observedIssues = [
      ...failures.map((item) => `${item.failureType} x${item.frequencyAcrossRuns}`),
      runtimeFailures.length ? `${runtimeFailures.length} runtime-dependent failure(s)` : null,
      coverageGaps.length ? `${coverageGaps.length} uncovered backend file(s)` : null,
      performanceAnomalies.length ? `${performanceAnomalies.length} performance anomaly/anomalies` : null,
    ].filter(Boolean);

    return {
      module: moduleName,
      originalRiskScore: ORIGINAL_RISK_SCORES[moduleName],
      observedIssues,
      likelihood,
      impact,
      detectability,
      updatedRiskScore,
      justification:
        observedIssues.length > 0
          ? `Updated from empirical evidence: ${observedIssues.join("; ")}. Backend coverage for this module is ${moduleCoverage?.coveragePercent ?? 0}%.`
          : `No direct failures were recorded in the available reports, but backend coverage remains ${moduleCoverage?.coveragePercent ?? 0}% and still limits detectability.`,
    };
  });
}

function buildMetricsSummary(playwrightSummary, coverageSummary, goSummary) {
  const nonRuntimeCritical = playwrightSummary.tests.filter(
    (item) => item.suiteKey === "critical" && !item.runtimeDependent && item.status !== "skipped",
  );
  const criticalPassRate = percent(
    nonRuntimeCritical.filter((item) => item.status === "passed").length,
    nonRuntimeCritical.length,
  );

  const defectsByModule = HIGH_RISK_MODULES.map((moduleName) => ({
    module: moduleName,
    defects: playwrightSummary.failures.filter((item) => item.module === moduleName).length,
    riskLevel: ORIGINAL_RISK_SCORES[moduleName] >= 20 ? "High" : "High",
  }));

  const currentExecutionSeconds = round(
    playwrightSummary.suiteSummary.reduce((sum, suite) => sum + suite.durationMs, 0) / 1000,
  );
  const assignment2ExecutionTable = path.join(tablesDir, "execution-time.csv");
  let previousExecutionSeconds = 0;
  if (fs.existsSync(assignment2ExecutionTable)) {
    const lines = fs.readFileSync(assignment2ExecutionTable, "utf8").split("\n").slice(1);
    previousExecutionSeconds = round(
      lines
        .filter(Boolean)
        .map((line) => line.split(",")[3]?.replaceAll('"', "") ?? "0s")
        .map((value) => Number.parseFloat(value.replace("s", "")) || 0)
        .reduce((sum, value) => sum + value, 0),
    );
  }

  return {
    generatedAt: new Date().toISOString(),
    suiteSummary: playwrightSummary.suiteSummary,
    goSummary,
    highRiskCoverage: coverageSummary.moduleSummary,
    defectsByModule,
    defectCount: playwrightSummary.failures.length,
    criticalPassRate,
    flakyRate: percent(playwrightSummary.flakyTests.length, playwrightSummary.tests.length || 1),
    currentExecutionSeconds,
    previousExecutionSeconds,
    executionDeltaSeconds: round(currentExecutionSeconds - previousExecutionSeconds),
    pipelineRuntimeSeconds: round(currentExecutionSeconds + goSummary.totalElapsedSeconds),
  };
}

function writeBarChart(filePath, rows, { title, valueSuffix = "%", fill = "#0f766e" }) {
  const width = 840;
  const height = 360;
  const maxValue = Math.max(...rows.map((row) => row.value), 1);
  const barWidth = 110;
  const gap = 28;
  const bars = rows
    .map((row, index) => {
      const barHeight = (row.value / maxValue) * 220;
      const x = 50 + index * (barWidth + gap);
      const y = 280 - barHeight;
      return `
        <rect x="${x}" y="${y}" width="${barWidth}" height="${barHeight}" fill="${fill}" rx="12" />
        <text x="${x + barWidth / 2}" y="${y - 8}" font-size="12" text-anchor="middle" fill="#0f172a">${row.value}${valueSuffix}</text>
        <text x="${x + barWidth / 2}" y="305" font-size="12" text-anchor="middle" fill="#334155">${row.label}</text>
      `;
    })
    .join("");

  fs.writeFileSync(
    filePath,
    `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
      <rect width="100%" height="100%" fill="#f8fafc" />
      <text x="40" y="32" font-size="20" font-family="Arial" fill="#0f172a">${title}</text>
      <line x1="40" y1="280" x2="780" y2="280" stroke="#cbd5e1" />
      ${bars}
    </svg>`,
  );
}

function writeLineChart(filePath, rows, title) {
  const width = 820;
  const height = 320;
  const maxValue = Math.max(...rows.map((row) => row.value), 1);
  const points = rows.map((row, index) => ({
    ...row,
    x: 80 + index * 170,
    y: 240 - (row.value / maxValue) * 140,
  }));
  const polyline = points.map((point) => `${point.x},${point.y}`).join(" ");
  const pointSvg = points
    .map(
      (point) => `
        <circle cx="${point.x}" cy="${point.y}" r="5" fill="#d97706" />
        <text x="${point.x}" y="${point.y - 10}" font-size="12" text-anchor="middle" fill="#0f172a">${point.value}s</text>
        <text x="${point.x}" y="270" font-size="12" text-anchor="middle" fill="#334155">${point.label}</text>
      `,
    )
    .join("");

  fs.writeFileSync(
    filePath,
    `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
      <rect width="100%" height="100%" fill="#fff7ed" />
      <text x="40" y="32" font-size="20" font-family="Arial" fill="#0f172a">${title}</text>
      <polyline fill="none" stroke="#d97706" stroke-width="3" points="${polyline}" />
      <line x1="60" y1="240" x2="760" y2="240" stroke="#fed7aa" />
      ${pointSvg}
    </svg>`,
  );
}

const playwrightRaw = loadPlaywrightTests();
const playwrightSummary = summarizePlaywright(playwrightRaw);
const coverageSummary = parseGoCoverage();
const goSummary = parseGoTestSummary();
const unexpectedBehavior = buildUnexpectedBehavior(playwrightSummary, coverageSummary);
const riskReevaluation = buildRiskReevaluation(playwrightSummary, coverageSummary);
const metricsSummary = buildMetricsSummary(playwrightSummary, coverageSummary, goSummary);

fs.writeFileSync(path.join(logsDir, "failed-tests-summary.json"), `${JSON.stringify(playwrightSummary.failures, null, 2)}\n`);
fs.writeFileSync(path.join(logsDir, "flaky-tests-summary.json"), `${JSON.stringify(playwrightSummary.flakyTests, null, 2)}\n`);
fs.writeFileSync(path.join(logsDir, "coverage-summary.json"), `${JSON.stringify(coverageSummary, null, 2)}\n`);
fs.writeFileSync(path.join(logsDir, "unexpected-behavior-summary.json"), `${JSON.stringify(unexpectedBehavior, null, 2)}\n`);
fs.writeFileSync(path.join(logsDir, "module-risk-reevaluation.json"), `${JSON.stringify(riskReevaluation, null, 2)}\n`);
fs.writeFileSync(path.join(logsDir, "midterm-metrics-summary.json"), `${JSON.stringify(metricsSummary, null, 2)}\n`);

fs.writeFileSync(
  path.join(tablesDir, "failed-tests.csv"),
  toCsv([
    ["test_id", "test_name", "module", "failure_type", "frequency_across_runs", "runtime_dependent", "failure_message"],
    ...playwrightSummary.failures.map((item) => [
      item.testId,
      item.testName,
      item.module,
      item.failureType,
      item.frequencyAcrossRuns,
      item.runtimeDependent,
      item.failureMessage,
    ]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "flaky-tests.csv"),
  toCsv([
    ["test_id", "test_name", "module", "passes", "failures", "suspected_cause"],
    ...playwrightSummary.flakyTests.map((item) => [
      item.testId,
      item.testName,
      item.module,
      item.passes,
      item.failures,
      item.suspectedCause,
    ]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "module-coverage.csv"),
  toCsv([
    ["module", "coverage_percent", "covered_statements", "total_statements", "below_70_percent"],
    ...coverageSummary.moduleSummary.map((item) => [
      item.module,
      item.coveragePercent,
      item.coveredStatements,
      item.totalStatements,
      item.belowThreshold,
    ]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "automation-coverage.csv"),
  toCsv([
    ["module_feature", "coverage_percent", "covered_statements", "total_statements", "below_70_percent"],
    ...coverageSummary.moduleSummary.map((item) => [
      item.module,
      item.coveragePercent,
      item.coveredStatements,
      item.totalStatements,
      item.belowThreshold,
    ]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "coverage-gaps.csv"),
  toCsv([
    ["module", "file", "coverage_percent", "uncovered_functions"],
    ...coverageSummary.gapEntries.map((item) => [
      item.module,
      item.file,
      item.coveragePercent,
      item.uncoveredFunctions.join("; "),
    ]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "unexpected-behavior.csv"),
  toCsv([
    ["category", "module", "description", "source"],
    ...unexpectedBehavior.map((item) => [item.category, item.module, item.description, item.source]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "risk-reevaluation.csv"),
  toCsv([
    ["module", "original_risk_score", "observed_issues", "likelihood", "impact", "detectability", "updated_risk_score", "justification"],
    ...riskReevaluation.map((item) => [
      item.module,
      item.originalRiskScore,
      item.observedIssues.join("; "),
      item.likelihood,
      item.impact,
      item.detectability,
      item.updatedRiskScore,
      item.justification,
    ]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "midterm-metrics-summary.csv"),
  toCsv([
    ["metric", "value", "notes"],
    ["High-risk average backend coverage", `${coverageSummary.averageCoveragePercent}%`, "Derived from logs/go-coverage.out"],
    ["Detected defects", metricsSummary.defectCount, "Derived from Playwright JSON reports"],
    ["Critical non-runtime pass rate", `${metricsSummary.criticalPassRate}%`, "Runtime-dependent live API cases excluded from blocking gate evaluation"],
    ["Flaky rate", `${metricsSummary.flakyRate}%`, "Across all parsed Playwright test results"],
    ["Current execution time", `${metricsSummary.currentExecutionSeconds}s`, "Aggregated from parsed Playwright reports"],
    ["Previous execution time", `${metricsSummary.previousExecutionSeconds}s`, "Derived from prior Assignment 2 execution table"],
    ["Pipeline runtime", `${metricsSummary.pipelineRuntimeSeconds}s`, "Playwright durations plus Go test elapsed time"],
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "defects-vs-risk.csv"),
  toCsv([
    ["module_feature", "high_risk_level", "defects_found", "runtime_dependent_failures", "notes"],
    ...metricsSummary.defectsByModule.map((item) => [
      item.module,
      item.riskLevel,
      item.defects,
      playwrightSummary.runtimeConstraintFailures.filter((failure) => failure.module === item.module).length,
      "Derived from parsed Playwright JSON failures",
    ]),
  ]),
);

fs.writeFileSync(
  path.join(tablesDir, "execution-time.csv"),
  toCsv([
    ["metric", "duration_seconds", "notes"],
    ["Assignment 2 baseline", metricsSummary.previousExecutionSeconds, "Derived from prior execution-time table"],
    ["Midterm parsed Playwright runtime", metricsSummary.currentExecutionSeconds, "Sum of parsed Playwright report durations"],
    ["Midterm estimated pipeline runtime", metricsSummary.pipelineRuntimeSeconds, "Parsed Playwright time plus Go test elapsed time"],
  ]),
);

writeBarChart(
  path.join(chartsDir, "automation-coverage.svg"),
  coverageSummary.moduleSummary.map((item) => ({ label: item.module.split(" ")[0], value: item.coveragePercent })),
  { title: "Midterm High-Risk Backend Coverage", valueSuffix: "%", fill: "#0f766e" },
);

writeBarChart(
  path.join(chartsDir, "defects-vs-risk.svg"),
  metricsSummary.defectsByModule.map((item) => ({ label: item.module.split(" ")[0], value: item.defects })),
  { title: "Detected Defects Per High-Risk Module", valueSuffix: "", fill: "#dc2626" },
);

writeLineChart(
  path.join(chartsDir, "execution-time.svg"),
  [
    { label: "Assignment 2", value: metricsSummary.previousExecutionSeconds },
    { label: "Midterm", value: metricsSummary.currentExecutionSeconds },
    { label: "Pipeline", value: metricsSummary.pipelineRuntimeSeconds },
  ],
  "Execution Time Comparison",
);

console.log("Midterm metrics artifacts generated.");
