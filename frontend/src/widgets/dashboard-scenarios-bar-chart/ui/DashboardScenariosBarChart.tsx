// ui/DashboardScenariosBarChart.tsx

import type { AgChartOptions } from 'ag-charts-community'
import { AgCharts } from 'ag-charts-react'

import { scenarioCategoryStats, type ScenarioCategoryStat } from '../model/scenarioCategoryStats'

import styles from './DashboardScenariosBarChart.module.scss'

const options: AgChartOptions<ScenarioCategoryStat> = {
    data: scenarioCategoryStats,

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
            max: 4,
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
            yKey: 'completed',
            yName: 'Пройдено сценариев',

            fill: '#00aaff',
            stroke: '#00aaff',
            strokeWidth: 0,

            cornerRadius: 12,

            tooltip: {
                renderer: ({ datum }) => ({
                    heading: datum.category,
                    data: [
                        {
                            label: 'Пройдено',
                            value: `${datum.completed} из 4`,
                        },
                    ],
                }),
            },
        },
    ],
}

export function DashboardScenariosBarChart() {
    const completedScenarios = scenarioCategoryStats.reduce(
        (total, category) => total + category.completed,
        0,
    )

    return (
        <section className={styles.card}>
            <header className={styles.header}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Пройденные сценарии</h2>

                    <p className={styles.description}>Количество тренировок по категориям</p>
                </div>

                <div className={styles.currentValue}>
                    <span className={styles.value}>{completedScenarios}</span>

                    <span className={styles.valueLabel}>всего</span>
                </div>
            </header>

            <div className={styles.chart}>
                <AgCharts options={options} />
            </div>
        </section>
    )
}
