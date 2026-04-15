import fs from "node:fs";
import path from "node:path";

const reportArg = process.argv[2];
if (!reportArg) {
  console.error("Usage: node scripts/print-playwright-failures.mjs <report.json>");
  process.exit(1);
}

const reportPath = path.resolve(process.cwd(), reportArg);
if (!fs.existsSync(reportPath)) {
  console.error(`Report file not found: ${reportPath}`);
  process.exit(1);
}

const report = JSON.parse(fs.readFileSync(reportPath, "utf8"));

function collect(node, parentTitle = "") {
  const currentTitle = [parentTitle, node.title].filter(Boolean).join(" > ");
  const rows = [];

  for (const suite of node.suites ?? []) {
    rows.push(...collect(suite, currentTitle));
  }

  for (const spec of node.specs ?? []) {
    const fullTitle = [currentTitle, spec.title].filter(Boolean).join(" > ");
    for (const test of spec.tests ?? []) {
      const result = test.results?.at(-1) ?? {};
      const status = result.status ?? "unknown";
      const message = (result.error?.message || result.errors?.[0]?.message || "")
        .split("\n")
        .slice(0, 4)
        .join("\n");
      rows.push({
        title: fullTitle,
        status,
        message,
      });
    }
  }

  return rows;
}

const rows = (report.suites ?? []).flatMap((suite) => collect(suite));
const failed = rows.filter((row) => row.status === "failed" || row.status === "timedOut");
const skipped = rows.filter((row) => row.status === "skipped");

console.log(`Playwright report: ${path.relative(process.cwd(), reportPath)}`);
console.log(`Total parsed tests: ${rows.length}`);
console.log(`Failed/timedOut: ${failed.length}`);
console.log(`Skipped: ${skipped.length}`);

if (failed.length) {
  console.log("\nFailed tests:");
  for (const item of failed) {
    console.log(`- ${item.title}`);
    if (item.message) {
      console.log(item.message);
    }
    console.log("");
  }
}

if (skipped.length) {
  console.log("Skipped tests:");
  for (const item of skipped) {
    console.log(`- ${item.title}`);
  }
}
