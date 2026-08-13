import { Link } from 'react-router-dom'

import { hasScamOrNotTopic } from '@/entities/scam-or-not'
import type { AttemptFeedback } from '@/entities/training-attempt'
import { Button } from '@/shared/ui/button'

import styles from './TrainingChat.module.scss'

type TrainingSummaryProps = {
    title?: string
    score?: number
    feedback?: AttemptFeedback
    isResultPending: boolean
    isFeedbackPending: boolean
    isResultError: boolean
    isFeedbackError: boolean
    logicalScenarioId?: string
}

export function TrainingSummary({
    title,
    score,
    feedback,
    isResultPending,
    isFeedbackPending,
    isResultError,
    isFeedbackError,
    logicalScenarioId,
}: TrainingSummaryProps) {
    const hasTopicGame = logicalScenarioId
        ? hasScamOrNotTopic(logicalScenarioId)
        : false

    return (
        <div className={styles.summary}>
            <div className={styles.top}>
                <div className={styles.summaryHeader}>
                    <span className={styles.summaryIcon}>✓</span>

                    <h2 className={styles.summaryTitle}>Тренировка завершена</h2>

                    {title && <p className={styles.summaryScenario}>{title}</p>}
                </div>

                {isResultPending ? (
                    <p>Подсчитываем результат...</p>
                ) : isResultError ? (
                    <p className={styles.summaryError}>Не удалось загрузить результат.</p>
                ) : (
                    <div className={styles.summaryScore}>
                        <strong className={styles.summaryScoreValue}>{score ?? 0}</strong>

                        <span className={styles.summaryScoreLabel}>баллов</span>
                    </div>
                )}
            </div>

            {isFeedbackPending && (
                <div className={styles.summaryLoading}>Анализируем прохождение...</div>
            )}

            {isFeedbackError && (
                <p className={styles.summaryError}>Не удалось загрузить персональный разбор.</p>
            )}

            {feedback && (
                <div className={styles.summaryContent}>
                    {feedback.motivation && (
                        <div className={styles.summaryMotivation}>{feedback.motivation}</div>
                    )}

                    {feedback.strengths.length > 0 && (
                        <div className={styles.summarySection}>
                            <h3>Что получилось хорошо</h3>

                            <ul>
                                {feedback.strengths.map((item) => (
                                    <li key={item}>{item}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    {feedback.weaknesses.length > 0 && (
                        <div className={styles.summarySection}>
                            <h3>На что обратить внимание</h3>

                            <ul>
                                {feedback.weaknesses.map((item) => (
                                    <li key={item}>{item}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    {feedback.riskProfile && (
                        <div className={styles.summarySection}>
                            <h3>Профиль риска</h3>

                            {feedback.riskProfile.dominantRisk && (
                                <p>
                                    <strong>{feedback.riskProfile.dominantRisk}</strong>
                                </p>
                            )}

                            {feedback.riskProfile.description && (
                                <p>{feedback.riskProfile.description}</p>
                            )}
                        </div>
                    )}

                    {feedback.recommendations.length > 0 && (
                        <div className={styles.summarySection}>
                            <h3>Рекомендации</h3>

                            <ul>
                                {feedback.recommendations.map((item) => (
                                    <li key={item}>{item}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    {feedback.learningTips.length > 0 && (
                        <div className={styles.summarySection}>
                            <h3>Что стоит изучить</h3>

                            <ul>
                                {feedback.learningTips.map((item) => (
                                    <li key={item}>{item}</li>
                                ))}
                            </ul>
                        </div>
                    )}
                </div>
            )}

            {hasTopicGame && (
                <div className={styles.summaryGameCta}>
                    <div>
                        <strong>Закрепите тему</strong>
                        <p>
                            Отличите мошеннические ситуации от безопасных в
                            короткой игре.
                        </p>
                    </div>
                    <Link to={`/training/scam-or-not/${logicalScenarioId}`}>
                        <Button>Закрепить тему в мини-игре</Button>
                    </Link>
                </div>
            )}
        </div>
    )
}
