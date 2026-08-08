import { apiRequest } from '@/shared/api';

export function logout() {
  return apiRequest<void>('/auth/logout', {
    method: 'POST',
  });
}