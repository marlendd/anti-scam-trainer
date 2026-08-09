import type {
    AttemptActor,
    AttemptStatus,
} from '../model/types'

export interface ChoiceOptionDto {
    id: string
    text: string
}

export interface AttemptMessageDto {
    author: AttemptActor
    text: string
}

export interface AttemptScenarioDto {
    title: string
    description: string
    role: AttemptActor
    product: {
        title: string
        price: number
    }
}

export interface CurrentNodeDto {
    id: string
    author: AttemptActor
    text: string
    messages: AttemptMessageDto[]
    choices: ChoiceOptionDto[]
}

export interface AttemptHistoryItemDto {
    node: {
        id: string
        author: AttemptActor
        text: string
        messages: AttemptMessageDto[]
    }
    selected_choice: ChoiceOptionDto
    consequence: string
    answered_at: string
}

export interface AttemptResponseDto {
    id: string
    scenario_id: string
    status: AttemptStatus
    current_node_id?: string
    ending_id?: string
    score?: number
    started_at: string
    updated_at: string
    completed_at?: string
}

export interface AttemptStateResponseDto
    extends AttemptResponseDto {
    scenario: AttemptScenarioDto
    current_node?: CurrentNodeDto
    history: AttemptHistoryItemDto[]
}