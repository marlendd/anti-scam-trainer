import { useQuery } from '@tanstack/react-query'

import { ApiError, apiRequest } from '@/shared/api'

import type {
    AttemptStatus,
    TrainingAttempt,
    TrainingAttemptState,
} from '../model/types'

interface AttemptResponseDto {
    id: string
    scenario_id: string
    status: AttemptStatus
    current_node_id?: string
    ending_id?: string
    score?: number
    started_at: string
    updated_at: string
    completed_at?: string
}

interface AttemptStateResponseDto extends AttemptResponseDto {
    current_node?: {
        id: string
        author: string
        text: string
        choices: {
            id: string
            text: string
        }[]
    }
}

function mapAttempt(dto: AttemptResponseDto): TrainingAttempt {
    return {
        id: dto.id,
        scenarioId: dto.scenario_id,
        status: dto.status,
        currentNodeId: dto.current_node_id,
        endingId: dto.ending_id,
        score: dto.score,
        startedAt: dto.started_at,
        updatedAt: dto.updated_at,
        completedAt: dto.completed_at,
    }
}

function mapAttemptState(
    dto: AttemptStateResponseDto,
): TrainingAttemptState {
    return {
        ...mapAttempt(dto),

        currentNode: dto.current_node
            ? {
                  id: dto.current_node.id,
                  author: dto.current_node.author,
                  text: dto.current_node.text,
                  choices: dto.current_node.choices.map((choice) => ({
                      id: choice.id,
                      text: choice.text,
                  })),
              }
            : undefined,
    }
}

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
        if (!(error instanceof ApiError) || error.status !== 404) {
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
        if (error instanceof ApiError && error.status === 409) {
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
                throw new Error('Scenario ID is required')
            }

            return getOrStartAttempt(scenarioId)
        },
        enabled: scenarioId !== null,
        retry: false,
    })

    const attemptId = attemptQuery.data?.id ?? null

    const stateQuery = useQuery({
        queryKey: [
            'training-attempt',
            'state',
            attemptId,
        ],
        queryFn: () => {
            if (!attemptId) {
                throw new Error('Attempt ID is required')
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
            (
                attemptQuery.isPending ||
                (
                    attemptQuery.isSuccess &&
                    stateQuery.isPending
                )
            ),

        isError:
            attemptQuery.isError ||
            stateQuery.isError,

        error:
            attemptQuery.error ??
            stateQuery.error,
    }
}