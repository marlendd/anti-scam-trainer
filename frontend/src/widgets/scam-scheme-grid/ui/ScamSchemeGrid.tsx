import styles from './ScamSchemeGrid.module.scss'
import { ScamSchemeCard } from '@/entities/scam-scheme-card'

export const scamSchemes = [
    {
        id: 'phishing-link',
        title: 'Фишинговая ссылка',
        description:
            'Мошенник отправляет ссылку на поддельную страницу оплаты или доставки, чтобы получить данные банковской карты.',
    },
    {
        id: 'advance-payment',
        title: 'Требование предоплаты',
        description:
            'Продавец просит заранее перевести всю сумму или внести залог, после чего перестаёт выходить на связь.',
    },
    {
        id: 'fake-delivery',
        title: 'Поддельная доставка',
        description:
            'Мошенник предлагает оформить доставку через сторонний сайт и присылает ссылку на фальшивую форму.',
    },
    {
        id: 'sms-code-request',
        title: 'Запрос кода из СМС',
        description:
            'Собеседник просит назвать код подтверждения, якобы необходимый для получения оплаты или оформления сделки.',
    },
    {
        id: 'fake-support',
        title: 'Ложная служба поддержки',
        description:
            'Мошенник представляется сотрудником поддержки и пытается получить данные карты, пароль или код подтверждения.',
    },
    {
        id: 'external-messenger',
        title: 'Переход в другой мессенджер',
        description:
            'Собеседник настойчиво предлагает продолжить общение вне платформы, где сложнее проверить его действия.',
    },
    {
        id: 'artificial-urgency',
        title: 'Искусственная срочность',
        description:
            'Мошенник торопит с оплатой или передачей данных, не оставляя времени спокойно проверить предложение.',
    },
    {
        id: 'excess-payment',
        title: 'Лишний перевод',
        description:
            'Мошенник утверждает, что случайно перевёл больше денег, и просит вернуть разницу на другую карту.',
    },
    {
        id: 'product-substitution',
        title: 'Подмена товара',
        description:
            'Покупатель возвращает другой, повреждённый или более дешёвый товар, утверждая, что именно его получил.',
    },
]

export function ScamSchemeGrid() {
    return (
        <section className={styles.grid}>
            {scamSchemes.map((scheme) => (
                <ScamSchemeCard
                    key={scheme.id}
                    id={scheme.id}
                    title={scheme.title}
                    description={scheme.description}
                />
            ))}
        </section>
    )
}
