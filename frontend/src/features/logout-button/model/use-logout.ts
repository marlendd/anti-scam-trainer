import { useMutation, useQueryClient } from '@tanstack/react-query'

import { currentUserQueryKey } from '@/entities/user'

import { logout } from '../api/logout'

export function useLogout() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: logout,

        onSuccess: () => {
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