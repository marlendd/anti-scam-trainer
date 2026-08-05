// src/entities/training-scenario/index.ts

export {
  fakeDeliveryScenario,
  trainingScenarioMocks,
} from './model/mocks'

export type {
  ScenarioAnalysis,
  ScenarioMessage,
  ScenarioParticipant,
  ScenarioParticipantRole,
  ScenarioProduct,
  ScenarioRedFlag,
  TrainingScenario,
} from './model/types'

export { mockScenarioRepository } from './api/mockScenarioRepository'

export type { ScenarioRepository } from './api/scenarioRepository'