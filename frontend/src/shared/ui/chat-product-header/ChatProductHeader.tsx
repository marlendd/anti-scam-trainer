// ChatProductHeader.tsx

import {
    faRotateRight,
} from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import styles from './ChatProductHeader.module.scss'

type ChatProductHeaderProps = {
    productTitle: string
    productPrice: number
    participantName: string
    participantStatus?: string
    imageUrl?: string
    className?: string

    onRestart?: () => void
    isRestarting?: boolean
}

const priceFormatter = new Intl.NumberFormat(
    'ru-RU',
)

export function ChatProductHeader({
    productTitle,
    productPrice,
    participantName,
    participantStatus,
    imageUrl,
    className,
    onRestart,
    isRestarting = false,
}: ChatProductHeaderProps) {
    const rootClassName = [
        styles.header,
        className,
    ]
        .filter(Boolean)
        .join(' ')

    return (
        <header className={rootClassName}>
            {imageUrl && (
                <img
                    className={styles.image}
                    src={imageUrl}
                    alt=""
                />
            )}

            <div className={styles.content}>
                <div className={styles.participant}>
                    <span
                        className={
                            styles.participantName
                        }
                    >
                        {participantName}
                    </span>

                    {participantStatus && (
                        <span
                            className={
                                styles.status
                            }
                        >
                            {participantStatus}
                        </span>
                    )}
                </div>

                <div className={styles.product}>
                    <span
                        className={
                            styles.productTitle
                        }
                    >
                        {productTitle}
                    </span>

                    <span
                        className={styles.price}
                    >
                        {priceFormatter.format(
                            productPrice,
                        )}{' '}
                        ₽
                    </span>
                </div>
            </div>

            {onRestart && (
                <button
                    type="button"
                    className={
                        styles.restartButton
                    }
                    disabled={isRestarting}
                    onClick={onRestart}
                >
                    <Icon icon={faRotateRight} />

                    <span>
                        {isRestarting
                            ? 'Перезапускаем...'
                            : 'Начать заново'}
                    </span>
                </button>
            )}
        </header>
    )
}