import { apiRequest } from '@/shared/api';
import type {User} from "@/entities/user";


export function getCurrentUser() {
  return apiRequest<User>('/users/me');
}