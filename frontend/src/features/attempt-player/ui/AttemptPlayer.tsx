import {
    useEffect,
    useMemo,
    useRef,
    useState,
} from 'react'

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

type MessageEntry = {
    type: 'message'
    id: string
    author: AttemptActor
    text: string
}

type ConsequenceEntry = {
    type: 'consequence'
    id: string
    text: string
}

type ChatEntry =
    | MessageEntry
    | ConsequenceEntry

type SubmittedAnswer = {
    nodeId: string
    choiceId: string
    choiceText: string
    consequence: string
}

function hashString(value: string) {
    let hash = 2166136261

    for (
        let index = 0;
        index < value.length;
        index += 1
    ) {
        hash ^= value.charCodeAt(index)
        hash = Math.imul(hash, 16777619)
    }

    return hash >>> 0
}

function createSeededRandom(seed: number) {
    let state = seed

    return () => {
        state += 0x6d2b79f5

        let value = state

        value = Math.imul(
            value ^ (value >>> 15),
            value | 1,
        )

        value ^=
            value +
            Math.imul(
                value ^ (value >>> 7),
                value | 61,
            )

        return (
            ((value ^ (value >>> 14)) >>> 0) /
            4294967296
        )
    }
}

function shuffleDeterministic<T>(
    items: T[],
    seed: string,
) {
    const result = [...items].sort((a, b) => {
        const first = String(
            (a as { id?: string }).id ?? '',
        )

        const second = String(
            (b as { id?: string }).id ?? '',
        )

        return first.localeCompare(second)
    })

    const random = createSeededRandom(
        hashString(seed),
    )

    for (
        let index = result.length - 1;
        index > 0;
        index -= 1
    ) {
        const targetIndex = Math.floor(
            random() * (index + 1),
        )

        ;[
            result[index],
            result[targetIndex],
        ] = [
            result[targetIndex],
            result[index],
        ]
    }

    return result
}

