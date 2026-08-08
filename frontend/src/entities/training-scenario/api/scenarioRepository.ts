// entities/training-scenario/api/scenarioRepository.ts

import type { TrainingScenario } from '../model/types'

export type ScenarioRepository = {
    getById: (scenarioId: string) => Promise<TrainingScenario>
}
