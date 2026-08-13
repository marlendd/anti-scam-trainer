import { apiRequest } from '@/shared/api'

import type { TrainingAttemptState } from '../model/types'

export function restartAttempt(
    scenarioId: any,
) {
    return apiRequest<TrainingAttemptState>(
        `/scenarios/${scenarioId}/attempts/restart`,
        {
            method: 'POST',
        },
    )
}