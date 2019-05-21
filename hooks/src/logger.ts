import * as pino from "pino";
import * as process from "process";

function initLogger(): any {
  const logOptions = {
    name: process.env["LOG_NAME"] || "ttlsh",
    safe: true,
    prettyPrint: process.env["LOG_PRETTY"] || false,
  };

  if (process.env["LOG_PRETTY"]) {
    const logger = pino(logOptions);
    logger.level = process.env["LOG_LEVEL"] || "warn";
    return logger;
  } else {
    const logger = pino(logOptions)
    logger.level = process.env["LOG_LEVEL"] || "warn";
    return logger;
  }
}

export const logger = initLogger();
