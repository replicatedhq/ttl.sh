import * as Express from "express";
import {
  Controller,
  Get,
  Res } from "ts-express-decorators";

interface ErrorResponse {
  error: any;
}

@Controller("/healthz")
export class HealthzAPI {
/**
 * /healthz handler
 *
 * @param request
 * @param response
 * @returns {{id: any, name: string}}
 */
  @Get("")
  public async check(
    @Res() response: Express.Response,
  ): Promise<{}> {
    response.status(200);
    return {};
  }
}
