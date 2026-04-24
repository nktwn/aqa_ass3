import fs from "node:fs";
import path from "node:path";
import {
  chartsDir,
  ensureDir,
  logsDir,
  readJsonIfExists,
  rootDir,
  round,
  tablesDir,
  toCsv,
  writeText,
} from "./experimental-common.mjs";

ensureDir(tablesDir);
ensureDir(chartsDir);

const performance = readJsonIfExists(path.join(logsDir, "performance-summary.json"), {
  scenarios: [],
  bottlenecks: ["Performance execution has not been recorded yet."],
});
const mutation = readJsonIfExists(path.join(logsDir, "mutation-results.json"), {
  results: [],
  moduleScores: [],
  totals: { totalMutants: 0, killedMutants: 0, survivedMutants: 0, overallMutationScorePercent: 0 },
});
const chaos = readJsonIfExists(path.join(logsDir, "chaos-summary.json"), {
  scenarios: [],
  recommendations: ["Chaos execution has not been recorded yet."],
});

function writePerformanceTable() {
  const rows = [
    [
      "scenario_id",
      "type",
      "duration_seconds",
      "average_response_time_ms",
      "median_response_time_ms",
      "p95_response_time_ms",
      "throughput_rps",
      "error_rate_percent",
      "passed",
    ],
    ...performance.scenarios.map((scenario) => [
      scenario.id,
      scenario.type,
      scenario.durationSeconds,
      scenario.metrics.averageResponseTimeMs,
      scenario.metrics.medianResponseTimeMs,
      scenario.metrics.p95ResponseTimeMs,
      scenario.metrics.throughputRps,
      scenario.metrics.errorRatePercent,
      scenario.passed,
    ]),
  ];

  writeText(path.join(tablesDir, "assignment3-performance-summary.csv"), toCsv(rows));
}

function writeMutationTables() {
  const summaryRows = [
    ["mutant_id", "module", "mutant_type", "description", "outcome", "file"],
    ...mutation.results.map((result) => [
      result.id,
      result.module,
      result.mutantType,
      result.description,
      result.outcome,
      result.file,
    ]),
  ];

  const scoreRows = [
    ["module", "total_mutants", "killed_mutants", "survived_mutants", "mutation_score_percent"],
    ...mutation.moduleScores.map((score) => [
      score.module,
      score.totalMutants,
      score.killedMutants,
      score.survivedMutants,
      score.mutationScorePercent,
    ]),
    [
      "OVERALL",
      mutation.totals.totalMutants,
      mutation.totals.killedMutants,
      mutation.totals.survivedMutants,
      mutation.totals.overallMutationScorePercent,
    ],
  ];

  writeText(path.join(tablesDir, "assignment3-mutation-summary.csv"), toCsv(summaryRows));
  writeText(path.join(tablesDir, "assignment3-mutation-scores.csv"), toCsv(scoreRows));
}

function writeChaosTable() {
  const rows = [
    [
      "scenario_id",
      "module",
      "fault_type",
      "availability_during_fault_percent",
      "recovered",
      "mttr_ms",
      "recovery_behavior",
      "graceful_degradation",
    ],
    ...chaos.scenarios.map((scenario) => [
      scenario.id,
      scenario.module,
      scenario.faultType,
      scenario.availabilityDuringFaultPercent,
      scenario.recovered,
      scenario.mttrMs ?? "n/a",
      scenario.recoveryBehavior,
      scenario.gracefulDegradation,
    ]),
  ];

  writeText(path.join(tablesDir, "assignment3-chaos-summary.csv"), toCsv(rows));
}

