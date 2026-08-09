import { apiRequest } from '@/shared/api'

import type {TrainingScenarioSummary, ScenariosResponseDto} from '../model/types'

export type TrainingRole = 'buyer' | 'seller'

export async function getScenarios(
    role: TrainingRole,
): Promise<TrainingScenarioSummary[]> {
    const response = await apiRequest<ScenariosResponseDto>(
        `/scenarios?role=${role}`,
    )

    return response.items.map((scenario) => ({
        id: scenario.id,
        logicalId: scenario.logical_id,
        version: scenario.version,
        role: scenario.role,
        title: scenario.title,
        description: scenario.description,
        status: scenario.status,
    }))
}