import { useMemo } from 'react'
import type { AgChartOptions } from 'ag-charts-community'
import { AgCharts } from 'ag-charts-react'

import { useRankHistory } from '@/entities/profile-progress'

import styles from './DashboardRatingChart.module.scss'

const IMPROVEMENT_COLOR = '#04e061'
const DECLINE_COLOR = '#ff4053'
const NEUTRAL_COLOR = '#858585'

const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
})

type RatingChartPoint = {
    date: string
    rank: number
}

function createOptions(
    data: RatingChartPoint[],
    chartColor: string,
): AgChartOptions<RatingChartPoint> {
    return {
        data,

        background: {
            fill: 'transparent',
        },

        padding: {
            top: 16,
            right: 0,
            bottom: 0,
            left: 0,
        },

        legend: {
            enabled: false,
        },

        axes: {
            x: {
                type: 'category',
                position: 'bottom',

                title: {
                    enabled: false,
                },

                gridLine: {
                    enabled: false,
                },
            },

            y: {
                type: 'number',
                position: 'left',
                reverse: true,
                nice: false,

                label: {
                    formatter: ({ value }) => `#${value}`,
                },

                gridLine: {
                    enabled: true,
                    style: [
                        {
                            lineDash: [4, 6],
                        },
                    ],
                },
            },
        },

        series: [
            {
                type: 'line',
                xKey: 'date',
                yKey: 'rank',
                yName: 'Место в рейтинге',

                interpolation: {
                    type: 'smooth',
                },

                stroke: chartColor,
                strokeWidth: 3,
                lineDash: [10, 7],

                marker: {
                    enabled: true,
                    size: 7,
                    fill: '#f7f7f7',
                    stroke: chartColor,
                    strokeWidth: 3,
                },

                tooltip: {
                    renderer: ({ datum }) => ({
                        data: [
                            {
                                label: 'Место',
                                value: `#${datum.rank}`,
                            },
                        ],
                    }),
                },
            },
        ],
    }
}

export function DashboardRatingChart() {
    const { data, isPending, isError } = useRankHistory()

    const rankingHistory = useMemo(
        () =>
            (data?.history ?? []).map((point) => ({
                date: dateFormatter.format(new Date(point.date)),
                rank: point.rank,
            })),
        [data],
    )

    const currentRank = rankingHistory[rankingHistory.length - 1]?.rank

    const previousRank = rankingHistory[rankingHistory.length - 2]?.rank

    const rankDifference =
        previousRank !== undefined && currentRank !== undefined ? previousRank - currentRank : 0

    const chartColor =
        rankDifference > 0 ? IMPROVEMENT_COLOR : rankDifference < 0 ? DECLINE_COLOR : NEUTRAL_COLOR

    const options = useMemo(
        () => createOptions(rankingHistory, chartColor),
        [rankingHistory, chartColor],
    )

    return (
        <section className={styles.card}>
            <header className={styles.header}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Позиция в рейтинге</h2>

                    <p className={styles.description}>Чем выше линия, тем лучше результат</p>
                </div>

                {currentRank !== undefined && (
                    <div className={styles.currentRank}>
                        <span className={styles.rank}>{currentRank}</span>

                        {rankDifference !== 0 && (
                            <span
                                className={rankDifference > 0 ? styles.improvement : styles.decline}
                            >
                                {rankDifference > 0 ? '↑' : '↓'} {Math.abs(rankDifference)}
                            </span>
                        )}
                    </div>
                )}
            </header>

            {isPending && <p className={styles.description}>Загружаем историю рейтинга...</p>}

            {isError && (
                <p className={styles.description}>Не удалось загрузить историю рейтинга.</p>
            )}

            {!isPending && !isError && rankingHistory.length === 0 && (
                <p className={styles.description}>Истории рейтинга пока нет.</p>
            )}

            {!isPending && !isError && rankingHistory.length > 0 && (
                <div className={styles.chart}>
                    <AgCharts options={options} />
                </div>
            )}
        </section>
    )
}
