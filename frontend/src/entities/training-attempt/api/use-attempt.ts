import { useQuery } from '@tanstack/react-query'

import { getAttempt } from './get-attempt'

export function useAttempt(attemptId: string | null) {
    return useQuery({
        queryKey: ['training-attempt', attemptId],
        queryFn: () => {
            if (!attemptId) {
                throw new Error('Attempt ID is required')
            }

            return getAttempt(attemptId)
        },
        enabled: attemptId !== null,
        retry: false,
    })
}
