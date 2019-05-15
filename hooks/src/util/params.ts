export class param {
  public static get(envName: string): string {
    return process.env[envName] || "";
  }
}
