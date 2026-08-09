// entities/user/api/use-current-user.ts

import { useQuery } from '@tanstack/react-query'

import { getCurrentUser } from './get-current-user'

export const currentUserQueryKey = ['current-user'] as const

export function useCurrentUser() {
    return useQuery({
        queryKey: currentUserQueryKey,
        queryFn: getCurrentUser,
    })
}
