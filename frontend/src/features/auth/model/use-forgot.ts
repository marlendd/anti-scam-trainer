import { useMutation, useQueryClient } from '@tanstack/react-query'

import { currentUserQueryKey } from '@/entities/user'

import { forgot } from '../api/forgot'

export function useForgot() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: forgot,
        onSuccess: async () => {
            await queryClient.invalidateQueries({
                queryKey: currentUserQueryKey,
            })
        },
    })
}
