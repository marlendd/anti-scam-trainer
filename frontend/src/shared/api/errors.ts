export class ApiError extends Error {
  status: number;
  body: unknown;

  constructor(status: number, body: unknown) {
    super(`API request failed with status ${status}`);

    this.name = 'ApiError';
    this.status = status;
    this.body = body;
  }
}