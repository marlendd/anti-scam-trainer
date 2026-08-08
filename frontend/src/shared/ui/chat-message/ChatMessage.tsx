import {
    faCheck,
    faCheckDouble,
    faCircleExclamation,
    faClock,
} from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import styles from './ChatMessage.module.scss'

export type ChatMessageDirection = 'incoming' | 'outgoing'

export type ChatMessageStatus = 'sending' | 'sent' | 'delivered' | 'read' | 'failed'

export type ChatMessagePosition = 'single' | 'first' | 'middle' | 'last'

type ChatMessageProps = {
    text: string
    direction: ChatMessageDirection
    position?: ChatMessagePosition
    avatarText?: string
    showAvatar?: boolean

    time?: string
    status?: ChatMessageStatus
}

const statusLabels: Record<ChatMessageStatus, string> = {
    sending: 'Отправляется',
    sent: 'Отправлено',
    delivered: 'Доставлено',
    read: 'Прочитано',
    failed: 'Ошибка отправки',
}

function MessageStatus({ status }: { status: ChatMessageStatus }) {
    const icon =
        status === 'sending'
            ? faClock
            : status === 'sent'
              ? faCheck
              : status === 'failed'
                ? faCircleExclamation
                : faCheckDouble

    return (
        <span
            className={styles.status}
            data-status={status}
            aria-label={statusLabels[status]}
            title={statusLabels[status]}
        >
            <Icon icon={icon} />
        </span>
    )
}

export function ChatMessage({
    text,
    time,
    direction,
    status,
    position = 'single',
    avatarText,
    showAvatar = false,
    // className,
}: ChatMessageProps) {
    const rootClassName = [
        styles.message,
        styles[direction],
        // className,
    ]
        .filter(Boolean)
        .join(' ')

    const meta = (
        <div className={styles.meta}>
            {direction === 'outgoing' && status && <MessageStatus status={status} />}

            <time className={styles.time}>{time}</time>
        </div>
    )

    return (
        <div className={rootClassName}>
            {direction === 'incoming' && (
                <div className={styles.avatarSlot}>
                    {showAvatar && (
                        <div className={styles.avatar} aria-hidden="true">
                            {avatarText?.slice(0, 2).toUpperCase() ?? '?'}
                        </div>
                    )}
                </div>
            )}

            <div className={styles.messageLine}>
                {direction === 'outgoing' && meta}

                <div className={styles.bubble} data-direction={direction} data-position={position}>
                    <p className={styles.text}>{text}</p>
                </div>

                {direction === 'incoming' && meta}
            </div>
        </div>
    )
}
