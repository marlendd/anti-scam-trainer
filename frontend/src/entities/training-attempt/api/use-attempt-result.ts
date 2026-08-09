import { useQuery } from '@tanstack/react-query'

import { getAttemptResult } from './get-attempt-result'

export function useAttemptResult(attemptId: string | null, enabled = true) {
    return useQuery({
        queryKey: ['training-attempt', 'result', attemptId],
        queryFn: () => {
            if (!attemptId) {
                throw new Error('Attempt ID is required')
            }

            return getAttemptResult(attemptId)
        },
        enabled: enabled && attemptId !== null,
    })
}
