import type { TestInfo } from "@playwright/test";

export type QaScenarioType =
  | "failure"
  | "edge"
  | "concurrency"
  | "invalid-user-behavior";

export type QaTestMetadata = {
  testId: string;
  module: string;
  scenarioType: QaScenarioType;
  inputData: string;
  expectedOutput: string;
  runtimeDependent?: boolean;
};

export async function attachQaMetadata(info: TestInfo, metadata: QaTestMetadata) {
  info.annotations.push(
    { type: "qa.test.id", description: metadata.testId },
    { type: "qa.module", description: metadata.module },
    { type: "qa.scenario", description: metadata.scenarioType },
    { type: "qa.input", description: metadata.inputData },
    { type: "qa.expected", description: metadata.expectedOutput },
    {
      type: "qa.runtimeDependent",
      description: String(Boolean(metadata.runtimeDependent)),
    },
  );

  await info.attach("qa-metadata", {
    body: Buffer.from(JSON.stringify(metadata, null, 2)),
    contentType: "application/json",
  });
}
