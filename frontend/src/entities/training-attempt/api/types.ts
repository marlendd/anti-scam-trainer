import type { AttemptStatus } from '../model/types'

export interface ChoiceOptionDto {
    id: string
    text: string
}

export interface CurrentNodeDto {
    id: string
    author: string
    text: string
    choices: ChoiceOptionDto[]
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
    current_node?: CurrentNodeDto
}