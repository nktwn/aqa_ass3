import fs from "node:fs";
import path from "node:path";
import {
  average,
  logsDir,
  nowIso,
  percentile,
  qaApiBaseUrl,
  readJson,
  rootDir,
  round,
  runCommand,
  writeJson,
} from "./experimental-common.mjs";

const configPath = path.join(rootDir, "qa", "mutation", "mutants.json");
const config = readJson(configPath);
const mutationLogsDir = path.join(logsDir, "mutation");
const resultPath = path.join(logsDir, "mutation-results.json");

fs.mkdirSync(mutationLogsDir, { recursive: true });

function replaceOccurrence(source, search, replace, occurrence = 1) {
  let fromIndex = 0;
  let currentOccurrence = 0;

  while (true) {
    const matchIndex = source.indexOf(search, fromIndex);
    if (matchIndex === -1) {
      throw new Error(`Could not find occurrence ${occurrence} of target text.`);
    }

    currentOccurrence += 1;
    if (currentOccurrence === occurrence) {
      return (
        source.slice(0, matchIndex) +
        replace +
        source.slice(matchIndex + search.length)
      );
    }

    fromIndex = matchIndex + search.length;
  }
}

function summarize(results) {
  const byModule = new Map();
  for (const result of results) {
    const moduleStats = byModule.get(result.module) || { total: 0, killed: 0, survived: 0 };
    moduleStats.total += 1;
    if (result.outcome === "killed") {
      moduleStats.killed += 1;
    } else {
      moduleStats.survived += 1;
    }
    byModule.set(result.module, moduleStats);
  }

  const moduleScores = Array.from(byModule.entries()).map(([module, stats]) => ({
    module,
    totalMutants: stats.total,
    killedMutants: stats.killed,
    survivedMutants: stats.survived,
    mutationScorePercent: round((stats.killed / Math.max(1, stats.total)) * 100),
  }));

  const totals = results.reduce(
    (accumulator, result) => {
      accumulator.totalMutants += 1;
      if (result.outcome === "killed") {
        accumulator.killedMutants += 1;
      } else {
        accumulator.survivedMutants += 1;
      }
      return accumulator;
    },
    { totalMutants: 0, killedMutants: 0, survivedMutants: 0 },
  );

  return {
    moduleScores,
    totals: {
      ...totals,
      overallMutationScorePercent: round((totals.killedMutants / Math.max(1, totals.totalMutants)) * 100),
    },
  };
}

const results = [];

for (const mutant of config) {
  const targetFile = path.join(rootDir, mutant.file);
  const originalSource = fs.readFileSync(targetFile, "utf8");
  const mutatedSource = replaceOccurrence(
    originalSource,
    mutant.search,
    mutant.replace,
    mutant.occurrence || 1,
  );

  fs.writeFileSync(targetFile, mutatedSource);

  const logPath = path.join(mutationLogsDir, `${mutant.id}.log`);
  let commandResult;
  try {
    commandResult = runCommand(mutant.testCommand[0], mutant.testCommand.slice(1), {
      cwd: path.join(rootDir, "backend"),
      env: {
        GOCACHE: "/tmp/go-build",
        GOMODCACHE: "/tmp/go-mod",
      },
    });
  } finally {
    fs.writeFileSync(targetFile, originalSource);
  }

  const outcome = commandResult.exitCode === 0 ? "survived" : "killed";
  fs.writeFileSync(
    logPath,
    [
      `Mutant: ${mutant.id}`,
      `Module: ${mutant.module}`,
      `Type: ${mutant.mutantType}`,
      `Description: ${mutant.description}`,
      `Outcome: ${outcome}`,
      `Command: ${[mutant.testCommand[0], ...mutant.testCommand.slice(1)].join(" ")}`,
      "",
      "--- stdout ---",
      commandResult.stdout,
      "",
      "--- stderr ---",
      commandResult.stderr,
    ].join("\n"),
  );

  results.push({
    ...mutant,
    generatedAt: nowIso(),
    testPackage: mutant.testCommand.slice(1).join(" "),
    exitCode: commandResult.exitCode,
    outcome,
    logPath,
  });
}

const summary = summarize(results);
const payload = {
  generatedAt: nowIso(),
  baseUrl: qaApiBaseUrl,
  results,
  ...summary,
};

writeJson(resultPath, payload);

console.log(`Mutation summary written to ${resultPath}`);
