import * as util from "util";
import { CronJob } from "cron";
import { logger } from "../logger";
import * as redis from "redis";
import { promisify } from "util";
import * as rp from "request-promise";

const client = redis.createClient({url: process.env["REDISCLOUD_URL"]});
const smembersAsync = promisify(client.smembers).bind(client);
const sremAsync = promisify(client.srem).bind(client);
const hgetAsync = promisify(client.hget).bind(client);
const delAsync = promisify(client.del).bind(client);

exports.name = "reap";
exports.describe = "find and purge expirable images";
exports.builder = {

};

exports.handler = async (argv) => {
  main(argv).catch((err) => {
    console.log(`Failed with error ${util.inspect(err)}`);
    process.exit(1);
  });
};

async function main(argv): Promise<any> {
  process.on('SIGTERM', function onSigterm () {
    logger.info(`Got SIGTERM, cleaning up`);
    process.exit();
  });

  const job = new CronJob({
    cronTime: "* * * * *",
    onTick: async () => {
      console.log("-----> beginning to reap expired images");

      const now = new Date().getTime();
      const images = await smembersAsync("current.images");
      console.log(`   there are ${images.length} total images to evaluate`);
      for (const image of images) {
        const expireAt = await hgetAsync(image, "expires");

        if (+expireAt > now) {
          const minutesLeft = (+expireAt - now) / 1000 / 60;
          console.log(`not expiring ${image} for another ~${Math.round(minutesLeft)} minute(s)`);
          continue;
        }

        const imageAndTag = image.split(":");
        const headers = {
          "Accept": "application/vnd.docker.distribution.manifest.v2+json",
        };

        // Get the manifest from the tag
        const getOptions = {
          method: "HEAD",
          uri: `https://replreg.is/v2/${imageAndTag[0]}/manifests/${imageAndTag[1]}`,
          headers,
          resolveWithFullResponse: true,
        }
        const getResponse = await rp(getOptions);

        const deleteURI = `https://replreg.is/v2/${imageAndTag[0]}/manifests/${getResponse.headers.etag.replace(/"/g,"")}`;

        // Remove from the registry
        const options = {
          method: "DELETE",
          uri: deleteURI,
          headers,
          simple: false,
        }

        await rp(options);

        console.log(`expiring ${image}`);
        await sremAsync("current.images", image);
        await delAsync(image);
      }
    },
    start: true,
  });

  job.start();
}
