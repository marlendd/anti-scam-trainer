import { apiRequest } from '@/shared/api'

export type AttemptResult = {
    score?: number
}

type AttemptResultResponseDto = {
    score?: number
}

export async function getAttemptResult(
    attemptId: string,
): Promise<AttemptResult> {
    return apiRequest<AttemptResultResponseDto>(
        `/attempts/${attemptId}/result`,
    )
}