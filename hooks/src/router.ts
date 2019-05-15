import * as _ from "lodash";

import * as express from "express";
import * as util from "util";
import * as uuid from "uuid";
import { logger } from "./logger";

export interface Response<T> {
  status: number;
  body: T;
  contentType?: string;
  headers?: object;
  filename?: string;
}

export interface RawResponse extends Response<string> {
}

export const Responses = {
  created(entity: any): RawResponse {
    return {
      status: 201,
      body: JSON.stringify(entity),
    };
  },
};

/*
 * This file contains express middleware functions
 * for Pre/Post request logging and response generation.
 */

export const onSuccess = (res: express.Response, reqId: string, statusCodeGetter?: () => number | undefined) => (result: any) => {

  if (result) {
    const statusToSend = result.status || (statusCodeGetter && statusCodeGetter()) || 200;
    const body = result.body || JSON.stringify(result);
    const contentType = result.contentType || "application/json";

    let bodyToLog = body;
    if (!bodyToLog) {
      bodyToLog = "";
    } else if (bodyToLog.length > 512) {
      bodyToLog = `${bodyToLog.substring(0, 512)} (... truncated, total ${bodyToLog.length} bytes)`;
    }
    if (res.statusCode !== 200) {
      logger.warn(`[${reqId}] WARN response already has statusCode ${res.statusCode}, a response might have already been sent!`);
      logger.warn(util.inspect(res));
    }
    logger.info(`[${reqId}] => ${statusToSend} ${bodyToLog}`);
    const respObj = res.status(statusToSend).type(contentType).set("X-Replreg-RequestId", reqId);
    if (result.filename) {
      respObj.attachment(result.filename);
    }
    if (result.headers) {
      _.forOwn(result.headers, (value, key) => {
        respObj.set(key!, value);
      });
    }
    respObj.send(body);
  } else {
    const statusToSend = (statusCodeGetter && statusCodeGetter()) || 200;
    logger.info(`[${reqId}] => ${statusToSend}`);
    res.status(statusToSend).set("X-Replreg-RequestId", reqId).json(result);
  }
};

export const onError = (res: express.Response, reqId: string) => (err: any) => {
  if (err.status) {
    handleFrameworkError(err, reqId, res);
  } else {
    handleUnexpectedError(err, reqId, res);
  }
};

function handleFrameworkError(err: any, reqId: string, res: express.Response) {
  // Structured error, specific status code.
  const errMsg = err.err ? err.err.message : err.message || "An unexpected error occurred";

  logger.info(`[${reqId}] !! ${err.status} ${errMsg} ${err.stack || util.inspect(err)}`);

  const errClass = err.constructor.name;
  const hasMeaningfulType = ["Object", "Error"].indexOf(errClass) === -1;

  const bodyToSend = {
    status: err.status,
    type: hasMeaningfulType ? errClass : "Error",
    error: errMsg,
    invalid: err.invalid,
  };
  res.status(err.status).set("X-Replreg-RequestId", reqId).json(bodyToSend);

}

function handleUnexpectedError(err: any, reqId: string, res: express.Response) {
  // Generic error, default code (500).
  const bodyToSend = {
    status: 500,
    error: err.message || "An unexpected error occurred",
    type: "Error",
  };
  if (
    (err.message ? err.message : "").indexOf("Can't set headers after they are sent") !== -1 ||
    (err.stack ? err.stack : "").indexOf("Can't set headers after they are sent") !== -1) {
    logger.error("Middleware error, current response object is", util.inspect(res));
  }
  logger.error(`[${reqId}] !! 500 ${err.stack || err.message || util.inspect(err)}`);
  res.status(500).set("X-Replreg-RequestId", reqId).json(bodyToSend);
}

export const preRequest = (req: express.Request, reqId: string) => {
  logger.info(`[${reqId}] <- ${req.method} ${req.originalUrl}`);
  if (!_.isEmpty(req.body)) {
    let bodyString = JSON.stringify(req.body);
    if (bodyString.length > 512) {
      bodyString = `${bodyString.substring(0, 512)} (... truncated, total ${bodyString.length} bytes)`;
    }
    logger.info(`[${reqId}] <- ${bodyString}`);
  }
};

export const requestId = (req: express.Request, handlerName: string) => {
  const id = `${handlerName}:${uuid.v4().replace("-", "").substring(0, 8)}`;

  const clientID: string = <string> req.headers["x-request-uuid"];

  if (clientID) {
    return `${id}.${clientID.substring(0, 8)}`;
  }

  return id;
};

export const wrapRoute = (route, handlerName) =>
  (req: express.Request, res: express.Response) => {
    const reqId = requestId(req, handlerName);
    preRequest(req, reqId);
    route.handler(req)
      .then(onSuccess(res, reqId))
      .catch(onError(res, reqId));
  };

export function register(route, handler, app) {
  // Register this route and callback with express.
  if (route.method === "get") {
    logger.debug(`GET    '${route.path}'`);
    app.get(route.path, handler);
  } else if (route.method === "post") {
    logger.debug(`POST   '${route.path}'`);
    app.post(route.path, handler);
  } else if (route.method === "put") {
    logger.debug(`PUT    '${route.path}'`);
    app.put(route.path, handler);
  } else if (route.method === "delete") {
    logger.debug(`DELETE '${route.path}'`);
    app.delete(route.path, handler);
  } else {
    logger.debug(`Unhandled HTTP method: '${route.method}'`);
  }
}
