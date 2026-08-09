import { useQuery } from '@tanstack/react-query'

import { getLeaderboard } from './get-leaderboard'

type UseLeaderboardParams = {
    limit?: number
    offset?: number
}

export function useLeaderboard({ limit = 50, offset = 0 }: UseLeaderboardParams = {}) {
    return useQuery({
        queryKey: ['leaderboard', limit, offset],
        queryFn: () =>
            getLeaderboard({
                limit,
                offset,
            }),
    })
}
