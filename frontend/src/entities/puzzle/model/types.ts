export type PuzzlePieceId =
  | 1
  | 2
  | 3
  | 4
  | 5
  | 6
  | 7
  | 8
  | 9

export type PuzzleReward = {
  title: string
  description?: string
}

export type PuzzleCollection = {
  id: string
  title: string
  description: string

  imageSrc: string

  unlockedPieces: PuzzlePieceId[]

  reward: PuzzleReward

  rewardClaimed?: boolean
}