import { apiRequest } from '@/shared/api'

import type { Leaderboard } from '../model/types'
import type { LeaderboardResponseDto } from './types'

type GetLeaderboardParams = {
    limit?: number
    offset?: number
}

export async function getLeaderboard({
    limit = 50,
    offset = 0,
}: GetLeaderboardParams = {}): Promise<Leaderboard> {
    const params = new URLSearchParams({
        limit: String(limit),
        offset: String(offset),
    })

    const response = await apiRequest<LeaderboardResponseDto>(`/leaderboard?${params.toString()}`)

    return {
        entries: response.entries.map((entry) => ({
            rank: entry.rank,
            player: entry.player,
            fragments: entry.fragments,
            score: entry.score,
            rankChange: entry.rank_change,
        })),
    }
}