function getMessagePosition(
    messages: ChatEntry[],
    index: number,
) {
    const current = messages[index]

    if (
        !current ||
        current.type !== 'message'
    ) {
        return 'single' as const
    }

    const previous = messages[index - 1]
    const next = messages[index + 1]

    const sameAsPrevious =
        previous?.type === 'message' &&
        previous.author === current.author

    const sameAsNext =
        next?.type === 'message' &&
        next.author === current.author

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
    const messagesEndRef =
        useRef<HTMLDivElement>(null)

    const submitAnswer = useSubmitAnswer()

    const [
        submittedAnswers,
        setSubmittedAnswers,
    ] = useState<SubmittedAnswer[]>([])

    const [
        isFinalAnswerSubmitted,
        setIsFinalAnswerSubmitted,
    ] = useState(false)

    const node = attempt.currentNode
    const playerRole = attempt.scenario.role

    const choices = useMemo(() => {
        if (!node) {
            return []
        }

        return shuffleDeterministic(
            node.choices,
            `${attempt.id}:${node.id}`,
        )
    }, [
        attempt.id,
        node,
    ])

    const messages = useMemo<ChatEntry[]>(() => {
        const historyNodeIds = new Set(
            attempt.history.map(
                (historyItem) =>
                    historyItem.node.id,
            ),
        )

        const historyMessages =
            attempt.history.flatMap<ChatEntry>(
                (historyItem) => {
                    const entries: ChatEntry[] =
                        [
                            ...historyItem.node.messages.map(
                                (
                                    message,
                                    index,
                                ) => ({
                                    type: 'message' as const,
                                    id: `${historyItem.node.id}-message-${index}`,
                                    author: message.author,
                                    text: message.text,
                                }),
                            ),

                            {
                                type: 'message',
                                id: `${historyItem.node.id}-choice-${historyItem.selectedChoice.id}`,
                                author: playerRole,
                                text: historyItem
                                    .selectedChoice
                                    .text,
                            },
                        ]

                    const submittedAnswer =
                        submittedAnswers.find(
                            (answer) =>
                                answer.nodeId ===
                                historyItem.node.id,
                        )

                    if (
                        submittedAnswer?.consequence
                    ) {
                        entries.push({
                            type: 'consequence',
                            id: `${historyItem.node.id}-consequence`,
                            text: submittedAnswer.consequence,
                        })
                    }

                    return entries
                },
            )

        const currentMessages: ChatEntry[] =
            node?.messages.map(
                (message, index) => ({
                    type: 'message',
                    id: `${node.id}-message-${index}`,
                    author: message.author,
                    text: message.text,
                }),
            ) ?? []

        const pendingMessages =
            submittedAnswers
                .filter(
                    (answer) =>
                        !historyNodeIds.has(
                            answer.nodeId,
                        ),
                )
                .flatMap<ChatEntry>(
                    (answer) => [
                        {
                            type: 'message',
                            id: `${answer.nodeId}-pending-choice-${answer.choiceId}`,
                            author: playerRole,
                            text: answer.choiceText,
                        },

                        {
                            type: 'consequence',
                            id: `${answer.nodeId}-pending-consequence`,
                            text: answer.consequence,
                        },
                    ],
                )

        return [
            ...historyMessages,
            ...currentMessages,
            ...pendingMessages,
        ]
    }, [
        attempt.history,
        node,
        playerRole,
        submittedAnswers,
    ])

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({
            behavior: 'smooth',
            block: 'end',
        })
    }, [
        messages.length,
        isFinalAnswerSubmitted,
    ])

    async function handleChoiceSelect(
        choiceId: string,
    ) {
        if (
            !node ||
            submitAnswer.isPending ||
            isFinalAnswerSubmitted
        ) {
            return
        }

        const choice = node.choices.find(
            (item) => item.id === choiceId,
        )

        if (!choice) {
            return
        }

        const result =
            await submitAnswer.mutateAsync({
                attemptId: attempt.id,
                nodeId: node.id,
                choiceId,
            })

        setSubmittedAnswers((current) => [
            ...current.filter(
                (answer) =>
                    answer.nodeId !==
                    result.nodeId,
            ),

            {
                nodeId: result.nodeId,
                choiceId: result.choiceId,
                choiceText: choice.text,
                consequence:
                    result.consequence,
            },
        ])

        if (result.completed) {
            setIsFinalAnswerSubmitted(true)
        }
    }

    if (
        attempt.status === 'completed' &&
        !isFinalAnswerSubmitted
    ) {
        return null
    }

    if (!node && !isFinalAnswerSubmitted) {
        return (
            <div className={styles.messages}>
                <div className={styles.error}>
                    Не удалось получить текущий шаг
                    сценария.
                </div>
            </div>
        )
    }

    return (
        <div className={styles.messages}>
            {messages.map(
                (message, index) => {
                    if (
                        message.type ===
                        'consequence'
                    ) {
                        return (
                            <div
                                key={message.id}
                                className={
                                    styles.consequence
                                }
                            >
                                <span
                                    className={
                                        styles.consequenceLabel
                                    }
                                >
                                    Последствие
                                    выбора
                                </span>

                                <span
                                    className={
                                        styles.consequenceText
                                    }
                                >
                                    {message.text}
                                </span>
                            </div>
                        )
                    }

                    const direction =
                        message.author ===
                        playerRole
                            ? 'outgoing'
                            : 'incoming'

                    const position =
                        getMessagePosition(
                            messages,
                            index,
                        )

                    const showAvatar =
                        direction ===
                            'incoming' &&
                        (position ===
                            'single' ||
                            position ===
                                'last')

                    return (
                        <ChatMessage
                            key={message.id}
                            direction={
                                direction
                            }
                            position={
                                position
                            }
                            text={
                                message.text
                            }
                            avatarText={getAuthorName(
                                message.author,
                            )}
                            showAvatar={
                                showAvatar
                            }
                        />
                    )
                },
            )}

            {!isFinalAnswerSubmitted &&
                node && (
                    <div
                        className={
                            styles.answers
                        }
                    >
                        <span
                            className={
                                styles.answersLabel
                            }
                        >
                            Выберите ответ
                        </span>

                        {choices.map(
                            (choice) => (
                                <button
                                    key={
                                        choice.id
                                    }
                                    type="button"
                                    className={
                                        styles.answer
                                    }
                                    disabled={
                                        submitAnswer.isPending
                                    }
                                    onClick={() => {
                                        void handleChoiceSelect(
                                            choice.id,
                                        )
                                    }}
                                >
                                    {
                                        choice.text
                                    }
                                </button>
                            ),
                        )}

                        {submitAnswer.isError && (
                            <span
                                className={
                                    styles.error
                                }
                            >
                                Не удалось
                                отправить ответ.
                                Попробуйте ещё
                                раз.
                            </span>
                        )}
                    </div>
                )}

            {isFinalAnswerSubmitted && (
                <div
                    className={
                        styles.complete
                    }
                >
                    <button
                        type="button"
                        className={
                            styles.completeButton
                        }
                        onClick={onComplete}
                    >
                        Посмотреть результат
                    </button>
                </div>
            )}

            <div ref={messagesEndRef} />
        </div>
    )
}