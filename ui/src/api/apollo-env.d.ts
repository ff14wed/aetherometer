declare module "apollo-env" {
  export function fetch(
    input: RequestInfo,
    init?: RequestInit,
  ): Promise<Response>;
}