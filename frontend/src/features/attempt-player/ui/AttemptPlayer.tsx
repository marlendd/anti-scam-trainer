import { useEffect, useRef, useState } from 'react'

import {
    type TrainingAttemptState,
    useSubmitAnswer,
} from '@/entities/training-attempt'
import { ChatMessage } from '@/shared/ui/chat-message'

import styles from './AttemptPlayer.module.scss'

type AttemptPlayerProps = {
    attempt: TrainingAttemptState
    onComplete?: () => void
}

type HistoryMessage = {
    id: string
    direction: 'incoming' | 'outgoing'
    text: string
    avatarText?: string
}

export function AttemptPlayer({
    attempt,
    onComplete,
}: AttemptPlayerProps) {
    const messagesEndRef = useRef<HTMLDivElement>(null)

    const [history, setHistory] = useState<HistoryMessage[]>([])

    const submitAnswer = useSubmitAnswer()

    const node = attempt.currentNode

    useEffect(() => {
        setHistory([])
    }, [attempt.id])

    useEffect(() => {
        if (!node) {
            return
        }

        setHistory((currentHistory) => {
            const messageId = `node-${node.id}`

            if (
                currentHistory.some(
                    (message) => message.id === messageId,
                )
            ) {
                return currentHistory
            }

            return [
                ...currentHistory,
                {
                    id: messageId,
                    direction: 'incoming',
                    text: node.text,
                    avatarText: node.author,
                },
            ]
        })
    }, [node])

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({
            behavior: 'smooth',
            block: 'end',
        })
    }, [history.length])

    async function handleChoiceSelect(
        choiceId: string,
        choiceText: string,
    ) {
        if (!node || submitAnswer.isPending) {
            return
        }

        setHistory((currentHistory) => [
            ...currentHistory,
            {
                id: `choice-${node.id}-${choiceId}`,
                direction: 'outgoing',
                text: choiceText,
            },
        ])

        try {
            const result = await submitAnswer.mutateAsync({
                attemptId: attempt.id,
                nodeId: node.id,
                choiceId,
            })

            if (result.completed) {
                onComplete?.()
            }
        } catch (error) {
            setHistory((currentHistory) =>
                currentHistory.filter(
                    (message) =>
                        message.id !==
                        `choice-${node.id}-${choiceId}`,
                ),
            )

            throw error
        }
    }

    return (
        <div className={styles.messages}>
            {history.map((message) => (
                <ChatMessage
                    key={message.id}
                    direction={message.direction}
                    position="single"
                    text={message.text}
                    avatarText={message.avatarText}
                    showAvatar={message.direction === 'incoming'}
                />
            ))}

            {attempt.status === 'in_progress' && node && (
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
                                    choice.text,
                                )
                            }}
                        >
                            {choice.text}
                        </button>
                    ))}

                    {submitAnswer.isError && (
                        <span className={styles.error}>
                            Не удалось отправить ответ.
                        </span>
                    )}
                </div>
            )}

            {attempt.status === 'completed' && (
                <div className={styles.completed}>
                    <strong>Сценарий завершён</strong>

                    {attempt.score !== undefined && (
                        <span>
                            Результат: {attempt.score}
                        </span>
                    )}
                </div>
            )}

            <div ref={messagesEndRef} />
        </div>
    )
}