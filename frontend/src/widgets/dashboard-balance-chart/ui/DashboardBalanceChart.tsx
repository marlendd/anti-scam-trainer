import {useMemo} from 'react'

import {usePuzzleProgress} from '@/entities/profile-progress'
import {
    PuzzleBoard,
    type PuzzlePieceId,
} from '@/entities/puzzle'

import PuzzleImage from '@/shared/assets/images/puzzle.webp'

import styles from './DashboardBalanceChart.module.scss'


export function DashboardBalanceChart() {
    const {
        data,
        isPending,
        isError,
    } = usePuzzleProgress()

    const unlockedPieces = useMemo(
        () =>
            Array.from(
                {
                    length: Math.min(
                        data?.earnedCount ?? 0,
                        9,
                    ),
                },
                (_, index) => (index + 1) as PuzzlePieceId,
            ),
        [data?.earnedCount],
    )

    const earnedCount = data?.earnedCount ?? 0
    const totalCount = data?.totalCount ?? 0

    const progress =
        totalCount > 0
            ? Math.round(
                (earnedCount / totalCount) * 100,
            )
            : 0

    return (
        <section className={styles.card}>
            <header className={styles.header}>
                <div className={styles.text}>
                    <h2 className={styles.title}>
                        Фрагменты пазла
                    </h2>

                    <p className={styles.description}>
                        Собранные фрагменты
                    </p>
                </div>

                {!isPending && !isError && (
                    <div className={styles.currentBalance}>
                        <span className={styles.balance}>
                            {earnedCount} из {totalCount}
                        </span>

                        <span className={styles.improvement}>
                            {progress}%
                        </span>
                    </div>
                )}
            </header>

            {isPending && (
                <p className={styles.description}>
                    Загружаем прогресс...
                </p>
            )}

            {isError && (
                <p className={styles.description}>
                    Не удалось загрузить прогресс пазла.
                </p>
            )}

            {}

            {!isPending && !isError && (
                <div className={styles.chart}>
                    <PuzzleBoard
                        imageSrc={PuzzleImage}
                        unlockedPieces={unlockedPieces}
                    />
                </div>
            )}
        </section>
    )
}