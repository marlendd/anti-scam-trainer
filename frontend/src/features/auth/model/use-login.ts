import { useMutation, useQueryClient } from '@tanstack/react-query'

import { currentUserQueryKey } from '@/entities/user'

import { login } from '../api/login'

export function useLogin() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: login,

        onSuccess: async () => {
            await queryClient.invalidateQueries({
                queryKey: currentUserQueryKey,
            })
        },
    })
}
