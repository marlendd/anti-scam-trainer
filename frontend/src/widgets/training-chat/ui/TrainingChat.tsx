import {
    type TrainingAttemptState,
    useAttemptFeedback,
    useAttemptResult,
    useRestartAttempt,
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

import GeForceImage from '@/shared/assets/images/geforce.jpg';
import IntelImage from '@/shared/assets/images/intel.jpg';
import IphoneImage from '@/shared/assets/images/iphone.webp';
import PlayStation from '@/shared/assets/images/playstation.jpg';



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
        isInteractive &&
        attempt?.status === 'completed'

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

    const restartAttempt = useRestartAttempt()

    if (isInteractive && !attempt) {
        return (
            <section className={styles.chat}>
                <div className={styles.empty}>
                    <h2 className={styles.emptyTitle}>
                        Выберите сценарий
                    </h2>

                    <p
                        className={
                            styles.emptyDescription
                        }
                    >
                        Выберите доступный сценарий,
                        чтобы начать тренировку.
                    </p>
                </div>
            </section>
        )
    }

    const imageMock = {
        'Дешёвый процессор и предоплата': IntelImage,
        'Видеокарта и фальшивая поддержка': GeForceImage,
        'Поддельная оплата и звонок поддержки': IphoneImage,
        'Фиктивная переплата и возврат разницы': PlayStation,
    }

    const image = imageMock[
        attempt?.scenario.title as keyof typeof imageMock
    ]

    if (isCompleted && attempt) {
        return (
            <section className={styles.chat}>
                <TrainingSummary
                    title={attempt.scenario.title}
                    score={
                        result?.score ??
                        attempt.score
                    }
                    ending={
                        attempt.ending ??
                        undefined
                    }
                    feedback={feedback}
                    logicalScenarioId={scenarioSummary?.logicalId}
                    isResultPending={
                        isResultPending
                    }
                    isFeedbackPending={
                        isFeedbackPending
                    }
                    isResultError={
                        isResultError
                    }
                    isFeedbackError={
                        isFeedbackError
                    }
                />
            </section>
        )
    }



    if (isInteractive && attempt) {
        // @ts-ignore
        return (
            <section className={styles.chat}>
                <ChatProductHeader
                    participantName="Собеседник"
                    participantStatus="В сети"
                    productTitle={
                        attempt.scenario.product
                            .title
                    }
                    productPrice={
                        attempt.scenario.product
                            .price
                    }
                    imageUrl={image}
                    onRestart={() => {
                        restartAttempt.mutate(
                            attempt.scenarioId,
                        )
                    }}
                    isRestarting={
                        restartAttempt.isPending
                    }
                />

                <AttemptPlayer
                    attempt={attempt}
                    onComplete={
                        onScenarioComplete
                    }
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

                    <p
                        className={
                            styles.emptyDescription
                        }
                    >
                        Выберите доступный сценарий,
                        чтобы начать тренировку.
                    </p>
                </div>
            </section>
        )
    }

    const interlocutor =
        scenario.participants.find(
            (participant) =>
                participant.id !==
                scenario.playerParticipantId,
        )


    return (
        <section className={styles.chat}>
            <ChatProductHeader
                participantName={
                    interlocutor?.name ??
                    'Илья Щебень'
                }
                participantStatus={
                    interlocutor?.status
                }
                productTitle={
                    scenario.product.title
                }
                productPrice={
                    scenario.product.price
                }
                imageUrl={
                    scenario.product.imageUrl
                }
            />

            <ScenarioPlayer
                scenario={scenario}
                mode={mode}
                onComplete={
                    onScenarioComplete
                }
            />
        </section>
    )
}
