// src/widgets/dashboard-path-progress/model/pathProgress.ts

export type TrainingPathId = 'buyer' | 'seller'

export type TrainingPathProgress = {
    id: TrainingPathId
    title: string
    description: string
    completedScenarios: number
    totalScenarios: number
    color: string
}

export const pathProgressMock: TrainingPathProgress[] = [
    {
        id: 'buyer',
        title: 'Путь покупателя',
        description: 'Безопасная покупка, оплата и получение товара',
        completedScenarios: 7,
        totalScenarios: 10,
        color: '#00aaff',
    },
    {
        id: 'seller',
        title: 'Путь продавца',
        description: 'Безопасная продажа, доставка и получение оплаты',
        completedScenarios: 4,
        totalScenarios: 10,
        color: '#ff4053',
    },
]
