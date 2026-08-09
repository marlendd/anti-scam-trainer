import type { ProgressRole } from '../model/types'

export type RoleProgressDto = {
    role: ProgressRole
    completed_count?: number
    in_progress_count?: number
    total_started?: number
}

export type CategoryProgressDto = {
    category?: string
    count?: number
}

export type CategoriesProgressResponseDto = {
    total_completed?: number
    stats?: CategoryProgressDto[]
}

export type PuzzleFragmentDto = {
    fragment_id?: string
    earned_at?: string
}

export type PuzzleProgressResponseDto = {
    earned_count?: number
    total_count?: number
    fragments?: PuzzleFragmentDto[]
}

export type RankHistoryPointDto = {
    date?: string
    rank?: number
}

export type RankHistoryResponseDto = {
    history?: RankHistoryPointDto[]
}
