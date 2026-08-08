import { apiRequest } from '@/shared/api';

export interface RegisterRequest {
  email: string;
  password: string;
}

export function register(data: RegisterRequest) {
  return apiRequest<void>('/auth/register', {
    method: 'POST',
    body: data,
  });
}