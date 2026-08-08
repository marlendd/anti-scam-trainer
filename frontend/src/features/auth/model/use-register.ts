import { useMutation, useQueryClient } from '@tanstack/react-query'

import { currentUserQueryKey } from '@/entities/user'

import { register } from '../api/register'

export function useRegister() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: register,
        onSuccess: async () => {
            await queryClient.invalidateQueries({
                queryKey: currentUserQueryKey,
            })
        },
    })
}