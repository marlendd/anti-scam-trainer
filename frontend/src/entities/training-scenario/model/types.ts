// src/entities/training-scenario/model/types.ts

export type ScenarioParticipantRole = 'buyer' | 'seller'

export type ScenarioParticipant = {
  id: string
  name: string
  role: ScenarioParticipantRole
  status?: string
}

export type ScenarioProduct = {
  id: string
  title: string
  price: number
  imageUrl?: string
}

export type ScenarioMessage = {
  id: string
  senderId: string
  text: string
  delayMs?: number
}

export type ScenarioRedFlag = {
  id: string
  title: string
  description: string
  accent?: string
}

export type ScenarioAnalysis = {
  title: string
  redFlags: ScenarioRedFlag[]
  safeActions: string[]
  goldenRule: string
}

export type TrainingScenario = {
  id: string
  title: string
  description: string
  product: ScenarioProduct
  analysis: ScenarioAnalysis
  playerParticipantId: string
  participants: ScenarioParticipant[]
  messages: ScenarioMessage[]
}