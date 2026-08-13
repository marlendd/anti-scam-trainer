export type AttemptStatus =
    | 'in_progress'
    | 'completed'
    | 'aborted'

export type AttemptActor = 'buyer' | 'seller'

export interface ChoiceOption {
    id: string
    text: string
}

export interface AttemptMessage {
    author: AttemptActor
    text: string
}

export interface AttemptScenarioProduct {
    title: string
    price: number
}

export interface AttemptScenario {
    title: string
    description: string
    role: AttemptActor
    product: AttemptScenarioProduct
}

export interface CurrentNode {
    id: string
    author: AttemptActor
    text: string
    messages: AttemptMessage[]
    choices: ChoiceOption[]
}

export interface AttemptHistoryNode {
    id: string
    author: AttemptActor
    text: string
    messages: AttemptMessage[]
}

export interface AttemptHistoryItem {
    node: AttemptHistoryNode
    selectedChoice: ChoiceOption
    consequence: string
    answeredAt: string
}

export interface TrainingAttempt {
    id: string
    scenarioId: string
    status: AttemptStatus
    currentNodeId?: string
    endingId?: string
    ending?: AttemptEnding
    score?: number
    startedAt: string
    updatedAt: string
    completedAt?: string
}

export type AttemptEnding = {
    id: string
    header: string
    result: string
}

export interface TrainingAttemptState extends TrainingAttempt {
    scenario: AttemptScenario
    currentNode?: CurrentNode
    ending?: AttemptEnding | null
    history: AttemptHistoryItem[]
}