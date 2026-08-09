import { apiRequest } from '@/shared/api'

import type { TrainingAttemptState } from '../model/types'

import { mapAttemptState } from './map-attempt'
import type { AttemptStateResponseDto } from './types'

export async function getAttempt(
    attemptId: string,
): Promise<TrainingAttemptState> {
    const response = await apiRequest<AttemptStateResponseDto>(
        `/attempts/${attemptId}`,
    )

    return mapAttemptState(response)
}