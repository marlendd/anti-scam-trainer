import type { TrainingScenario } from './types'

export const phishingLinkScenario: TrainingScenario = {
    id: 'phishing-link',
    title: 'Фишинговая ссылка',
    description: 'Поддельная страница оплаты под видом официального сервиса.',

    product: {
        id: 'matebook',
        title: 'Huawei MateBook 14s',
        price: 59990,
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
            name: 'Илья',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'phishing-message-1',
            senderId: 'buyer',
            text: 'Здравствуйте! Ноутбук ещё продаётся?',
            delayMs: 500,
        },
        {
            type: 'choice',
            id: 'phishing-choice-1',
            previewOptionId: 'phishing-choice-1-answer-1',
            options: [
                {
                    id: 'phishing-choice-1-answer-1',
                    text: 'Здравствуйте. Да, объявление актуально.',
                    isCorrect: true,
                },
                {
                    id: 'phishing-choice-1-answer-2',
                    text: 'Да, можете оплачивать.',
                    isCorrect: false,
                },
            ],
        },
        {
            type: 'message',
            id: 'phishing-message-2',
            senderId: 'buyer',
            text: 'Я уже оформил доставку. Вам осталось получить оплату по ссылке.',
            delayMs: 800,
        },
        {
            type: 'message',
            id: 'phishing-message-3',
            senderId: 'buyer',
            text: 'https://avito-secure-pay.example/order/18492',
            delayMs: 500,
        },
        {
            type: 'choice',
            id: 'phishing-choice-2',
            previewOptionId: 'phishing-choice-2-answer-2',
            options: [
                {
                    id: 'phishing-choice-2-answer-1',
                    text: 'Хорошо, сейчас открою.',
                    isCorrect: false,
                },
                {
                    id: 'phishing-choice-2-answer-2',
                    text: 'Оплату я проверю только внутри приложения.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Фишинговая ссылка',

        redFlags: [
            {
                id: 'external-link',
                title: 'Внешняя ссылка',
                description:
                    'Собеседник предлагает перейти на сторонний сайт для получения денег.',
            },
            {
                id: 'fake-domain',
                title: 'Поддельный адрес',
                description:
                    'Домен лишь имитирует официальный сервис и не относится к платформе.',
            },
            {
                id: 'fake-payment',
                title: 'Оплата вне приложения',
                description:
                    'Получать деньги или подтверждать доставку через сторонние формы не требуется.',
            },
        ],

        safeActions: [
            'Не переходить по ссылкам из сообщений',
            'Проверять заказ внутри приложения',
            'Не вводить данные банковской карты на сторонних сайтах',
        ],

        goldenRule:
            'Все действия с оплатой и доставкой выполняйте только внутри официального приложения.',
    },
}

export const advancePaymentScenario: TrainingScenario = {
    id: 'advance-payment',
    title: 'Требование предоплаты',
    description: 'Продавец просит заранее перевести деньги или внести залог.',

    product: {
        id: 'iphone',
        title: 'iPhone 15 Pro 256 ГБ',
        price: 84990,
    },

    playerParticipantId: 'buyer',

    participants: [
        {
            id: 'buyer',
            name: 'Александр',
            role: 'buyer',
        },
        {
            id: 'seller',
            name: 'Максим',
            role: 'seller',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'advance-message-1',
            senderId: 'buyer',
            text: 'Здравствуйте. Телефон ещё в продаже?',
            delayMs: 400,
        },
        {
            type: 'message',
            id: 'advance-message-2',
            senderId: 'seller',
            text: 'Да. Но желающих много, могу придержать за вами.',
            delayMs: 700,
        },
        {
            type: 'message',
            id: 'advance-message-3',
            senderId: 'seller',
            text: 'Переведите 5000 ₽ залога на карту, и я сниму объявление.',
            delayMs: 700,
        },
        {
            type: 'choice',
            id: 'advance-choice-1',
            previewOptionId: 'advance-choice-1-answer-2',
            options: [
                {
                    id: 'advance-choice-1-answer-1',
                    text: 'Хорошо, пришлите номер карты.',
                    isCorrect: false,
                },
                {
                    id: 'advance-choice-1-answer-2',
                    text: 'Предоплату переводить не буду. Проведём сделку через платформу.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Требование предоплаты',

        redFlags: [
            {
                id: 'deposit',
                title: 'Залог до сделки',
                description:
                    'Продавец требует деньги до осмотра товара или оформления безопасной сделки.',
            },
            {
                id: 'direct-transfer',
                title: 'Перевод напрямую',
                description:
                    'Деньги предлагают отправить непосредственно на банковскую карту.',
            },
            {
                id: 'scarcity',
                title: 'Давление спросом',
                description:
                    'Фраза о большом количестве желающих используется, чтобы ускорить решение.',
            },
        ],

        safeActions: [
            'Не переводить залог незнакомым продавцам',
            'Использовать встроенные способы оплаты',
            'Не принимать решение под давлением',
        ],

        goldenRule:
            'Не переводите деньги напрямую незнакомому человеку до безопасного оформления сделки.',
    },
}

export const fakeDeliveryScenario: TrainingScenario = {
    id: 'fake-delivery',
    title: 'Поддельная доставка',
    description: 'Собеседник предлагает оформить доставку через сторонний сервис.',

    product: {
        id: 'camera',
        title: 'Фотоаппарат Sony Alpha',
        price: 45000,
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
            name: 'Антон',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'delivery-message-1',
            senderId: 'buyer',
            text: 'Добрый день. Готов купить с доставкой.',
            delayMs: 400,
        },
        {
            type: 'choice',
            id: 'delivery-choice-1',
            previewOptionId: 'delivery-choice-1-answer-1',
            options: [
                {
                    id: 'delivery-choice-1-answer-1',
                    text: 'Хорошо, оформляйте доставку через приложение.',
                    isCorrect: true,
                },
                {
                    id: 'delivery-choice-1-answer-2',
                    text: 'Хорошо. Что нужно сделать?',
                    isCorrect: false,
                },
            ],
        },
        {
            type: 'message',
            id: 'delivery-message-2',
            senderId: 'buyer',
            text: 'У меня приложение сегодня не работает. Оформил через партнёрскую службу.',
            delayMs: 700,
        },
        {
            type: 'message',
            id: 'delivery-message-3',
            senderId: 'buyer',
            text: 'Перейдите на avito-delivery.example и подтвердите получение денег.',
            delayMs: 600,
        },
        {
            type: 'choice',
            id: 'delivery-choice-2',
            previewOptionId: 'delivery-choice-2-answer-2',
            options: [
                {
                    id: 'delivery-choice-2-answer-1',
                    text: 'Хорошо, попробую.',
                    isCorrect: false,
                },
                {
                    id: 'delivery-choice-2-answer-2',
                    text: 'Доставку оформляем только через платформу.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Поддельная доставка',

        redFlags: [
            {
                id: 'partner-service',
                title: 'Неизвестный сервис доставки',
                description:
                    'Покупатель предлагает использовать стороннюю службу вместо встроенной доставки.',
            },
            {
                id: 'external-page',
                title: 'Сторонняя форма',
                description:
                    'Для подтверждения сделки предлагается открыть внешний сайт.',
            },
        ],

        safeActions: [
            'Оформлять доставку только внутри приложения',
            'Не вводить платёжные данные на сторонних страницах',
            'Отказаться от сделки при попытке обойти платформу',
        ],

        goldenRule:
            'Безопасная доставка оформляется и отслеживается внутри платформы.',
    },
}

export const smsCodeRequestScenario: TrainingScenario = {
    id: 'sms-code-request',
    title: 'Запрос кода из СМС',
    description: 'Собеседник пытается получить одноразовый код подтверждения.',

    product: {
        id: 'console',
        title: 'PlayStation 5',
        price: 49990,
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
            name: 'Денис',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'sms-message-1',
            senderId: 'buyer',
            text: 'Я оформил заказ, но система просит подтверждение от продавца.',
            delayMs: 500,
        },
        {
            type: 'message',
            id: 'sms-message-2',
            senderId: 'buyer',
            text: 'Сейчас вам придёт СМС. Скиньте мне код, чтобы завершить оплату.',
            delayMs: 800,
        },
        {
            type: 'choice',
            id: 'sms-choice-1',
            previewOptionId: 'sms-choice-1-answer-2',
            options: [
                {
                    id: 'sms-choice-1-answer-1',
                    text: 'Хорошо, сейчас пришлю.',
                    isCorrect: false,
                },
                {
                    id: 'sms-choice-1-answer-2',
                    text: 'Коды из СМС никому не сообщаю.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Запрос кода из СМС',

        redFlags: [
            {
                id: 'sms-secret',
                title: 'Запрос секретного кода',
                description:
                    'Одноразовые коды используются для подтверждения действий владельца аккаунта.',
            },
            {
                id: 'fake-reason',
                title: 'Выдуманное подтверждение',
                description:
                    'Для получения оплаты покупателю не нужен код из СМС продавца.',
            },
        ],

        safeActions: [
            'Никому не сообщать коды подтверждения',
            'Читать текст СМС перед любым действием',
            'При подозрении прекратить общение',
        ],

        goldenRule:
            'Код из СМС — такой же секрет, как пароль. Его нельзя передавать другим людям.',
    },
}

export const fakeSupportScenario: TrainingScenario = {
    id: 'fake-support',
    title: 'Ложная служба поддержки',
    description: 'Мошенник представляется сотрудником поддержки платформы.',

    product: {
        id: 'headphones',
        title: 'AirPods Pro',
        price: 15990,
    },

    playerParticipantId: 'seller',

    participants: [
        {
            id: 'seller',
            name: 'Александр',
            role: 'seller',
        },
        {
            id: 'support',
            name: 'Служба поддержки',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'support-message-1',
            senderId: 'support',
            text: 'Здравствуйте. Мы обнаружили подозрительную операцию по вашему объявлению.',
            delayMs: 500,
        },
        {
            type: 'message',
            id: 'support-message-2',
            senderId: 'support',
            text: 'Чтобы избежать блокировки, подтвердите номер карты и код из СМС.',
            delayMs: 800,
        },
        {
            type: 'choice',
            id: 'support-choice-1',
            previewOptionId: 'support-choice-1-answer-2',
            options: [
                {
                    id: 'support-choice-1-answer-1',
                    text: 'Какие данные карты вам нужны?',
                    isCorrect: false,
                },
                {
                    id: 'support-choice-1-answer-2',
                    text: 'Я самостоятельно обращусь в поддержку через приложение.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Ложная служба поддержки',

        redFlags: [
            {
                id: 'sensitive-data',
                title: 'Запрос платёжных данных',
                description:
                    'Настоящей поддержке не нужны полный номер карты, CVV или одноразовый код.',
            },
            {
                id: 'account-threat',
                title: 'Угроза блокировкой',
                description:
                    'Мошенник создаёт чувство страха и требует действовать немедленно.',
            },
        ],

        safeActions: [
            'Открыть поддержку самостоятельно через приложение',
            'Не сообщать пароли, коды и данные карты',
            'Проверить наличие официального уведомления',
        ],

        goldenRule:
            'Не доверяйте человеку только потому, что он представился сотрудником поддержки.',
    },
}

export const externalMessengerScenario: TrainingScenario = {
    id: 'external-messenger',
    title: 'Переход в другой мессенджер',
    description: 'Собеседник пытается вывести общение за пределы платформы.',

    product: {
        id: 'bike',
        title: 'Велосипед Trek',
        price: 38000,
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
            name: 'Роман',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'messenger-message-1',
            senderId: 'buyer',
            text: 'Здравствуйте. Можно подробнее про велосипед?',
            delayMs: 400,
        },
        {
            type: 'choice',
            id: 'messenger-choice-1',
            previewOptionId: 'messenger-choice-1-answer-1',
            options: [
                {
                    id: 'messenger-choice-1-answer-1',
                    text: 'Конечно, задавайте вопросы.',
                    isCorrect: true,
                },
                {
                    id: 'messenger-choice-1-answer-2',
                    text: 'Да, напишите куда удобнее.',
                    isCorrect: false,
                },
            ],
        },
        {
            type: 'message',
            id: 'messenger-message-2',
            senderId: 'buyer',
            text: 'Тут неудобно. Давайте в Telegram, мой ник @buyer_example.',
            delayMs: 700,
        },
        {
            type: 'message',
            id: 'messenger-message-3',
            senderId: 'buyer',
            text: 'Там я скину данные для оплаты и доставки.',
            delayMs: 500,
        },
        {
            type: 'choice',
            id: 'messenger-choice-2',
            previewOptionId: 'messenger-choice-2-answer-2',
            options: [
                {
                    id: 'messenger-choice-2-answer-1',
                    text: 'Хорошо, сейчас напишу.',
                    isCorrect: false,
                },
                {
                    id: 'messenger-choice-2-answer-2',
                    text: 'Предпочитаю продолжить общение здесь.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Переход в другой мессенджер',

        redFlags: [
            {
                id: 'leave-platform',
                title: 'Уход с платформы',
                description:
                    'Собеседник без объективной причины настаивает на переходе в сторонний мессенджер.',
            },
            {
                id: 'payment-outside',
                title: 'Оплата вне платформы',
                description:
                    'В другом мессенджере мошеннику проще прислать поддельную ссылку или реквизиты.',
            },
        ],

        safeActions: [
            'Продолжать переписку внутри платформы',
            'Не переходить по присланным внешним ссылкам',
            'Не отправлять контакты без необходимости',
        ],

        goldenRule:
            'Если сделка началась на платформе, безопаснее вести там и переписку, и оформление.',
    },
}

export const artificialUrgencyScenario: TrainingScenario = {
    id: 'artificial-urgency',
    title: 'Искусственная срочность',
    description: 'Мошенник заставляет принять решение, не оставляя времени на проверку.',

    product: {
        id: 'gpu',
        title: 'GeForce RTX 4070',
        price: 48990,
    },

    playerParticipantId: 'buyer',

    participants: [
        {
            id: 'buyer',
            name: 'Александр',
            role: 'buyer',
        },
        {
            id: 'seller',
            name: 'Виктор',
            role: 'seller',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'urgency-message-1',
            senderId: 'seller',
            text: 'Видеокарта осталась последняя. Уже есть другой покупатель.',
            delayMs: 500,
        },
        {
            type: 'message',
            id: 'urgency-message-2',
            senderId: 'seller',
            text: 'Если прямо сейчас переведёте деньги, отдам вам.',
            delayMs: 700,
        },
        {
            type: 'message',
            id: 'urgency-message-3',
            senderId: 'seller',
            text: 'Решайте за пять минут, потом будет поздно.',
            delayMs: 500,
        },
        {
            type: 'choice',
            id: 'urgency-choice-1',
            previewOptionId: 'urgency-choice-1-answer-2',
            options: [
                {
                    id: 'urgency-choice-1-answer-1',
                    text: 'Хорошо, куда переводить?',
                    isCorrect: false,
                },
                {
                    id: 'urgency-choice-1-answer-2',
                    text: 'Сначала я проверю товар и условия сделки.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Искусственная срочность',

        redFlags: [
            {
                id: 'deadline',
                title: 'Искусственный дедлайн',
                description:
                    'Пользователю дают несколько минут на решение без объективной причины.',
            },
            {
                id: 'other-buyers',
                title: 'Давление конкуренцией',
                description:
                    'История о другом покупателе используется для отключения критического мышления.',
            },
        ],

        safeActions: [
            'Не принимать финансовые решения в спешке',
            'Проверять продавца и товар',
            'Отказываться от сделки при сильном давлении',
        ],

        goldenRule:
            'Надёжная сделка не требует принимать решение за несколько минут.',
    },
}

export const excessPaymentScenario: TrainingScenario = {
    id: 'excess-payment',
    title: 'Лишний перевод',
    description: 'Мошенник сообщает о якобы ошибочном переводе и просит вернуть разницу.',

    product: {
        id: 'monitor',
        title: 'Монитор LG UltraGear',
        price: 32000,
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
            name: 'Никита',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'excess-message-1',
            senderId: 'buyer',
            text: 'Я случайно отправил вам 42 000 ₽ вместо 32 000 ₽.',
            delayMs: 500,
        },
        {
            type: 'message',
            id: 'excess-message-2',
            senderId: 'buyer',
            text: 'Верните, пожалуйста, лишние 10 000 ₽ на другую карту. Очень срочно.',
            delayMs: 700,
        },
        {
            type: 'choice',
            id: 'excess-choice-1',
            previewOptionId: 'excess-choice-1-answer-2',
            options: [
                {
                    id: 'excess-choice-1-answer-1',
                    text: 'Хорошо, пришлите номер карты.',
                    isCorrect: false,
                },
                {
                    id: 'excess-choice-1-answer-2',
                    text: 'Сначала проверю фактическое поступление денег и обращусь в поддержку.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Лишний перевод',

        redFlags: [
            {
                id: 'unverified-payment',
                title: 'Непроверенный перевод',
                description:
                    'Слова собеседника или скриншот не подтверждают фактическое поступление денег.',
            },
            {
                id: 'different-card',
                title: 'Возврат на другую карту',
                description:
                    'Просьба отправить деньги по другим реквизитам особенно подозрительна.',
            },
            {
                id: 'urgency',
                title: 'Срочность',
                description:
                    'Мошенник старается добиться перевода до проверки операции.',
            },
        ],

        safeActions: [
            'Проверять поступление только в банковском приложении',
            'Не возвращать деньги на сторонние реквизиты',
            'При спорной операции обращаться в банк или поддержку',
        ],

        goldenRule:
            'Никогда не возвращайте якобы лишний перевод, пока самостоятельно не подтвердили операцию.',
    },
}

export const productSubstitutionScenario: TrainingScenario = {
    id: 'product-substitution',
    title: 'Подмена товара',
    description: 'Покупатель пытается вернуть другой или повреждённый товар.',

    product: {
        id: 'keyboard',
        title: 'Клавиатура Keychron Q1',
        price: 14990,
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
            name: 'Сергей',
            role: 'buyer',
            status: 'в сети',
        },
    ],

    timeline: [
        {
            type: 'message',
            id: 'substitution-message-1',
            senderId: 'buyer',
            text: 'Здравствуйте. Получил клавиатуру, она вся в царапинах и не включается.',
            delayMs: 500,
        },
        {
            type: 'message',
            id: 'substitution-message-2',
            senderId: 'buyer',
            text: 'Хочу вернуть её и получить деньги обратно.',
            delayMs: 600,
        },
        {
            type: 'choice',
            id: 'substitution-choice-1',
            previewOptionId: 'substitution-choice-1-answer-2',
            options: [
                {
                    id: 'substitution-choice-1-answer-1',
                    text: 'Хорошо, отправляйте обратно.',
                    isCorrect: false,
                },
                {
                    id: 'substitution-choice-1-answer-2',
                    text: 'Давайте оформим возврат через платформу и проверим товар.',
                    isCorrect: true,
                },
            ],
        },
        {
            type: 'message',
            id: 'substitution-message-3',
            senderId: 'buyer',
            text: 'Не хочу оформлять официально. Просто верните деньги, а я отправлю посылку потом.',
            delayMs: 700,
        },
        {
            type: 'choice',
            id: 'substitution-choice-2',
            previewOptionId: 'substitution-choice-2-answer-2',
            options: [
                {
                    id: 'substitution-choice-2-answer-1',
                    text: 'Ладно, так быстрее.',
                    isCorrect: false,
                },
                {
                    id: 'substitution-choice-2-answer-2',
                    text: 'Возврат проведём только через предусмотренную процедуру.',
                    isCorrect: true,
                },
            ],
        },
    ],

    analysis: {
        title: 'Подмена товара',

        redFlags: [
            {
                id: 'unofficial-return',
                title: 'Возврат вне платформы',
                description:
                    'Покупатель пытается обойти официальный процесс возврата.',
            },
            {
                id: 'money-first',
                title: 'Сначала деньги',
                description:
                    'Покупатель требует возврат средств до проверки возвращённого товара.',
            },
        ],

        safeActions: [
            'Фиксировать состояние товара перед отправкой',
            'Проводить возврат через официальную процедуру',
            'Проверять полученный обратно товар до завершения возврата',
        ],

        goldenRule:
            'При возврате важно сохранить возможность подтвердить, какой именно товар был отправлен и возвращён.',
    },
}

export const trainingScenarioMocks: TrainingScenario[] = [
    phishingLinkScenario,
    advancePaymentScenario,
    fakeDeliveryScenario,
    smsCodeRequestScenario,
    fakeSupportScenario,
    externalMessengerScenario,
    artificialUrgencyScenario,
    excessPaymentScenario,
    productSubstitutionScenario,
]

export const trainingScenarioById = Object.fromEntries(
    trainingScenarioMocks.map((scenario) => [
        scenario.id,
        scenario,
    ]),
) as Record<string, TrainingScenario>