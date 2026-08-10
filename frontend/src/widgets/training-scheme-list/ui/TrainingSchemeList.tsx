import { useState } from 'react'
import { faCheck, faPlay } from '@fortawesome/free-solid-svg-icons'
import { useNavigate } from 'react-router-dom'

import { startAttempt } from '@/entities/training-attempt'
import type { TrainingScenarioSummary } from '@/entities/training-scenario'
import { Icon } from '@/shared/ui/icon'

import styles from './TrainingSchemeList.module.scss'

type TrainingScenarioListProps = {
    scenarios: TrainingScenarioSummary[]
}

export function TrainingSchemeList({ scenarios }: TrainingScenarioListProps) {
    const navigate = useNavigate()

    const [startingScenarioId, setStartingScenarioId] = useState<string | null>(null)

    async function handleScenarioClick(scenario: TrainingScenarioSummary) {
        if (startingScenarioId) {
            return
        }

        try {
            if (scenario.status !== 'in_progress') {
                setStartingScenarioId(scenario.id)

                await startAttempt(scenario.id)
            }

            navigate(`/training/path/${scenario.role}/${scenario.id}`)
        } finally {
            setStartingScenarioId(null)
        }
    }

    return (
        <div className={styles.list}>
            {scenarios.map((scenario) => {
                const isCompleted = scenario.status === 'completed'

                const isStarting = startingScenarioId === scenario.id

                return (
                    <button
                        key={scenario.id}
                        type="button"
                        className={styles.scenario}
                        data-status={scenario.status}
                        disabled={isStarting}
                        onClick={() => handleScenarioClick(scenario)}
                    >
                        <span className={styles.scenarioIcon} aria-hidden="true">
                            <Icon icon={isCompleted ? faCheck : faPlay} />
                        </span>

                        <span className={styles.scenarioContent}>
                            <strong className={styles.scenarioTitle}>{scenario.title}</strong>

                            <span className={styles.scenarioDescription}>
                                {scenario.description}
                            </span>
                        </span>
                    </button>
                )
            })}
        </div>
    )
}
