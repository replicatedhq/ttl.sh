import * as util from "util";
import { Server } from "../server/server";
import { logger } from "../logger";

exports.name = "hooks";
exports.describe = "Start and run the hook api server";
exports.builder = {

};

exports.handler = async (argv) => {
  main(argv).catch((err) => {
    console.log(`Failed with error ${util.inspect(err)}`);
    process.exit(1);
  });
};

async function main(argv): Promise<any> {
  process.on("SIGTERM", function onSigterm() {
    logger.info(`Got SIGTERM, cleaning up`);
    process.exit();
  });

  await new Server().start();
}
