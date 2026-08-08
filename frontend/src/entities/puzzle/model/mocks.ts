import safeDealImage from '@/shared/assets/images/puzzle.webp'

import type { PuzzleCollection } from './types'

export const puzzleCollectionsMock: PuzzleCollection[] = [
    {
        id: 'safe-deal',
        title: 'Безопасная сделка',
        description: 'Соберите пазл, проходя сценарии безопасных сделок.',

        imageSrc: safeDealImage,

        unlockedPieces: [1, 2, 4, 5, 7],

        reward: {
            title: '−30% на продвижение',
            description: 'Скидка на одно продвижение объявления.',
        },
    },

    {
        id: 'delivery',
        title: 'Авито Доставка',
        description: 'Разберитесь в схемах мошенничества с доставкой.',

        imageSrc: safeDealImage,

        unlockedPieces: [1, 2, 3, 4, 5, 6, 7, 8, 9],

        reward: {
            title: 'Бесплатное поднятие',
            description: 'Одно бесплатное поднятие объявления.',
        },
    },

    {
        id: 'phishing',
        title: 'Защита от фишинга',
        description: 'Научитесь отличать настоящие страницы от поддельных.',

        imageSrc: safeDealImage,

        unlockedPieces: [1, 4],

        reward: {
            title: '−20% на доставку',
            description: 'Скидка на следующую Авито Доставку.',
        },
    },
]
