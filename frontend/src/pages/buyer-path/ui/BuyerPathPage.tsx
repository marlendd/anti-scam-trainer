import { useMatch } from 'react-router-dom'

import { buyerSchemesMock } from '@/entities/training-path'
import {
  fakeDeliveryScenario,
  type TrainingScenario,
} from '@/entities/training-scenario'
import { useDocumentTitle } from '@/shared/lib/use-document-title'
import { TrainingChat } from '@/widgets/training-chat'
import { TrainingSchemeList } from '@/widgets/training-scheme-list'

import styles from './BuyerPathPage.module.scss'

export function BuyerPathPage() {
  const scenario: TrainingScenario =
    fakeDeliveryScenario

  const sessionMatch = useMatch({
    path: '/training/path/:pathId/:schemeId',
    end: true,
  })

  const isSession = Boolean(sessionMatch)

  useDocumentTitle(scenario.title)

  return (
    <main
      className={styles.page}
      data-session={isSession}
    >
      <div className={styles.pathColumn}>
        <TrainingSchemeList
          schemes={buyerSchemesMock}
        />
      </div>

      <div className={styles.chatColumn}>
        <TrainingChat
          mode="interactive"
          scenario={isSession ? scenario : null}
          onScenarioComplete={() => {}}
        />
      </div>
    </main>
  )
}