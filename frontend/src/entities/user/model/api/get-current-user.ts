import { apiRequest } from '@/shared/api';
import type {User} from "@/entities/user/model/types.ts";


export function getCurrentUser() {
  return apiRequest<User>('/users/me');
}