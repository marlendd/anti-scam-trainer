// entities/training-scenario/api/useScenarios.ts

import { useQuery } from '@tanstack/react-query'

import { getScenarios } from './get-scenarios'

import type { TrainingRole } from '../model/types'

export const scenariosQueryKey = (role: TrainingRole | null) => ['scenarios', role] as const

export function useScenarios(role: TrainingRole | null) {
    return useQuery({
        queryKey: scenariosQueryKey(role),
        queryFn: () => getScenarios(role as TrainingRole),
        enabled: role !== null,
    })
}
