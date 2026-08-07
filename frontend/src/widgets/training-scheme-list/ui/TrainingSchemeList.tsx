// widgets/training-scheme-list/ui/TrainingSchemeList.tsx

import { useState } from 'react'
import {
  faCheck,
  faChevronDown,
  faLock,
  faPlay,
} from '@fortawesome/free-solid-svg-icons'
import { useNavigate } from 'react-router-dom'

import {
  getScenarioStatus,
  type TrainingScheme,
} from '@/entities/training-path'
import { Icon } from '@/shared/ui/icon'

import styles from './TrainingSchemeList.module.scss'

type TrainingSchemeListProps = {
  schemes: TrainingScheme[]
}

export function TrainingSchemeList({
  schemes,
}: TrainingSchemeListProps) {
  const navigate = useNavigate()

  const [openedSchemeId, setOpenedSchemeId] = useState<
    string | null
  >(schemes[0]?.id ?? null)

  function toggleScheme(schemeId: string) {
    setOpenedSchemeId((currentId) =>
      currentId === schemeId ? null : schemeId,
    )
  }

  return (
    <div className={styles.list}>
      {schemes.map((scheme) => {
        const isOpened = openedSchemeId === scheme.id

        const completedCount = scheme.scenarios.filter(
          (scenario) => scenario.isCompleted,
        ).length

        const totalCount = scheme.scenarios.length

        const progress =
          totalCount === 0
            ? 0
            : Math.round(
                (completedCount / totalCount) * 100,
              )

        const panelId = `scheme-panel-${scheme.id}`

        return (
          <section
            key={scheme.id}
            className={styles.scheme}
          >
            <button
              type="button"
              className={styles.schemeHeader}
              aria-expanded={isOpened}
              aria-controls={panelId}
              onClick={() => toggleScheme(scheme.id)}
            >
              <div className={styles.schemeInformation}>
                <h2 className={styles.schemeTitle}>
                  {scheme.title}
                </h2>

                <p className={styles.schemeDescription}>
                  {scheme.description}
                </p>

                <div className={styles.progressRow}>
                  <div className={styles.progressTrack}>
                    <span
                      className={styles.progressFill}
                      style={{ width: `${progress}%` }}
                    />
                  </div>

                  <span className={styles.progressValue}>
                    {completedCount} из {totalCount}
                  </span>
                </div>
              </div>

              <span
                className={styles.chevron}
                data-opened={isOpened}
                aria-hidden="true"
              >
                <Icon icon={faChevronDown} />
              </span>
            </button>

            {isOpened && (
              <div
                id={panelId}
                className={styles.scenarios}
              >
                {scheme.scenarios.map(
                  (scenario, scenarioIndex) => {
                    const status = getScenarioStatus(
                      scheme.scenarios,
                      scenarioIndex,
                    )

                    const isLocked = status === 'locked'

                    const icon =
                      status === 'completed'
                        ? faCheck
                        : status === 'available'
                          ? faPlay
                          : faLock

                    return (
                      <button
                        key={scenario.id}
                        type="button"
                        className={styles.scenario}
                        data-status={status}
                        disabled={isLocked}
                        onClick={() => {
                          navigate(
                            `/training/path/buyer/${scenario.id}`,
                          )
                        }}
                      >
                        <span
                          className={styles.scenarioIcon}
                          aria-hidden="true"
                        >
                          <Icon icon={icon} />
                        </span>

                        <span
                          className={styles.scenarioContent}
                        >
                          <strong
                            className={styles.scenarioTitle}
                          >
                            {scenario.title}
                          </strong>

                          <span
                            className={
                              styles.scenarioDescription
                            }
                          >
                            {scenario.description}
                          </span>
                        </span>

                        {scenario.durationMinutes && (
                          <span
                            className={
                              styles.scenarioDuration
                            }
                          >
                            {scenario.durationMinutes} мин
                          </span>
                        )}
                      </button>
                    )
                  },
                )}
              </div>
            )}
          </section>
        )
      })}
    </div>
  )
}