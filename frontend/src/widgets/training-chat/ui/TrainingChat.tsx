import type { TrainingAttemptState } from '@/entities/training-attempt'
import {
    useAttemptFeedback,
    useAttemptResult,
} from '@/entities/training-attempt'
import type {
    TrainingScenario,
    TrainingScenarioSummary,
} from '@/entities/training-scenario'
import { AttemptPlayer } from '@/features/attempt-player'
import {
    ScenarioPlayer,
    type ScenarioPlaybackMode,
} from '@/features/scenario-player'
import { ChatProductHeader } from '@/shared/ui/chat-product-header'

import { TrainingSummary } from './TrainingSummary'

import styles from './TrainingChat.module.scss'

type TrainingChatProps = {
    scenario?: TrainingScenario | null
    scenarioSummary?: TrainingScenarioSummary | null
    attempt?: TrainingAttemptState | null
    mode?: ScenarioPlaybackMode
    onScenarioComplete?: () => void
}

export function TrainingChat({
    scenario = null,
    scenarioSummary = null,
    attempt = null,
    mode = 'preview',
    onScenarioComplete,
}: TrainingChatProps) {
    const isInteractive = mode === 'interactive'
    const isCompleted =
        isInteractive && attempt?.status === 'completed'

    const {
        data: result,
        isPending: isResultPending,
        isError: isResultError,
    } = useAttemptResult(
        attempt?.id ?? null,
        isCompleted,
    )

    const {
        data: feedback,
        isPending: isFeedbackPending,
        isError: isFeedbackError,
    } = useAttemptFeedback(
        attempt?.id ?? null,
        isCompleted,
    )

    if (isInteractive && !attempt) {
        return (
            <section className={styles.chat}>
                <div className={styles.empty}>
                    <h2 className={styles.emptyTitle}>
                        Выберите сценарий
                    </h2>

                    <p className={styles.emptyDescription}>
                        Выберите доступный сценарий, чтобы начать тренировку.
                    </p>
                </div>
            </section>
        )
    }

    const tempImage = 'https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fmedia.karousell.com%2Fmedia%2Fphotos%2Fproducts%2F2023%2F6%2F24%2Fbarely_used_huawei_matebook_d__1687577820_f8691867_progressive.jpg&f=1&nofb=1&ipt=528ad03097a65b4fb0bbf0b26a65a431ef3c04eceb43ccaa58e140d83db78151'

    if (isCompleted && attempt) {
        return (
            <section className={styles.chat}>
                <TrainingSummary
                    title={attempt.scenario.title}
                    score={result?.score ?? attempt.score}
                    feedback={feedback}
                    isResultPending={isResultPending}
                    isFeedbackPending={isFeedbackPending}
                    isResultError={isResultError}
                    isFeedbackError={isFeedbackError}
                    logicalScenarioId={scenarioSummary?.logicalId}
                />
            </section>
        )
    }

    if (isInteractive && attempt) {
        return (
            <section className={styles.chat}>
                <ChatProductHeader
                    participantName="Собеседник"
                    participantStatus="В сети"
                    productTitle={attempt.scenario.product.title}
                    productPrice={attempt.scenario.product.price}
                    imageUrl={tempImage}
                />

                <AttemptPlayer
                    attempt={attempt}
                    onComplete={onScenarioComplete}
                />
            </section>
        )
    }

    if (!scenario) {
        return (
            <section className={styles.chat}>
                <div className={styles.empty}>
                    <h2 className={styles.emptyTitle}>
                        Выберите сценарий
                    </h2>

                    <p className={styles.emptyDescription}>
                        Выберите доступный сценарий, чтобы начать тренировку.
                    </p>
                </div>
            </section>
        )
    }

    const interlocutor = scenario.participants.find(
        (participant) =>
            participant.id !== scenario.playerParticipantId,
    )

    return (
        <section className={styles.chat}>
            <ChatProductHeader
                participantName={interlocutor?.name ?? 'Илья Щебень'}
                participantStatus={interlocutor?.status}
                productTitle={scenario.product.title}
                productPrice={scenario.product.price}
                imageUrl={scenario.product.imageUrl}
            />

            <ScenarioPlayer
                scenario={scenario}
                mode={mode}
                onComplete={onScenarioComplete}
            />
        </section>
    )
}
