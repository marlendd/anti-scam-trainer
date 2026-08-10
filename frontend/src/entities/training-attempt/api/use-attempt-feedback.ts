import { useQuery } from '@tanstack/react-query'

import { getAttemptFeedback } from './get-attempt-feedback'

export function useAttemptFeedback(attemptId: string | null, enabled = true) {
    return useQuery({
        queryKey: ['training-attempt', 'feedback', attemptId],
        queryFn: () => {
            if (!attemptId) {
                throw new Error('Attempt ID is required')
            }

            return getAttemptFeedback(attemptId)
        },
        enabled: enabled && attemptId !== null,
    })
}
