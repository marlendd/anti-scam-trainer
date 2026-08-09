import { apiRequest } from '@/shared/api'

import type { AttemptStatus, TrainingAttempt } from '../model/types'

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

export async function startAttempt(scenarioId: string): Promise<TrainingAttempt> {
    const response = await apiRequest<AttemptResponseDto>(`/scenarios/${scenarioId}/attempts`, {
        method: 'POST',
    })

    return {
        id: response.id,
        scenarioId: response.scenario_id,
        status: response.status,
        currentNodeId: response.current_node_id,
        endingId: response.ending_id,
        score: response.score,
        startedAt: response.started_at,
        updatedAt: response.updated_at,
        completedAt: response.completed_at,
    }
}