function buildSvgChart({ title, labels, series, yLabel = "", maxValueOverride = null }) {
  const width = 920;
  const height = 420;
  const margin = { top: 60, right: 30, bottom: 100, left: 70 };
  const innerWidth = width - margin.left - margin.right;
  const innerHeight = height - margin.top - margin.bottom;
  const allValues = series.flatMap((item) => item.values);
  const maxValue = maxValueOverride ?? Math.max(1, ...allValues);
  const groupWidth = innerWidth / Math.max(1, labels.length);
  const barWidth = Math.min(42, (groupWidth - 20) / Math.max(1, series.length));
  const palette = ["#0b6e4f", "#c84c09", "#1d4ed8", "#9a3412"];

  const bars = [];
  const xLabels = [];
  const yTicks = [];

  for (let tick = 0; tick <= 5; tick += 1) {
    const value = (maxValue / 5) * tick;
    const y = margin.top + innerHeight - (value / maxValue) * innerHeight;
    yTicks.push(`
      <line x1="${margin.left}" y1="${y}" x2="${width - margin.right}" y2="${y}" stroke="#d7e3dd" stroke-width="1" />
      <text x="${margin.left - 10}" y="${y + 4}" text-anchor="end" font-size="12" fill="#334155">${round(value)}</text>
    `);
  }

  labels.forEach((label, labelIndex) => {
    const groupX = margin.left + labelIndex * groupWidth;
    xLabels.push(
      `<text x="${groupX + groupWidth / 2}" y="${height - 35}" text-anchor="middle" font-size="12" fill="#334155">${label}</text>`,
    );

    series.forEach((item, seriesIndex) => {
      const value = item.values[labelIndex] ?? 0;
      const barHeight = (value / maxValue) * innerHeight;
      const x = groupX + 10 + seriesIndex * barWidth;
      const y = margin.top + innerHeight - barHeight;
      bars.push(
        `<rect x="${x}" y="${y}" width="${barWidth - 4}" height="${barHeight}" fill="${palette[seriesIndex % palette.length]}" rx="4" />`,
      );
      bars.push(
        `<text x="${x + (barWidth - 4) / 2}" y="${y - 8}" text-anchor="middle" font-size="11" fill="#0f172a">${round(value)}</text>`,
      );
    });
  });

  const legend = series
    .map(
      (item, index) => `
        <rect x="${margin.left + index * 180}" y="22" width="14" height="14" fill="${palette[index % palette.length]}" rx="3" />
        <text x="${margin.left + index * 180 + 22}" y="34" font-size="12" fill="#334155">${item.label}</text>
      `,
    )
    .join("");

  return `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
  <rect width="${width}" height="${height}" fill="#f7fbf9" />
  <text x="${margin.left}" y="18" font-size="20" font-weight="700" fill="#0f172a">${title}</text>
  ${legend}
  ${yTicks.join("")}
  <line x1="${margin.left}" y1="${margin.top + innerHeight}" x2="${width - margin.right}" y2="${margin.top + innerHeight}" stroke="#64748b" stroke-width="1.5" />
  <line x1="${margin.left}" y1="${margin.top}" x2="${margin.left}" y2="${margin.top + innerHeight}" stroke="#64748b" stroke-width="1.5" />
  ${bars.join("")}
  ${xLabels.join("")}
  <text x="20" y="${margin.top - 12}" font-size="12" fill="#475569">${yLabel}</text>
</svg>
`;
}

