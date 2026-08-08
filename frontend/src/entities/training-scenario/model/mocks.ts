import type { TrainingScenario } from './types'

export const fakeDeliveryScenario: TrainingScenario = {
    id: 'fishing-link',

    title: 'Фишинговая ссылка',
    description: 'Сценарий с поддельной страницей оплаты',

    product: {
        id: 'huawei-matebook-parts',
        title: 'Запчасти для Huawei MateBook 14s HKF-X',
        imageUrl:
            'https://20.img.avito.st/image/1/1.PbRgtLaxkV1OFUFeFugMrxcWk13SAZVf.c2Kr5fPnMZpzko45kbOYI-36E9lMfvOeCcWRKDfVkIQ?cqp=2.jyhhz7zK6lMkGfOGBa4Nrold2Hxx--U_poa2YsxOn6Tkc-5hBhg-gigekeJxvEbQ',
        price: 9999,
    },

    playerParticipantId: 'seller',

    participants: [
        {
            id: 'seller',
            name: 'Александр',
            role: 'seller',
        },
        {
            id: 'buyer',
            name: 'Это разборка Питерская',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'message-1',
            senderId: 'buyer',
            text: 'Здравствуйте! Товар ещё продаётся?',
            delayMs: 500,
        },

        {
            type: 'choice',
            id: 'choice-1',

            previewOptionId: 'choice-1-answer-1',

            options: [
                {
                    id: 'choice-1-answer-1',
                    text: 'Здравствуйте. Да, ещё актуально.',
                    isCorrect: true,
                },
                {
                    id: 'choice-1-answer-2',
                    text: 'Да. Оплачивайте сразу.',
                    isCorrect: false,
                },
                {
                    id: 'choice-1-answer-3',
                    text: 'Да, напишите мне в Telegram.',
                    isCorrect: false,
                },
            ],
        },

        {
            type: 'message',
            id: 'message-2',
            senderId: 'buyer',
            text: 'Отлично. Я уже оплатил доставку. Перейдите по ссылке и получите деньги.',
            delayMs: 900,
        },

        {
            type: 'message',
            id: 'message-3',
            senderId: 'buyer',
            text: 'https://avito-pay-confirm.ru/order/18492',
            delayMs: 600,
        },

        {
            type: 'choice',
            id: 'choice-2',

            previewOptionId: 'choice-2-answer-2',

            options: [
                {
                    id: 'choice-2-answer-1',
                    text: 'Хорошо, сейчас перейду.',
                    isCorrect: false,
                    feedback: 'Не переходите по ссылкам для получения оплаты.',
                },
                {
                    id: 'choice-2-answer-2',
                    text: 'Я проверю оплату только внутри приложения.',
                    isCorrect: true,
                    feedback: 'Верно. Все операции необходимо проверять внутри приложения.',
                },
                {
                    id: 'choice-2-answer-3',
                    text: 'Пришлите ссылку ещё раз, эта не открывается.',
                    isCorrect: false,
                },
            ],
        },
    ],

    analysis: {
        title: 'Фишинговая ссылка',

        redFlags: [
            {
                id: 'urgency',
                title: 'Красный флаг 1',
                description:
                    'Покупатель заявляет, что уже совершил операцию без согласования деталей, и требует немедленных действий.',
            },
            {
                id: 'external-link',
                title: 'Красный флаг 2',
                description:
                    'Собеседник отправляет ссылку для получения денег. Все операции должны проходить внутри приложения.',
            },
            {
                id: 'phishing-domain',
                title: 'Красный флаг 3',
                description:
                    'Адрес avito-pay-confirm.ru имитирует официальный сервис и используется для кражи данных карты.',
                accent: 'avito-pay-confirm.ru',
            },
        ],

        safeActions: [
            'Не переходить по ссылкам на оплату из сообщений',
            'Проверять статус заказа только внутри приложения',
            'Не вводить данные банковской карты на сторонних сайтах',
        ],

        goldenRule:
            'Все операции с деньгами и доставкой проходят внутри приложения. Любая внешняя ссылка — угроза.',
    },
}

export const trainingScenarioMocks: TrainingScenario[] = [fakeDeliveryScenario]
