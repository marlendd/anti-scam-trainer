import { useMemo, useRef, useEffect } from 'react'

import {
    type AttemptActor,
    type TrainingAttemptState,
    useSubmitAnswer,
} from '@/entities/training-attempt'
import { ChatMessage } from '@/shared/ui/chat-message'

import styles from './AttemptPlayer.module.scss'

type AttemptPlayerProps = {
    attempt: TrainingAttemptState
    onComplete?: () => void
}

type ChatEntry = {
    id: string
    author: AttemptActor
    text: string
}

function getMessagePosition(
    messages: ChatEntry[],
    index: number,
) {
    const current = messages[index]
    const previous = messages[index - 1]
    const next = messages[index + 1]

    const sameAsPrevious =
        previous?.author === current?.author

    const sameAsNext =
        next?.author === current?.author

    if (!sameAsPrevious && !sameAsNext) {
        return 'single' as const
    }

    if (!sameAsPrevious && sameAsNext) {
        return 'first' as const
    }

    if (sameAsPrevious && sameAsNext) {
        return 'middle' as const
    }

    return 'last' as const
}

function getAuthorName(author: AttemptActor) {
    return author === 'buyer'
        ? 'Покупатель'
        : 'Продавец'
}

export function AttemptPlayer({
    attempt,
    onComplete,
}: AttemptPlayerProps) {
    const messagesEndRef = useRef<HTMLDivElement>(null)

    const submitAnswer = useSubmitAnswer()

    const node = attempt.currentNode
    const playerRole = attempt.scenario.role

    const messages = useMemo<ChatEntry[]>(() => {
        const historyMessages =
            attempt.history.flatMap((historyItem) => [
                ...historyItem.node.messages.map(
                    (message, index) => ({
                        id: `${historyItem.node.id}-message-${index}`,
                        author: message.author,
                        text: message.text,
                    }),
                ),

                {
                    id: `${historyItem.node.id}-choice-${historyItem.selectedChoice.id}`,
                    author: playerRole,
                    text: historyItem.selectedChoice.text,
                },
            ])

        const currentMessages =
            node?.messages.map((message, index) => ({
                id: `${node.id}-message-${index}`,
                author: message.author,
                text: message.text,
            })) ?? []

        return [
            ...historyMessages,
            ...currentMessages,
        ]
    }, [attempt.history, node, playerRole])

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({
            behavior: 'smooth',
            block: 'end',
        })
    }, [messages.length])

    async function handleChoiceSelect(
        choiceId: string,
    ) {
        if (!node || submitAnswer.isPending) {
            return
        }

        const result = await submitAnswer.mutateAsync({
            attemptId: attempt.id,
            nodeId: node.id,
            choiceId,
        })

        if (result.completed) {
            onComplete?.()
        }
    }

    if (attempt.status === 'completed') {
        return null
    }

    if (!node) {
        return (
            <div className={styles.messages}>
                <div className={styles.error}>
                    Не удалось получить текущий шаг сценария.
                </div>
            </div>
        )
    }

    return (
        <div className={styles.messages}>
            {messages.map((message, index) => {
                const direction =
                    message.author === playerRole
                        ? 'outgoing'
                        : 'incoming'

                const position = getMessagePosition(
                    messages,
                    index,
                )

                const showAvatar =
                    direction === 'incoming' &&
                    (position === 'single' ||
                        position === 'last')

                return (
                    <ChatMessage
                        key={message.id}
                        direction={direction}
                        position={position}
                        text={message.text}
                        avatarText={getAuthorName(
                            message.author,
                        )}
                        showAvatar={showAvatar}
                    />
                )
            })}

            <div className={styles.answers}>
                <span className={styles.answersLabel}>
                    Выберите ответ
                </span>

                {node.choices.map((choice) => (
                    <button
                        key={choice.id}
                        type="button"
                        className={styles.answer}
                        disabled={submitAnswer.isPending}
                        onClick={() => {
                            void handleChoiceSelect(
                                choice.id,
                            )
                        }}
                    >
                        {choice.text}
                    </button>
                ))}

                {submitAnswer.isError && (
                    <span className={styles.error}>
                        Не удалось отправить ответ.
                        Попробуйте ещё раз.
                    </span>
                )}
            </div>

            <div ref={messagesEndRef} />
        </div>
    )
}