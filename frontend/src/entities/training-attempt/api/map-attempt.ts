import type { TrainingAttempt, TrainingAttemptState } from '../model/types'

import type { AttemptResponseDto, AttemptStateResponseDto } from './types'

export function mapAttempt(dto: AttemptResponseDto): TrainingAttempt {
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

export function mapAttemptState(dto: AttemptStateResponseDto): TrainingAttemptState {
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
