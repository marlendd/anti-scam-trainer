import { useMutation, useQueryClient } from '@tanstack/react-query'

import { currentUserQueryKey } from '@/entities/user'
import { setAuthSession } from '@/shared/lib/auth-session'

import { login } from '../api/login'

export function useLogin() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: login,

        onSuccess: () => {
            setAuthSession()

            void queryClient.invalidateQueries({
                queryKey: currentUserQueryKey,
            })
        },
    })
}