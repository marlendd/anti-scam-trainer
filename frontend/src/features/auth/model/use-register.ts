import { useMutation, useQueryClient } from '@tanstack/react-query'

import { currentUserQueryKey } from '@/entities/user'

import { register } from '../api/register'
import {setAuthSession} from "@/shared/lib/auth-session.ts";

export function useRegister() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: register,
        onSuccess: async () => {
            setAuthSession()

            await queryClient.invalidateQueries({
                queryKey: currentUserQueryKey,
            })
        },
    })
}
