export type ProgressRole = 'buyer' | 'seller'

export type RoleProgress = {
    role: ProgressRole
    completedCount: number
    inProgressCount: number
    totalStarted: number
}

export type CategoryProgress = {
    category: string
    count: number
}

export type CategoriesProgress = {
    totalCompleted: number
    stats: CategoryProgress[]
}

export type PuzzleFragment = {
    fragmentId: string
    earnedAt: string
}

export type PuzzleProgress = {
    earnedCount: number
    totalCount: number
    fragments: PuzzleFragment[]
}

export type RankHistoryPoint = {
    date: string
    rank: number
}

export type RankHistory = {
    history: RankHistoryPoint[]
}
