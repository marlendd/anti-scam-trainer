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
    const isCompleted = attempt.status === 'completed'

    useEffect(() => {
        setHistory([])
    }, [attempt.id])

    useEffect(() => {
        if (!node) {
            return
        }

        setHistory((currentHistory) => {
            const messageId = `node-${node.id}`

            const alreadyExists = currentHistory.some(
                (message) => message.id === messageId,
            )

            if (alreadyExists) {
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

        const messageId = `choice-${node.id}-${choiceId}`

        setHistory((currentHistory) => [
            ...currentHistory,
            {
                id: messageId,
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
        } catch {
            setHistory((currentHistory) =>
                currentHistory.filter(
                    (message) => message.id !== messageId,
                ),
            )
        }
    }

    if (isCompleted) {
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
                        Не удалось отправить ответ. Попробуйте ещё раз.
                    </span>
                )}
            </div>

            <div ref={messagesEndRef} />
        </div>
    )
}