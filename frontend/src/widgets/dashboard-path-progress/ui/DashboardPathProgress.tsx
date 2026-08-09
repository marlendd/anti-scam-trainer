import type { CSSProperties } from 'react'
import { faCartShopping, faStore, type IconDefinition } from '@fortawesome/free-solid-svg-icons'

import { type ProgressRole, useRoleProgress } from '@/entities/profile-progress'
import { useScenarios } from '@/entities/training-scenario'
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
                        style={{
                            color: path.color,
                        }}
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

function getCompletedCount(
    roleProgress:
        | {
              role: ProgressRole
              completedCount: number
          }[]
        | undefined,
    role: ProgressRole,
) {
    return roleProgress?.find((item) => item.role === role)?.completedCount ?? 0
}

export function DashboardPathProgress() {
    const {
        data: roleProgress,
        isPending: isRoleProgressPending,
        isError: isRoleProgressError,
    } = useRoleProgress()

    const {
        data: buyerScenarios = [],
        isPending: isBuyerPending,
        isError: isBuyerError,
    } = useScenarios('buyer')

    const {
        data: sellerScenarios = [],
        isPending: isSellerPending,
        isError: isSellerError,
    } = useScenarios('seller')

    const isPending = isRoleProgressPending || isBuyerPending || isSellerPending

    const isError = isRoleProgressError || isBuyerError || isSellerError

    if (isPending) {
        return null
    }

    if (isError) {
        return null
    }

    const paths = pathProgressMock.map((path) => {
        const role = path.id as ProgressRole

        const totalScenarios = role === 'buyer' ? buyerScenarios.length : sellerScenarios.length

        return {
            ...path,
            completedScenarios: getCompletedCount(roleProgress, role),
            totalScenarios,
        }
    })

    return (
        <div className={styles.paths}>
            {paths.map((path) => (
                <PathProgressItem key={path.id} path={path} />
            ))}
        </div>
    )
}
