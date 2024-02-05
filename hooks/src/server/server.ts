import "source-map-support/register";
import * as cors from "cors";
import {ServerLoader, ServerSettings} from "ts-express-decorators";
import {$log} from "ts-log-debug";
import Path = require("path");
import * as bodyParser from "body-parser";

let port = process.env.PORT;
if (port == null || port == "") {
  port = "8000";
}

@ServerSettings({
  rootDir: Path.resolve(__dirname),
  mount: {
    "/": "${rootDir}/../controllers/**/*.js",
  },
  acceptMimes: ["application/json"],
  port: Number(port),
  httpsPort: 0,
  debug: false,
  logger: {
    level: "warn",
  },
})

export class Server extends ServerLoader {
  /**
   * This method let you configure the middleware required by your application to works.
   * @returns {Server}
   */
  public async $onMountingMiddlewares(): Promise<null> {
    this.expressApp.enable("trust proxy");  // so we get the real ip from the ELB in amaazon

    this.use(bodyParser.json({
      type: [
        "application/json",
        "application/vnd.docker.distribution.events.v1+json",
      ],
    }));

    this.use(cors());

    if (process.env["NODE_ENV"] === "production") {
      $log.level = "OFF";
    }

    return null;
  }

  public $onReady() {
      console.log("Server started...");
  }

  public $onServerInitError(err) {
      console.error(err);
  }
}
