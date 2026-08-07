import { useEffect, useState } from 'react'

import {
  fakeDeliveryScenario,
  type TrainingScenario,
} from '@/entities/training-scenario'
import { useDocumentTitle } from '@/shared/lib/use-document-title'
import { ScenarioAnalysisPanel } from '@/widgets/scenario-analysis-panel'
import { TrainingChat } from '@/widgets/training-chat'

import styles from './ScamSchemePage.module.scss'

export function ScamSchemePage() {
  const [scenario] = useState<TrainingScenario>(fakeDeliveryScenario)
  const [isScenarioCompleted, setIsScenarioCompleted] = useState(false)

  useDocumentTitle(scenario.title)

  useEffect(() => {
    setIsScenarioCompleted(true)
  }, [scenario.id])

  return (
    <main className={styles.page}>
      <section className={styles.header}>
        <div className={styles.text}>
          <h2 className={styles.title}>{scenario.title}</h2>
          <span className={styles.description}>{scenario.description}</span>
        </div>
      </section>

      <div className={styles.training}>
        <TrainingChat
          scenario={scenario}
          mode="preview"
          onScenarioComplete={() => {
            setIsScenarioCompleted(true)
          }}
        />

        <ScenarioAnalysisPanel
          analysis={scenario.analysis}
          revealed={isScenarioCompleted}
        />
      </div>
    </main>
  )
}