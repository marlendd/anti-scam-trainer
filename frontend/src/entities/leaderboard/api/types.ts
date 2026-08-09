export type LeaderboardEntryDto = {
    rank: number
    player: string
    fragments: number
    score: number
    rank_change: number | null
}

export type LeaderboardResponseDto = {
    entries: LeaderboardEntryDto[]
}
