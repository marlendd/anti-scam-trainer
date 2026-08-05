import type { TrainingScenario } from '@/entities/training-scenario'
import { ScenarioPlayer } from '@/features/scenario-player'
import { ChatProductHeader } from '@/shared/ui/chat-product-header'

import styles from './TrainingChat.module.scss'

type TrainingChatProps = {
  scenario: TrainingScenario
  onScenarioComplete?: () => void
}

export function TrainingChat({
  scenario,
  onScenarioComplete,
}: TrainingChatProps) {
  const interlocutor = scenario.participants.find(
    (participant) =>
      participant.id !== scenario.playerParticipantId,
  )

  return (
    <section className={styles.chat}>
      <ChatProductHeader
        participantName={
          interlocutor?.name ?? 'Собеседник'
        }
        participantStatus={interlocutor?.status}
        productTitle={scenario.product.title}
        productPrice={scenario.product.price}
        imageUrl={scenario.product.imageUrl}
      />

      <ScenarioPlayer
        scenario={scenario}
        onComplete={onScenarioComplete}
      />
    </section>
  )
}