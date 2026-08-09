import type {TrainingAttemptState} from '@/entities/training-attempt'
import type {TrainingScenario} from '@/entities/training-scenario'
import {AttemptPlayer} from '@/features/attempt-player'

import messageImage from '@/shared/assets/images/message.png';

import {
    ScenarioPlayer,
    type ScenarioPlaybackMode,
} from '@/features/scenario-player'
import {ChatProductHeader} from '@/shared/ui/chat-product-header'

import styles from './TrainingChat.module.scss'

type TrainingChatProps = {
    scenario?: TrainingScenario | null
    attempt?: TrainingAttemptState | null
    mode?: ScenarioPlaybackMode
    onScenarioComplete?: () => void
}

export function TrainingChat({
                                 scenario = null,
                                 attempt = null,
                                 mode = 'preview',
                                 onScenarioComplete,
                             }: TrainingChatProps) {
    if (mode === 'interactive') {
        if (!attempt) {
            return (
                <section className={styles.chat}>
                    <div className={styles.empty}>
                        <img src={messageImage} alt="" className={styles.image}/>
                        <h2 className={styles.emptyTitle}>
                            Чат не выбран
                        </h2>

                        <p className={styles.emptyDescription}>
                            Выберите сценарий, чтобы начать тренировку.
                        </p>
                    </div>
                </section>
            )
        }

        return (
            <section className={styles.chat}>
                <ChatProductHeader
                    participantName="Илья Щебень"
                    participantStatus="В сети"
                    productTitle={scenario?.title ?? 'Видео карта GeForce Nvidia 8600 Ti RTX Ryzen'}
                    productPrice={56000}
                    imageUrl=""
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
                participantName={interlocutor?.name ?? 'Собеседник'}
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