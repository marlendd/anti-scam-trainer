// features/scenario-player/model/useScenarioPlayback.ts

import { useEffect, useState } from 'react'

import type { ScenarioMessage } from '@/entities/training-scenario'

export function useScenarioPlayback(
  messages: ScenarioMessage[],
) {
  const [visibleCount, setVisibleCount] = useState(0)

  useEffect(() => {
    setVisibleCount(0)
  }, [messages])

  useEffect(() => {
    if (visibleCount >= messages.length) {
      return
    }

    const nextMessage = messages[visibleCount]
    const delay = nextMessage?.delayMs ?? 800

    const timeoutId = window.setTimeout(() => {
      setVisibleCount((count) => count + 1)
    }, delay)

    return () => {
      window.clearTimeout(timeoutId)
    }
  }, [messages, visibleCount])

  return {
    visibleMessages: messages.slice(0, visibleCount),
    isCompleted: visibleCount >= messages.length,
    reset() {
      setVisibleCount(0)
    },
  }
}