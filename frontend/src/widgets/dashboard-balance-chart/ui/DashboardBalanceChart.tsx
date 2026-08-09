import { useMemo } from 'react'
import type { AgChartOptions } from 'ag-charts-community'
import { AgCharts } from 'ag-charts-react'

import { usePuzzleProgress } from '@/entities/profile-progress'

import styles from './DashboardBalanceChart.module.scss'

const CHART_COLOR = '#00aaff'

type PuzzleChartPoint = {
    date: string
    fragments: number
}

const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
})

function createChartData(
    fragments: {
        fragmentId: string
        earnedAt: string
    }[],
): PuzzleChartPoint[] {
    const sortedFragments = [...fragments].sort(
        (first, second) => new Date(first.earnedAt).getTime() - new Date(second.earnedAt).getTime(),
    )

    const fragmentsByDate = new Map<string, number>()

    for (const fragment of sortedFragments) {
        const date = dateFormatter.format(new Date(fragment.earnedAt))

        fragmentsByDate.set(date, (fragmentsByDate.get(date) ?? 0) + 1)
    }

    let total = 0

    return Array.from(fragmentsByDate.entries()).map(([date, count]) => {
        total += count

        return {
            date,
            fragments: total,
        }
    })
}

function createOptions(
    data: PuzzleChartPoint[],
    totalCount: number,
): AgChartOptions<PuzzleChartPoint> {
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

                min: 0,
                max: totalCount > 0 ? totalCount : undefined,
                nice: false,

                title: {
                    enabled: false,
                },

                interval: {
                    step: 1,
                },

                label: {
                    formatter: ({ value }) => `${value}`,
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
                yKey: 'fragments',
                yName: 'Фрагменты',

                interpolation: {
                    type: 'smooth',
                },

                stroke: CHART_COLOR,
                strokeWidth: 3,
                lineDash: [10, 7],

                marker: {
                    enabled: true,
                    size: 7,
                    fill: '#f7f7f7',
                    stroke: CHART_COLOR,
                    strokeWidth: 3,
                },

                tooltip: {
                    renderer: ({ datum }) => ({
                        data: [
                            {
                                label: 'Собрано',
                                value: `${datum.fragments} из ${totalCount}`,
                            },
                        ],
                    }),
                },
            },
        ],
    }
}

export function DashboardBalanceChart() {
    const { data, isPending, isError } = usePuzzleProgress()

    const chartData = useMemo(() => createChartData(data?.fragments ?? []), [data])

    const options = useMemo(
        () => createOptions(chartData, data?.totalCount ?? 0),
        [chartData, data?.totalCount],
    )

    const earnedCount = data?.earnedCount ?? 0
    const totalCount = data?.totalCount ?? 0

    const progress = totalCount > 0 ? Math.round((earnedCount / totalCount) * 100) : 0

    return (
        <section className={styles.card}>
            <header className={styles.header}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Фрагменты пазла</h2>

                    <p className={styles.description}>История сбора фрагментов</p>
                </div>

                {!isPending && !isError && (
                    <div className={styles.currentBalance}>
                        <span className={styles.balance}>
                            {earnedCount} из {totalCount}
                        </span>

                        <span className={styles.improvement}>{progress}%</span>
                    </div>
                )}
            </header>

            {isPending && <p className={styles.description}>Загружаем прогресс...</p>}

            {isError && <p className={styles.description}>Не удалось загрузить прогресс пазла.</p>}

            {!isPending && !isError && chartData.length === 0 && (
                <p className={styles.description}>Фрагменты пока не собраны.</p>
            )}

            {!isPending && !isError && chartData.length > 0 && (
                <div className={styles.chart}>
                    <AgCharts options={options} />
                </div>
            )}
        </section>
    )
}
