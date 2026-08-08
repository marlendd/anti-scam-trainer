// src/widgets/dashboard-path-progress/ui/DashboardPathProgress.tsx

import type { CSSProperties } from 'react'
import { faCartShopping, faStore, type IconDefinition } from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import {
    pathProgressMock,
    type TrainingPathId,
    type TrainingPathProgress,
} from '../model/pathProgress'

import styles from './DashboardPathProgress.module.scss'

type ProgressRingStyle = CSSProperties & {
    '--progress': string
    '--progress-color': string
}

type PathProgressItemProps = {
    path: TrainingPathProgress
}

const pathIcons: Record<TrainingPathId, IconDefinition> = {
    buyer: faCartShopping,
    seller: faStore,
}

function getProgressPercentage(completedScenarios: number, totalScenarios: number) {
    if (totalScenarios <= 0) {
        return 0
    }

    const percentage = (completedScenarios / totalScenarios) * 100

    return Math.min(100, Math.max(0, Math.round(percentage)))
}

function PathProgressItem({ path }: PathProgressItemProps) {
    const progress = getProgressPercentage(path.completedScenarios, path.totalScenarios)

    const ringStyle: ProgressRingStyle = {
        '--progress': `${progress}%`,
        '--progress-color': path.color,
    }

    return (
        <article className={styles.path}>
            <div
                className={styles.progressRing}
                style={ringStyle}
                role="progressbar"
                aria-label={`Прогресс: ${path.title}`}
                aria-valuemin={0}
                aria-valuemax={100}
                aria-valuenow={progress}
            >
                <div className={styles.progressRingInner}>
                    <span className={styles.progressValue}>{progress}%</span>
                </div>
            </div>

            <div className={styles.pathContent}>
                <div className={styles.pathHeading}>
                    <span
                        className={styles.pathIcon}
                        style={{ color: path.color }}
                        aria-hidden="true"
                    >
                        <Icon icon={pathIcons[path.id]} />
                    </span>

                    <h3 className={styles.pathTitle}>{path.title}</h3>
                </div>

                <p className={styles.pathDescription}>{path.description}</p>

                <div className={styles.pathStatistics}>
                    <span>Пройдено</span>

                    <strong>
                        {path.completedScenarios} из {path.totalScenarios}
                    </strong>
                </div>

                <div className={styles.progressTrack} aria-hidden="true">
                    <span
                        className={styles.progressFill}
                        style={{
                            width: `${progress}%`,
                            backgroundColor: path.color,
                        }}
                    />
                </div>
            </div>
        </article>
    )
}

type DashboardPathProgressProps = {
    paths?: TrainingPathProgress[]
}

export function DashboardPathProgress({ paths = pathProgressMock }: DashboardPathProgressProps) {
    return (
        <div className={styles.paths}>
            {paths.map((path) => (
                <PathProgressItem key={path.id} path={path} />
            ))}
        </div>
    )
}
