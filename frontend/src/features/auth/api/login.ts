import { apiRequest } from '@/shared/api';

export interface LoginRequest {
  email: string;
  password: string;
}

export function login(data: LoginRequest) {
  return apiRequest<void>('/auth/login', {
    method: 'POST',
    body: data,
  });
}