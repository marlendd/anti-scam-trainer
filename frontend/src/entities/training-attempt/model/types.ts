export type AttemptStatus =
    | 'in_progress'
    | 'completed'
    | 'aborted'

export interface ChoiceOption {
    id: string
    text: string
}

export interface CurrentNode {
    id: string
    author: string
    text: string
    choices: ChoiceOption[]
}

export interface TrainingAttempt {
    id: string
    scenarioId: string
    status: AttemptStatus
    currentNodeId?: string
    endingId?: string
    score?: number
    startedAt: string
    updatedAt: string
    completedAt?: string
}

export interface TrainingAttemptState extends TrainingAttempt {
    currentNode?: CurrentNode
}