function writeCharts() {
  if (performance.scenarios.length) {
    const labels = performance.scenarios.map((scenario) => scenario.id);
    writeText(
      path.join(chartsDir, "performance-response-time.svg"),
      buildSvgChart({
        title: "Assignment 3 Performance Response Times",
        labels,
        yLabel: "Milliseconds",
        series: [
          {
            label: "Average",
            values: performance.scenarios.map((scenario) => scenario.metrics.averageResponseTimeMs),
          },
          {
            label: "Median",
            values: performance.scenarios.map((scenario) => scenario.metrics.medianResponseTimeMs),
          },
          {
            label: "p95",
            values: performance.scenarios.map((scenario) => scenario.metrics.p95ResponseTimeMs),
          },
        ],
      }),
    );

    writeText(
      path.join(chartsDir, "performance-throughput.svg"),
      buildSvgChart({
        title: "Assignment 3 Performance Throughput and Error Rate",
        labels,
        yLabel: "RPS / Percent",
        series: [
          {
            label: "Throughput (rps)",
            values: performance.scenarios.map((scenario) => scenario.metrics.throughputRps),
          },
          {
            label: "Error Rate (%)",
            values: performance.scenarios.map((scenario) => scenario.metrics.errorRatePercent),
          },
        ],
      }),
    );
  }

  if (mutation.moduleScores.length) {
    writeText(
      path.join(chartsDir, "mutation-score.svg"),
      buildSvgChart({
        title: "Assignment 3 Mutation Score by Module",
        labels: mutation.moduleScores.map((item) => item.module),
        yLabel: "Percent",
        maxValueOverride: 100,
        series: [
          {
            label: "Mutation Score %",
            values: mutation.moduleScores.map((item) => item.mutationScorePercent),
          },
        ],
      }),
    );
  }

  if (chaos.scenarios.length) {
    writeText(
      path.join(chartsDir, "chaos-availability.svg"),
      buildSvgChart({
        title: "Assignment 3 Chaos Availability During Fault",
        labels: chaos.scenarios.map((item) => item.id),
        yLabel: "Percent",
        maxValueOverride: 100,
        series: [
          {
            label: "Availability During Fault %",
            values: chaos.scenarios.map((item) => item.availabilityDuringFaultPercent),
          },
        ],
      }),
    );
  }
}

function writeMarkdownSummaries() {
  const performanceMd = `# Performance Results

Generated from \`logs/performance-summary.json\`.

## Scenario Summary

| Scenario | Avg (ms) | Median (ms) | p95 (ms) | Throughput (rps) | Error Rate (%) | Passed |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
${performance.scenarios
  .map(
    (scenario) =>
      `| ${scenario.id} | ${scenario.metrics.averageResponseTimeMs} | ${scenario.metrics.medianResponseTimeMs} | ${scenario.metrics.p95ResponseTimeMs} | ${scenario.metrics.throughputRps} | ${scenario.metrics.errorRatePercent} | ${scenario.passed ? "Yes" : "No"} |`,
  )
  .join("\n")}

## Bottleneck Analysis

${performance.bottlenecks.map((item) => `- ${item}`).join("\n")}
`;

  const mutationMd = `# Mutation Results

Generated from \`logs/mutation-results.json\`.

## Module Scores

| Module | Total | Killed | Survived | Score (%) |
| --- | ---: | ---: | ---: | ---: |
${mutation.moduleScores
  .map(
    (score) =>
      `| ${score.module} | ${score.totalMutants} | ${score.killedMutants} | ${score.survivedMutants} | ${score.mutationScorePercent} |`,
  )
  .join("\n")}
| Overall | ${mutation.totals.totalMutants} | ${mutation.totals.killedMutants} | ${mutation.totals.survivedMutants} | ${mutation.totals.overallMutationScorePercent} |

## Surviving Mutants

${mutation.results
  .filter((result) => result.outcome === "survived")
  .map((result) => `- ${result.id} (${result.module}): ${result.description}`)
  .join("\n") || "- None."}
`;

  const chaosMd = `# Chaos Results

Generated from \`logs/chaos-summary.json\`.

## Scenario Summary

| Scenario | Fault | Availability During Fault (%) | Recovered | MTTR (ms) |
| --- | --- | ---: | --- | ---: |
${chaos.scenarios
  .map(
    (scenario) =>
      `| ${scenario.id} | ${scenario.faultType} | ${scenario.availabilityDuringFaultPercent} | ${scenario.recovered ? "Yes" : "No"} | ${scenario.mttrMs ?? "n/a"} |`,
  )
  .join("\n")}

## Resilience Recommendations

${chaos.recommendations.map((item) => `- ${item}`).join("\n")}
`;

  writeText(path.join(rootDir, "qa-docs", "performance-results.md"), performanceMd);
  writeText(path.join(rootDir, "qa-docs", "mutation-results.md"), mutationMd);
  writeText(path.join(rootDir, "qa-docs", "chaos-results.md"), chaosMd);
}

writePerformanceTable();
writeMutationTables();
writeChaosTable();
writeCharts();
writeMarkdownSummaries();

console.log("Assignment 3 tables and charts generated.");
