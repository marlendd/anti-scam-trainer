import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'

import {
    getMixedScamOrNotCards,
    getScamOrNotTopic,
    getScamOrNotTopicCards,
    shuffleScamOrNotCards,
    type ScamOrNotCard,
} from '@/entities/scam-or-not'
import { useDocumentTitle } from '@/shared/lib/use-document-title'
import { Button } from '@/shared/ui/button'

import styles from './ScamOrNotPage.module.scss'

function createDeck(topicId?: string): ScamOrNotCard[] {
    if (!topicId) {
        return getMixedScamOrNotCards()
    }

    return shuffleScamOrNotCards(getScamOrNotTopicCards(topicId))
}

export function ScamOrNotPage() {
    const { logicalScenarioId } = useParams()

    return (
        <ScamOrNotGame
            key={logicalScenarioId ?? 'mixed'}
            logicalScenarioId={logicalScenarioId}
        />
    )
}

type ScamOrNotGameProps = {
    logicalScenarioId?: string
}

function ScamOrNotGame({ logicalScenarioId }: ScamOrNotGameProps) {
    const topic = logicalScenarioId
        ? getScamOrNotTopic(logicalScenarioId)
        : undefined
    const isUnknownTopic = Boolean(logicalScenarioId && !topic)

    const [deck, setDeck] = useState(() => createDeck(logicalScenarioId))
    const [currentIndex, setCurrentIndex] = useState(0)
    const [selectedAnswer, setSelectedAnswer] = useState<boolean | null>(null)
    const [correctCount, setCorrectCount] = useState(0)

    useDocumentTitle('Мошенник или нет?')

    if (isUnknownTopic || deck.length === 0) {
        return (
            <main className={styles.page}>
                <section className={styles.empty}>
                    <span className={styles.eyebrow}>Мини-игра</span>
                    <h1 className={styles.title}>Набор карточек не найден</h1>
                    <p className={styles.subtitle}>
                        Для этой темы пока нет отдельной игры. Попробуйте
                        смешанный режим.
                    </p>
                    <Link to="/training/scam-or-not">
                        <Button size="large">Играть в смешанном режиме</Button>
                    </Link>
                </section>
            </main>
        )
    }

    const isCompleted = currentIndex >= deck.length
    const currentCard = deck[currentIndex]
    const progress = isCompleted
        ? 100
        : Math.round((currentIndex / deck.length) * 100)
    const resultPercent = Math.round((correctCount / deck.length) * 100)

    function handleAnswer(answer: boolean) {
        if (selectedAnswer !== null || !currentCard) {
            return
        }

        setSelectedAnswer(answer)

        if (answer === currentCard.isScam) {
            setCorrectCount((count) => count + 1)
        }
    }

    function handleNext() {
        setCurrentIndex((index) => index + 1)
        setSelectedAnswer(null)
    }

    function handleRestart() {
        setDeck(createDeck(logicalScenarioId))
        setCurrentIndex(0)
        setSelectedAnswer(null)
        setCorrectCount(0)
    }

    if (isCompleted) {
        return (
            <main className={styles.page}>
                <section
                    className={styles.result}
                    aria-labelledby="scam-game-result-title"
                >
                    <span className={styles.resultIcon} aria-hidden="true">
                        {resultPercent >= 70 ? '✓' : '!'}
                    </span>
                    <span className={styles.eyebrow}>
                        {topic?.title ?? 'Смешанный режим'}
                    </span>
                    <h1 id="scam-game-result-title" className={styles.title}>
                        Игра завершена
                    </h1>
                    <p className={styles.resultScore}>{resultPercent}%</p>
                    <p className={styles.subtitle}>
                        Вы верно определили {correctCount} из {deck.length}{' '}
                        ситуаций.
                    </p>
                    <div className={styles.resultActions}>
                        <Button size="large" onClick={handleRestart}>
                            Сыграть ещё раз
                        </Button>
                        <Link to="/training/role-selection">
                            <Button variant="secondary" size="large">
                                К сценариям
                            </Button>
                        </Link>
                    </div>
                </section>
            </main>
        )
    }

    const isCorrect = selectedAnswer === currentCard.isScam

    return (
        <main className={styles.page}>
            <header className={styles.heading}>
                <span className={styles.eyebrow}>
                    {topic?.title ?? 'Смешанный режим'}
                </span>
                <h1 className={styles.title}>Мошенник или нет?</h1>
                <p className={styles.subtitle}>
                    Оцените ситуацию, затем сверьтесь с разбором.
                </p>
            </header>

            <section className={styles.game} aria-label="Игровая карточка">
                <div className={styles.progressRow}>
                    <span>
                        Ситуация {currentIndex + 1} из {deck.length}
                    </span>
                    <span>{correctCount} верно</span>
                </div>
                <div
                    className={styles.progressTrack}
                    role="progressbar"
                    aria-label="Прогресс игры"
                    aria-valuemin={0}
                    aria-valuemax={100}
                    aria-valuenow={progress}
                >
                    <span
                        className={styles.progressFill}
                        style={{ width: `${progress}%` }}
                    />
                </div>

                <article className={styles.card}>
                    <p className={styles.situation}>{currentCard.situation}</p>

                    <div className={styles.answers} aria-label="Выберите ответ">
                        <button
                            type="button"
                            className={styles.answer}
                            data-answer="scam"
                            data-selected={selectedAnswer === true}
                            disabled={selectedAnswer !== null}
                            onClick={() => handleAnswer(true)}
                        >
                            Мошенник
                        </button>
                        <button
                            type="button"
                            className={styles.answer}
                            data-answer="safe"
                            data-selected={selectedAnswer === false}
                            disabled={selectedAnswer !== null}
                            onClick={() => handleAnswer(false)}
                        >
                            Не мошенник
                        </button>
                    </div>

                    {selectedAnswer !== null && (
                        <div
                            className={styles.feedback}
                            data-correct={isCorrect}
                            role="status"
                            aria-live="polite"
                        >
                            <strong className={styles.feedbackTitle}>
                                {isCorrect ? 'Верно!' : 'Не совсем.'}{' '}
                                {currentCard.isScam
                                    ? 'Это мошенническая схема.'
                                    : 'Это безопасное поведение.'}
                            </strong>
                            <p>{currentCard.explanation}</p>
                            <ul
                                className={styles.signs}
                                aria-label="Ключевые признаки"
                            >
                                {currentCard.riskSigns.map((sign) => (
                                    <li key={sign}>{sign}</li>
                                ))}
                            </ul>
                            <Button
                                className={styles.nextButton}
                                onClick={handleNext}
                            >
                                {currentIndex === deck.length - 1
                                    ? 'Узнать результат'
                                    : 'Следующая ситуация'}
                            </Button>
                        </div>
                    )}
                </article>
            </section>
        </main>
    )
}
