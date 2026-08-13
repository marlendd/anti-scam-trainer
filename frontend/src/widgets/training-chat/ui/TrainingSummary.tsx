import type { AttemptFeedback } from '@/entities/training-attempt'

import styles from './TrainingChat.module.scss'

type ScenarioEnding = {
    id: string
    header: string
    result: string
}

type TrainingSummaryProps = {
    title?: string
    score?: number
    ending?: ScenarioEnding
    feedback?: AttemptFeedback
    isResultPending: boolean
    isFeedbackPending: boolean
    isResultError: boolean
    isFeedbackError: boolean
}

export function TrainingSummary({
    title,
    score,
    ending,
    feedback,
    isResultPending,
    isFeedbackPending,
    isResultError,
    isFeedbackError,
}: TrainingSummaryProps) {
    return (
        <div className={styles.summary}>
            <div className={styles.top}>
                <div className={styles.summaryHeader}>
                    <span className={styles.summaryIcon}>
                        ✓
                    </span>

                    <div className={styles.title}>
                        <h2
                            className={
                                styles.summaryTitle
                            }
                        >
                            Тренировка завершена
                        </h2>

                        {title && (
                            <p
                                className={
                                    styles.summaryScenario
                                }
                            >
                                {title}
                            </p>
                        )}
                    </div>
                </div>

                {isResultPending ? (
                    <p>
                        Подсчитываем результат...
                    </p>
                ) : isResultError ? (
                    <p
                        className={
                            styles.summaryError
                        }
                    >
                        Не удалось загрузить
                        результат.
                    </p>
                ) : (
                    <div
                        className={
                            styles.summaryScore
                        }
                    >
                        <strong
                            className={
                                styles.summaryScoreValue
                            }
                        >
                            {score ?? 0}
                        </strong>

                        <span
                            className={
                                styles.summaryScoreLabel
                            }
                        >
                            баллов
                        </span>
                    </div>
                )}
            </div>

            {ending && (
                <div
                    className={
                        styles.summaryEnding
                    }
                >
                    <span
                        className={
                            styles.summaryEndingLabel
                        }
                    >
                        Результат сценария
                    </span>

                    <h3
                        className={
                            styles.summaryEndingTitle
                        }
                    >
                        {ending.header}
                    </h3>

                    <p
                        className={
                            styles.summaryEndingResult
                        }
                    >
                        {ending.result}
                    </p>
                </div>
            )}

            {isFeedbackPending && (
                <div
                    className={
                        styles.summaryLoading
                    }
                >
                    Анализируем прохождение...
                </div>
            )}

            {isFeedbackError && (
                <p
                    className={
                        styles.summaryError
                    }
                >
                    Не удалось загрузить
                    персональный разбор.
                </p>
            )}

            {feedback && (
                <div
                    className={
                        styles.summaryContent
                    }
                >
                    {feedback.motivation && (
                        <div
                            className={
                                styles.summaryMotivation
                            }
                        >
                            {feedback.motivation}
                        </div>
                    )}

                    {feedback.strengths.length >
                        0 && (
                        <div
                            className={
                                styles.summarySection
                            }
                        >
                            <h3>
                                Что получилось хорошо
                            </h3>

                            <ul>
                                {feedback.strengths.map(
                                    (item) => (
                                        <li
                                            key={
                                                item
                                            }
                                        >
                                            {
                                                item
                                            }
                                        </li>
                                    ),
                                )}
                            </ul>
                        </div>
                    )}

                    {feedback.weaknesses.length >
                        0 && (
                        <div
                            className={
                                styles.summarySection
                            }
                        >
                            <h3>
                                На что обратить
                                внимание
                            </h3>

                            <ul>
                                {feedback.weaknesses.map(
                                    (item) => (
                                        <li
                                            key={
                                                item
                                            }
                                        >
                                            {
                                                item
                                            }
                                        </li>
                                    ),
                                )}
                            </ul>
                        </div>
                    )}

                    {feedback.riskProfile && (
                        <div
                            className={
                                styles.summarySection
                            }
                        >
                            <h3>
                                Профиль риска
                            </h3>

                            {feedback.riskProfile
                                .dominantRisk && (
                                <p>
                                    <strong>
                                        {
                                            feedback
                                                .riskProfile
                                                .dominantRisk
                                        }
                                    </strong>
                                </p>
                            )}

                            {feedback.riskProfile
                                .description && (
                                <p>
                                    {
                                        feedback
                                            .riskProfile
                                            .description
                                    }
                                </p>
                            )}
                        </div>
                    )}

                    {feedback.recommendations
                        .length > 0 && (
                        <div
                            className={
                                styles.summarySection
                            }
                        >
                            <h3>
                                Рекомендации
                            </h3>

                            <ul>
                                {feedback.recommendations.map(
                                    (item) => (
                                        <li
                                            key={
                                                item
                                            }
                                        >
                                            {
                                                item
                                            }
                                        </li>
                                    ),
                                )}
                            </ul>
                        </div>
                    )}

                    {feedback.learningTips.length >
                        0 && (
                        <div
                            className={
                                styles.summarySection
                            }
                        >
                            <h3>
                                Что стоит изучить
                            </h3>

                            <ul>
                                {feedback.learningTips.map(
                                    (item) => (
                                        <li
                                            key={
                                                item
                                            }
                                        >
                                            {
                                                item
                                            }
                                        </li>
                                    ),
                                )}
                            </ul>
                        </div>
                    )}
                </div>
            )}
        </div>
    )
}