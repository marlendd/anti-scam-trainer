import { useQuery } from '@tanstack/react-query'

import { getActiveAttempt } from './get-active-attempt'

export function useActiveAttempt(scenarioId: string | null) {
    return useQuery({
        queryKey: ['training-attempt', 'active', scenarioId],
        queryFn: () => {
            if (!scenarioId) {
                throw new Error('Scenario ID is required')
            }

            return getActiveAttempt(scenarioId)
        },
        enabled: scenarioId !== null,
        retry: false,
    })
}