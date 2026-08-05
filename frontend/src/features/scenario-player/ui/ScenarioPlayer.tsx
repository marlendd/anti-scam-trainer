// features/scenario-player/ui/ScenarioPlayer.tsx

import { useEffect, useRef } from 'react'

import type {
  ScenarioMessage,
  TrainingScenario,
} from '@/entities/training-scenario'
import { ChatMessage } from '@/shared/ui/chat-message'

import { useScenarioPlayback } from '../model/useScenarioPlayback'

import styles from './ScenarioPlayer.module.scss'

type ScenarioPlayerProps = {
  scenario: TrainingScenario
  onComplete?: () => void
}

function getMessagePosition(
  messages: ScenarioMessage[],
  index: number,
) {
  const current = messages[index]
  const previous = messages[index - 1]
  const next = messages[index + 1]

  const sameAsPrevious =
    previous?.senderId === current?.senderId

  const sameAsNext =
    next?.senderId === current?.senderId

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

export function ScenarioPlayer({
  scenario,
  onComplete,
}: ScenarioPlayerProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const {
    visibleMessages,
    isCompleted,
  } = useScenarioPlayback(scenario.messages)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({
      behavior: 'smooth',
      block: 'end',
    })
  }, [visibleMessages.length])

  useEffect(() => {
    if (isCompleted) {
      onComplete?.()
    }
  }, [isCompleted, onComplete])

  return (
    <div className={styles.messages}>
      {visibleMessages.map((message, index) => {
        const sender = scenario.participants.find(
          (participant) =>
            participant.id === message.senderId,
        )

        const direction =
          message.senderId === scenario.playerParticipantId
            ? 'outgoing'
            : 'incoming'

        const position = getMessagePosition(
          visibleMessages,
          index,
        )

        const showAvatar =
          direction === 'incoming' &&
          (position === 'single' || position === 'last')

        return (
          <ChatMessage
            key={message.id}
            direction={direction}
            position={position}
            text={message.text}
            avatarText={sender?.name}
            showAvatar={showAvatar}
          />
        )
      })}

      {!isCompleted && (
        <div className={styles.typing}>
          Собеседник печатает…
        </div>
      )}

      <div ref={messagesEndRef} />
    </div>
  )
}