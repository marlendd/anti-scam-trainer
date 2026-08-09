import { useQuery } from '@tanstack/react-query'

import { ApiError, apiRequest } from '@/shared/api'

import type { TrainingAttempt, TrainingAttemptState } from '../model/types'
import { mapAttempt, mapAttemptState } from './map-attempt'
import type {
    AttemptResponseDto,
    AttemptStateResponseDto,
} from './types'

async function getActiveAttempt(
    scenarioId: string,
): Promise<TrainingAttempt> {
    const response = await apiRequest<AttemptResponseDto>(
        `/scenarios/${scenarioId}/attempts/active`,
    )

    return mapAttempt(response)
}

async function startAttempt(
    scenarioId: string,
): Promise<TrainingAttempt> {
    const response = await apiRequest<AttemptResponseDto>(
        `/scenarios/${scenarioId}/attempts`,
        {
            method: 'POST',
        },
    )

    return mapAttempt(response)
}

async function getAttempt(
    attemptId: string,
): Promise<TrainingAttemptState> {
    const response = await apiRequest<AttemptStateResponseDto>(
        `/attempts/${attemptId}`,
    )

    return mapAttemptState(response)
}

async function getOrStartAttempt(
    scenarioId: string,
): Promise<TrainingAttempt> {
    try {
        return await getActiveAttempt(scenarioId)
    } catch (error) {
        if (
            !(error instanceof ApiError) ||
            error.status !== 404
        ) {
            throw error
        }
    }

    try {
        return await startAttempt(scenarioId)
    } catch (error) {
        /*
         * На случай гонки: между GET active и POST
         * активная попытка уже могла появиться.
         */
        if (
            error instanceof ApiError &&
            error.status === 409
        ) {
            return getActiveAttempt(scenarioId)
        }

        throw error
    }
}

export function useTrainingSession(
    scenarioId: string | null,
) {
    const attemptQuery = useQuery({
        queryKey: [
            'training-attempt',
            'active',
            scenarioId,
        ],
        queryFn: () => {
            if (!scenarioId) {
                throw new Error(
                    'Scenario ID is required',
                )
            }

            return getOrStartAttempt(scenarioId)
        },
        enabled: scenarioId !== null,
        retry: false,
    })

    const attemptId =
        attemptQuery.data?.id ?? null

    const stateQuery = useQuery({
        queryKey: [
            'training-attempt',
            'state',
            attemptId,
        ],
        queryFn: () => {
            if (!attemptId) {
                throw new Error(
                    'Attempt ID is required',
                )
            }

            return getAttempt(attemptId)
        },
        enabled: attemptId !== null,
        retry: false,
    })

    return {
        data: stateQuery.data ?? null,

        isPending:
            scenarioId !== null &&
            (attemptQuery.isPending ||
                (attemptQuery.isSuccess &&
                    stateQuery.isPending)),

        isError:
            attemptQuery.isError ||
            stateQuery.isError,

        error:
            attemptQuery.error ??
            stateQuery.error,
    }
}