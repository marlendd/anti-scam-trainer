import type { TrainingScenario } from '@/entities/training-scenario'
import { ScenarioPlayer, type ScenarioPlaybackMode } from '@/features/scenario-player'
import { ChatProductHeader } from '@/shared/ui/chat-product-header'

import styles from './TrainingChat.module.scss'

type TrainingChatProps = {
    scenario: TrainingScenario | null
    mode?: ScenarioPlaybackMode
    onScenarioComplete?: () => void
}

export function TrainingChat({
    scenario,
    mode = 'preview',
    onScenarioComplete,
}: TrainingChatProps) {
    if (!scenario) {
        return (
            <section className={styles.chat}>
                <div className={styles.empty}>
                    <h2 className={styles.emptyTitle}>Выберите сценарий</h2>

                    <p className={styles.emptyDescription}>
                        Выберите доступный сценарий, чтобы начать тренировку.
                    </p>
                </div>
            </section>
        )
    }

    const interlocutor = scenario.participants.find(
        (participant) => participant.id !== scenario.playerParticipantId,
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

            <ScenarioPlayer scenario={scenario} mode={mode} onComplete={onScenarioComplete} />
        </section>
    )
}
