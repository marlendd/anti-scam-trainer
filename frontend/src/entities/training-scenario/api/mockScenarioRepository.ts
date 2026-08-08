// entities/training-scenario/api/mockScenarioRepository.ts

import { trainingScenarioMocks } from '../model/mocks'

import type { ScenarioRepository } from './scenarioRepository'

function wait(delayMs: number) {
    return new Promise<void>((resolve) => {
        window.setTimeout(resolve, delayMs)
    })
}

export const mockScenarioRepository: ScenarioRepository = {
    async getById(scenarioId) {
        await wait(500)

        const scenario = trainingScenarioMocks.find((item) => item.id === scenarioId)

        if (!scenario) {
            throw new Error(`Scenario "${scenarioId}" not found`)
        }

        return structuredClone(scenario)
    },
}
