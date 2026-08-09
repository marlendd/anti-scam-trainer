import { apiRequest } from '@/shared/api'

export interface LoginRequest {
    email: string
}

export function forgot(data: LoginRequest) {
    return apiRequest<void>('/auth/forgot-password', {
        method: 'POST',
        body: data,
    })
}
