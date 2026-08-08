// entities/training-path/model/types.ts

export type TrainingScenario = {
    id: string
    title: string
    description: string
    durationMinutes?: number
    isCompleted: boolean
}

export type TrainingScheme = {
    id: string
    title: string
    description: string
    scenarios: TrainingScenario[]
}
