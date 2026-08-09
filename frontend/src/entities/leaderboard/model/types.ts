export type LeaderboardEntry = {
    rank: number
    player: string
    fragments: number
    score: number
    rankChange: number | null
}

export type Leaderboard = {
    entries: LeaderboardEntry[]
}
