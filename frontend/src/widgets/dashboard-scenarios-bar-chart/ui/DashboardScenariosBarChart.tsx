import { useMemo } from 'react'
import type { AgChartOptions } from 'ag-charts-community'
import { AgCharts } from 'ag-charts-react'

import { useCategoriesProgress } from '@/entities/profile-progress'

import styles from './DashboardScenariosBarChart.module.scss'

type CategoryChartPoint = {
    category: string
    count: number
}

function createOptions(data: CategoryChartPoint[]): AgChartOptions<CategoryChartPoint> {
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

                tick: {
                    enabled: false,
                },

                label: {
                    color: '#858585',
                    fontSize: 12,
                    autoRotate: false,
                },

                gridLine: {
                    enabled: false,
                },
            },

            y: {
                type: 'number',
                position: 'left',
                min: 0,
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
                type: 'bar',
                xKey: 'category',
                yKey: 'count',
                yName: 'Количество',

                fill: '#00aaff',
                stroke: '#00aaff',
                strokeWidth: 0,

                cornerRadius: 12,

                tooltip: {
                    renderer: ({ datum }) => ({
                        heading: datum.category,
                        data: [
                            {
                                label: 'Количество',
                                value: datum.count,
                            },
                        ],
                    }),
                },
            },
        ],
    }
}

export function DashboardScenariosBarChart() {
    const { data, isPending, isError } = useCategoriesProgress()

    const chartData = data?.stats ?? []

    const options = useMemo(() => createOptions(chartData), [chartData])

    return (
        <section className={styles.card}>
            <header className={styles.header}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Категории риска</h2>

                    <p className={styles.description}>
                        Уязвимости по категориям мошеннических сценариев
                    </p>
                </div>

                {data && (
                    <div className={styles.currentValue}>
                        <span className={styles.value}>{data.totalCompleted}</span>

                        <span className={styles.valueLabel}>пройдено</span>
                    </div>
                )}
            </header>

            {isPending && <p className={styles.description}>Загружаем статистику...</p>}

            {isError && <p className={styles.description}>Не удалось загрузить статистику.</p>}

            {!isPending && !isError && chartData.length === 0 && (
                <p className={styles.description}>Статистики пока нет.</p>
            )}

            {!isPending && !isError && chartData.length > 0 && (
                <div className={styles.chart}>
                    <AgCharts options={options} />
                </div>
            )}
        </section>
    )
}
