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
    type: 'message'
    id: string
    senderId: string
    text: string
    delayMs?: number
}

export type TrainingRole = 'buyer' | 'seller'

export type TrainingScenarioStatus =
    | 'not_started'
    | 'in_progress'
    | 'completed'

export interface TrainingScenarioSummary {
    id: string
    logicalId: string
    version: number
    role: TrainingRole
    title: string
    description: string
    status: TrainingScenarioStatus
}

export interface ScenarioDto {
    id: string
    logical_id: string
    version: number
    role: TrainingRole
    title: string
    description: string
    status: TrainingScenarioStatus
}

export interface ScenariosResponseDto {
    items: ScenarioDto[]
}

export type ScenarioAnswerOption = {
    id: string
    text: string
    isCorrect: boolean

    /**
     * Потом сюда можно добавить объяснение,
     * изменение score и т.д.
     */
    feedback?: string
}

export type ScenarioChoice = {
    type: 'choice'
    id: string

    options: ScenarioAnswerOption[]

    /**
     * Какой ответ автоматически показывать
     * в режиме preview / glossary.
     */
    previewOptionId?: string
}

export type ScenarioTimelineItem = ScenarioMessage | ScenarioChoice

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

    playerParticipantId: string
    participants: ScenarioParticipant[]

    timeline: ScenarioTimelineItem[]

    analysis: ScenarioAnalysis
}
