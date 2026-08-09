import { apiRequest } from '@/shared/api'

import type { TrainingAttempt } from '../model/types'

import { mapAttempt } from './map-attempt'
import type { AttemptResponseDto } from './types'

export async function getActiveAttempt(scenarioId: string): Promise<TrainingAttempt> {
    const response = await apiRequest<AttemptResponseDto>(
        `/scenarios/${scenarioId}/attempts/active`,
    )

    return mapAttempt(response)
}
