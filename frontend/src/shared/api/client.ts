import {ApiError} from './errors';

const API_BASE_URL = '/api';

type ApiRequestOptions = Omit<RequestInit, 'body'> & {
    body?: unknown;
};

export async function apiRequest<T>(
    path: string,
    options: ApiRequestOptions = {},
): Promise<T> {
    const {body, headers, ...rest} = options;

    const response = await fetch(`${API_BASE_URL}${path}`, {
        ...rest,
        credentials: 'same-origin',
        headers: {
            ...(body !== undefined && {'Content-Type': 'application/json'}),
            ...headers,
        },
        body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
        let errorBody: unknown;

        try {
            errorBody = await response.json();
        } catch {
            errorBody = await response.text();
        }

        throw new ApiError(response.status, errorBody);
    }

    if (response.status === 204) {
        return undefined as T;
    }

    return response.json() as Promise<T>;
}