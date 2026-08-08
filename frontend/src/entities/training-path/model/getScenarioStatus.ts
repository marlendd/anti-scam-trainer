// entities/training-path/model/getScenarioStatus.ts

import type { TrainingScenario } from './types'

export type ScenarioStatus = 'completed' | 'available' | 'locked'

export function getScenarioStatus(
    scenarios: TrainingScenario[],
    scenarioIndex: number,
): ScenarioStatus {
    const scenario = scenarios[scenarioIndex]

    if (scenario.isCompleted) {
        return 'completed'
    }

    const firstIncompleteIndex = scenarios.findIndex((item) => !item.isCompleted)

    if (scenarioIndex === firstIncompleteIndex) {
        return 'available'
    }

    return 'locked'
}
