import { useMutation, useQueryClient } from '@tanstack/react-query'

import { currentUserQueryKey } from '@/entities/user'

import { logout } from '../api/logout'
import {clearAuthSession} from "@/shared/lib/auth-session.ts";

export function useLogout() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: logout,

        onSuccess: () => {
            clearAuthSession()

            queryClient.setQueryData(currentUserQueryKey, null)

            queryClient.removeQueries({
                queryKey: ['profile-progress'],
            })

            queryClient.removeQueries({
                queryKey: ['training-attempt'],
            })

            queryClient.removeQueries({
                queryKey: ['training-scenarios'],
            })
        },
    })
}