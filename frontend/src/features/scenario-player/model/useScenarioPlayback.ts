import {
  useCallback,
  useEffect,
  useState,
} from 'react'

import type {
  ScenarioAnswerOption,
  ScenarioMessage,
  ScenarioTimelineItem,
} from '@/entities/training-scenario'

export type ScenarioResult = {
  correctAnswers: number
  totalAnswers: number
}

export type ScenarioPlaybackMode =
  | 'preview'
  | 'interactive'

type SelectedAnswer = {
  choiceId: string
  answer: ScenarioAnswerOption
}

type UseScenarioPlaybackParams = {
  timeline: ScenarioTimelineItem[]
  playerParticipantId: string
  mode: ScenarioPlaybackMode
}

export function useScenarioPlayback({
  timeline,
  playerParticipantId,
  mode,
}: UseScenarioPlaybackParams) {
  const [currentIndex, setCurrentIndex] = useState(0)

  const [visibleMessages, setVisibleMessages] = useState<
    ScenarioMessage[]
  >([])

  const [selectedAnswers, setSelectedAnswers] = useState<
    SelectedAnswer[]
  >([])

  const [isWaitingForAnswer, setIsWaitingForAnswer] =
    useState(false)

  const [isTyping, setIsTyping] = useState(false)

  const currentItem = timeline[currentIndex]

  const currentChoice =
    currentItem?.type === 'choice'
      ? currentItem
      : null

  const isCompleted =
    currentIndex >= timeline.length

  const appendAnswer = useCallback(
    (
      choiceId: string,
      answer: ScenarioAnswerOption,
    ) => {
      setVisibleMessages((current) => [
        ...current,
        {
          type: 'message',
          id: `answer-${choiceId}-${answer.id}`,
          senderId: playerParticipantId,
          text: answer.text,
        },
      ])

      setSelectedAnswers((current) => [
        ...current,
        {
          choiceId,
          answer,
        },
      ])

      setIsWaitingForAnswer(false)
      setCurrentIndex((current) => current + 1)
    },
    [playerParticipantId],
  )

  const selectAnswer = useCallback(
    (answer: ScenarioAnswerOption) => {
      if (!currentChoice) {
        return
      }

      appendAnswer(currentChoice.id, answer)
    },
    [appendAnswer, currentChoice],
  )

  useEffect(() => {
    if (!currentItem || isCompleted) {
      return
    }

    if (currentItem.type === 'choice') {
      if (mode === 'interactive') {
        setIsWaitingForAnswer(true)
        return
      }

      const previewAnswer =
        currentItem.options.find(
          (option) =>
            option.id === currentItem.previewOptionId,
        ) ?? currentItem.options[0]

      if (!previewAnswer) {
        setCurrentIndex((current) => current + 1)
        return
      }

      const timeout = window.setTimeout(() => {
        appendAnswer(
          currentItem.id,
          previewAnswer,
        )
      }, 500)

      return () => {
        window.clearTimeout(timeout)
      }
    }

    setIsTyping(true)

    const timeout = window.setTimeout(() => {
      setVisibleMessages((current) => [
        ...current,
        currentItem,
      ])

      setIsTyping(false)
      setCurrentIndex((current) => current + 1)
    }, currentItem.delayMs ?? 500)

    return () => {
      window.clearTimeout(timeout)
    }
  }, [
    appendAnswer,
    currentItem,
    isCompleted,
    mode,
  ])

  return {
    visibleMessages,
    currentChoice,
    selectedAnswers,

    isWaitingForAnswer,
    isTyping,
    isCompleted,

    selectAnswer,
  }
}