import { useEffect, useRef } from 'react'

import type {
  ScenarioMessage,
  TrainingScenario,
} from '@/entities/training-scenario'
import { ChatMessage } from '@/shared/ui/chat-message'

import {
  type ScenarioResult,
  useScenarioPlayback,
} from '../model/useScenarioPlayback'

import styles from './ScenarioPlayer.module.scss'

export type ScenarioPlayerMode =
  | 'preview'
  | 'interactive'

type ScenarioPlayerProps = {
  scenario: TrainingScenario
  mode?: ScenarioPlayerMode
  onComplete?: (result: ScenarioResult) => void
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
  mode = 'preview',
  onComplete,
}: ScenarioPlayerProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const {
    visibleMessages,
    currentChoice,
    selectedAnswers,
    isWaitingForAnswer,
    isTyping,
    isCompleted,
    selectAnswer,
  } = useScenarioPlayback({
    timeline: scenario.timeline,
    playerParticipantId:
      scenario.playerParticipantId,
    mode,
  })

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({
      behavior: 'smooth',
      block: 'end',
    })
  }, [
    visibleMessages.length,
    isWaitingForAnswer,
  ])

  useEffect(() => {
    if (!isCompleted) {
      return
    }

    const totalAnswers = selectedAnswers.length

    const correctAnswers =
      selectedAnswers.filter(
        ({ answer }) => answer.isCorrect,
      ).length

    onComplete?.({
      totalAnswers,
      correctAnswers,
    })
  }, [
    isCompleted,
    onComplete,
    selectedAnswers,
  ])

  return (
    <div className={styles.messages}>
      {visibleMessages.map((message, index) => {
        const sender = scenario.participants.find(
          (participant) =>
            participant.id === message.senderId,
        )

        const direction =
          message.senderId ===
          scenario.playerParticipantId
            ? 'outgoing'
            : 'incoming'

        const position = getMessagePosition(
          visibleMessages,
          index,
        )

        const showAvatar =
          direction === 'incoming' &&
          (position === 'single' ||
            position === 'last')

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

      {isTyping && (
        <div className={styles.typing}>
          Собеседник печатает…
        </div>
      )}

      {mode === 'interactive' &&
        isWaitingForAnswer &&
        currentChoice && (
          <div className={styles.answers}>
            <span className={styles.answersLabel}>
              Выберите ответ
            </span>

            {currentChoice.options.map((option) => (
              <button
                key={option.id}
                type="button"
                className={styles.answer}
                onClick={() => {
                  selectAnswer(option)
                }}
              >
                {option.text}
              </button>
            ))}
          </div>
        )}

      <div ref={messagesEndRef} />
    </div>
  )
}