import { apiRequest } from '@/shared/api'

import type { CategoriesProgress, PuzzleProgress, RankHistory, RoleProgress } from '../model/types'
import type {
    CategoriesProgressResponseDto,
    PuzzleProgressResponseDto,
    RankHistoryResponseDto,
    RoleProgressDto,
} from './types'

export async function getRoleProgress(): Promise<RoleProgress[]> {
    const response = await apiRequest<RoleProgressDto[]>('/profile/role-progress')

    return response.map((item) => ({
        role: item.role,
        completedCount: item.completed_count ?? 0,
        inProgressCount: item.in_progress_count ?? 0,
        totalStarted: item.total_started ?? 0,
    }))
}

export async function getCategoriesProgress(): Promise<CategoriesProgress> {
    const response = await apiRequest<CategoriesProgressResponseDto>('/profile/categories-progress')

    return {
        totalCompleted: response.total_completed ?? 0,
        stats: (response.stats ?? []).map((item) => ({
            category: item.category ?? '',
            count: item.count ?? 0,
        })),
    }
}

export async function getPuzzleProgress(): Promise<PuzzleProgress> {
    const response = await apiRequest<PuzzleProgressResponseDto>('/profile/puzzle')

    return {
        earnedCount: response.earned_count ?? 0,
        totalCount: response.total_count ?? 0,
        fragments: (response.fragments ?? []).map((fragment) => ({
            fragmentId: fragment.fragment_id ?? '',
            earnedAt: fragment.earned_at ?? '',
        })),
    }
}

export async function getRankHistory(): Promise<RankHistory> {
    const response = await apiRequest<RankHistoryResponseDto>('/profile/rank-history')

    return {
        history: (response.history ?? []).map((point) => ({
            date: point.date ?? '',
            rank: point.rank ?? 0,
        })),
    }
}
