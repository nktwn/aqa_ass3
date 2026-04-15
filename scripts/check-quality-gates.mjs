import fs from "node:fs";
import path from "node:path";

const rootDir = process.cwd();
const logsDir = path.join(rootDir, "logs");
const tablesDir = path.join(rootDir, "qa-docs", "tables");

function readJson(fileName) {
  return JSON.parse(fs.readFileSync(path.join(logsDir, fileName), "utf8"));
}

const metrics = readJson("midterm-metrics-summary.json");
const coverage = readJson("coverage-summary.json");
const failures = readJson("failed-tests-summary.json");
const flakyTests = readJson("flaky-tests-summary.json");

const nonRuntimeFailures = failures.filter((item) => !item.runtimeDependent);
const averageCoverage = coverage.averageCoveragePercent;
const smokeSuite = metrics.suiteSummary.find((suite) => suite.key === "smoke");

const gates = [
  {
    id: "QG-MID-001",
    criterion: "Backend Go test pass rate",
    threshold: ">= 95%",
    actual: `${metrics.goSummary.testPassRate}%`,
    passed: metrics.goSummary.testPassRate >= 95,
    rationale: "Unit and integration evidence must stay mostly green to trust backend coverage outputs.",
  },
  {
    id: "QG-MID-002",
    criterion: "Critical non-runtime failures",
    threshold: "0",
    actual: String(nonRuntimeFailures.length),
    passed: nonRuntimeFailures.length === 0,
    rationale: "Blocking failures in high-risk flows should not be masked by unrelated infrastructure constraints.",
  },
  {
    id: "QG-MID-003",
    criterion: "Average high-risk backend coverage",
    threshold: ">= 20%",
    actual: `${averageCoverage}%`,
    passed: averageCoverage >= 20,
    rationale: "Midterm coverage should materially improve visibility across the five highest-risk modules.",
  },
  {
    id: "QG-MID-004",
    criterion: "Observed flaky rate",
    threshold: "<= 5%",
    actual: `${metrics.flakyRate}%`,
    passed: metrics.flakyRate <= 5,
    rationale: "Empirical analysis becomes unreliable when too many tests oscillate across runs.",
  },
  {
    id: "QG-MID-005",
    criterion: "Non-runtime smoke pass rate",
    threshold: ">= 85%",
    actual: smokeSuite ? `${smokeSuite.passRate}%` : "n/a",
    passed: smokeSuite ? smokeSuite.passRate >= 85 : false,
    rationale: "Smoke feedback should remain useful even when optional live services are not enabled.",
  },
];

const overallPass = gates.every((gate) => gate.passed);
const evaluation = {
  generatedAt: new Date().toISOString(),
  overallPass,
  gates,
  assessment: {
    thresholdsTooStrictOrLenient:
      averageCoverage < 20
        ? "The coverage threshold would be too strict for the current repository state, but 20% is still defensible as an incremental midterm target."
        : "The thresholds are balanced: strict enough to catch regressions without blocking the repo on known runtime-dependent API gaps.",
    failureCause:
      nonRuntimeFailures.length > 0
        ? "Current blocking failures point to test or product quality issues in non-runtime-dependent paths."
        : "The remaining observed failures are dominated by infrastructure/runtime availability rather than functional regressions in mock-backed flows.",
    insufficientTests:
      coverage.moduleSummary.some((entry) => entry.coveragePercent < 10)
        ? "Yes. At least one high-risk module still has less than 10% backend coverage, which means the test portfolio is still thin in those areas."
        : "Current backend instrumentation is broad enough for a midterm baseline, though more depth is still desirable.",
  },
};

fs.writeFileSync(path.join(logsDir, "quality-gates-summary.json"), `${JSON.stringify(evaluation, null, 2)}\n`);
fs.writeFileSync(
  path.join(tablesDir, "quality-gates-evaluation.csv"),
  [
    ["gate_id", "criterion", "threshold", "actual", "passed", "rationale"],
    ...gates.map((gate) => [gate.id, gate.criterion, gate.threshold, gate.actual, gate.passed, gate.rationale]),
    ["ASSESSMENT", "Threshold fit", "n/a", evaluation.assessment.thresholdsTooStrictOrLenient, "n/a", "Midterm interpretation"],
    ["ASSESSMENT", "Failure cause", "n/a", evaluation.assessment.failureCause, "n/a", "Midterm interpretation"],
    ["ASSESSMENT", "Insufficient tests", "n/a", evaluation.assessment.insufficientTests, "n/a", "Midterm interpretation"],
  ]
    .map((row) => row.map((value) => `"${String(value).replaceAll('"', '""')}"`).join(","))
    .join("\n") + "\n",
);

if (!overallPass) {
  for (const gate of gates.filter((item) => !item.passed)) {
    console.error(`QUALITY GATE FAILED: ${gate.id} (${gate.criterion}) actual=${gate.actual} threshold=${gate.threshold}`);
  }
  process.exit(1);
}

console.log("Quality gates passed.");
