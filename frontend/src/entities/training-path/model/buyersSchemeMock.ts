import type { TrainingScheme } from './types'

export const buyerSchemesMock: TrainingScheme[] = [
    {
        id: 'fake-delivery',
        title: 'Поддельная доставка',
        description: 'Мошенники имитируют оформление доставки и получение оплаты.',
        scenarios: [
            {
                id: 'fake-delivery-link',
                title: 'Ссылка от продавца',
                description: 'Проверьте подозрительную ссылку на доставку.',
                durationMinutes: 3,
                isCompleted: true,
            },
            {
                id: 'fake-payment-page',
                title: 'Страница получения денег',
                description: 'Определите поддельную форму оплаты.',
                durationMinutes: 4,
                isCompleted: false,
            },
            {
                id: 'bank-card-details',
                title: 'Запрос банковских данных',
                description: 'Не передавайте данные карты мошеннику.',
                durationMinutes: 4,
                isCompleted: false,
            },
        ],
    },
    {
        id: 'fake-support',
        title: 'Поддельная поддержка',
        description: 'Мошенник представляется сотрудником службы поддержки.',
        scenarios: [
            {
                id: 'support-chat',
                title: 'Сообщение от поддержки',
                description: 'Проверьте личность собеседника.',
                durationMinutes: 3,
                isCompleted: false,
            },
            {
                id: 'sms-code',
                title: 'Запрос кода из СМС',
                description: 'Распознайте попытку получить код подтверждения.',
                durationMinutes: 4,
                isCompleted: false,
            },
            {
                id: 'support-chat',
                title: 'Сообщение от поддержки',
                description: 'Проверьте личность собеседника.',
                durationMinutes: 3,
                isCompleted: false,
            },
            {
                id: 'sms-code',
                title: 'Запрос кода из СМС',
                description: 'Распознайте попытку получить код подтверждения.',
                durationMinutes: 4,
                isCompleted: false,
            },
            {
                id: 'support-chat',
                title: 'Сообщение от поддержки',
                description: 'Проверьте личность собеседника.',
                durationMinutes: 3,
                isCompleted: false,
            },
            {
                id: 'sms-code',
                title: 'Запрос кода из СМС',
                description: 'Распознайте попытку получить код подтверждения.',
                durationMinutes: 4,
                isCompleted: false,
            },
            {
                id: 'support-chat',
                title: 'Сообщение от поддержки',
                description: 'Проверьте личность собеседника.',
                durationMinutes: 3,
                isCompleted: false,
            },
            {
                id: 'sms-code',
                title: 'Запрос кода из СМС',
                description: 'Распознайте попытку получить код подтверждения.',
                durationMinutes: 4,
                isCompleted: false,
            },
            {
                id: 'support-chat',
                title: 'Сообщение от поддержки',
                description: 'Проверьте личность собеседника.',
                durationMinutes: 3,
                isCompleted: false,
            },
            {
                id: 'sms-code',
                title: 'Запрос кода из СМС',
                description: 'Распознайте попытку получить код подтверждения.',
                durationMinutes: 4,
                isCompleted: false,
            },
        ],
    },
]
