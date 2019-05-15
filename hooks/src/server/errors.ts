import * as Express from "express";
import * as util from "util";
import {
  Err,
  IMiddlewareError,
  MiddlewareError,
  Next,
  Request,
  Response,
} from "ts-express-decorators";
import { logger } from "../logger";

export class HTTPError extends Error {
  public static requireMatch(field: string, pattern: RegExp, name: string) {
    const isValid = pattern.test(field);
    if (!isValid) {
      throw new HTTPError(400, {
        message: `Missing or invalid parameters: ${name}`,
        code: "bad_request",
        extra: {name},
      });
    }
  }
  constructor(
    public readonly status: number,
    public readonly body: any,
  ) {
    super((body && body.message) || "Internal Error");
  }
}

export class ServerError extends HTTPError {
  constructor() {
    super(500, {
      error: {
        message: "A server error has occurred",
      },
    });
  }
}

export class Unauthorized extends HTTPError {
  constructor() {
    super(401, {
      error: {
        message: "Unauthorized",
      },
    });
  }
}

export class Errors {
  public static Unauthorized = Unauthorized;
  public static ServerError = ServerError;
}

@MiddlewareError()
export class ErrorMiddleware implements IMiddlewareError {

  public use(
    @Err() error: any,
    @Request() request: Express.Request,
    @Response() response: Express.Response,
    @Next() next: Express.NextFunction,
  ): any {
    logger.debug("Handling error", error);

    if (response.headersSent) {
      logger.debug("Headers sent, skipping error handling");
      return next(error);
    }

    if (!(error instanceof HTTPError)) {
      // its an unhandled error so log it and then return a regular 500
      logger.error("Handling internal server error " + util.inspect(error));
      error = new ServerError();
    }

    response.status(error.status).send(JSON.stringify(error.body));
    return next();

  }
}
