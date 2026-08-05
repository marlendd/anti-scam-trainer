// src/widgets/dashboard-heatmap/ui/DashboardHeatmap.tsx

import {useMemo} from 'react'
import type {AgChartOptions} from 'ag-charts-enterprise'
import {AgCharts} from 'ag-charts-react'

import {
    dashboardHeatMapData,
    type DashboardHeatMapDatum,
} from '../model/dashboardHeatMapData'

import styles from './DashboardHeatMap.module.scss'

export function DashboardHeatMap() {
    const statistics = useMemo(() => {
        return dashboardHeatMapData.reduce(
            (result, item) => {
                result.completedScenarios += item.activity

                if (item.activity > 0) {
                    result.activeDays += 1
                }

                return result
            },
            {
                completedScenarios: 0,
                activeDays: 0,
            },
        )
    }, [])

    const options = useMemo<
        AgChartOptions<DashboardHeatMapDatum>
    >(
        () => ({
            data: dashboardHeatMapData,

            background: {
                fill: 'transparent',
            },

            padding: {
                top: 12,
                right: 12,
                bottom: 8,
                left: 8,
            },

            legend: {
                enabled: false,
            },

            gradientLegend: {
                enabled: false,
            },

            axes: {
                x: {
                    position: 'bottom',

                    interval: {
                        step: 2,
                    },

                    tick: {
                        enabled: false,
                    },

                    label: {
                        fontSize: 11,
                        color: '#858585',
                        autoRotate: false,
                    },

                    gridLine: {
                        enabled: false,
                    },
                },

                y: {
                    position: 'left',
                    reverse: true,

                    tick: {
                        enabled: false,
                    },

                    label: {
                        fontSize: 11,
                        color: '#858585',
                    },

                    gridLine: {
                        enabled: false,
                    },
                },
            },

            series: [
                {
                    type: 'heatmap',

                    xKey: 'week',
                    yKey: 'weekday',
                    colorKey: 'activity',

                    xName: 'Неделя',
                    yName: 'День недели',
                    colorName: 'Пройдено сценариев',

                    itemPadding: 3,

                    stroke: '#f7f7f7',
                    strokeWidth: 2,

                    colorScale: {
                        mode: 'discrete',
                        domain: [0, 8],

                        fills: [
                            {
                                color: '#fff1f2',
                                stop: 1,
                            },
                            {
                                color: '#ffd9dd',
                                stop: 3,
                            },
                            {
                                color: '#ffadb5',
                                stop: 5,
                            },
                            {
                                color: '#ff7380',
                                stop: 7,
                            },
                            {
                                color: '#ff4053',
                            },
                        ],
                    },

                    label: {
                        enabled: false,
                    },

                    tooltip: {
                        renderer: ({datum}) => ({
                            heading: datum.dateLabel,
                            data: [
                                {
                                    label: 'Пройдено сценариев',
                                    value: String(datum.activity),
                                },
                            ],
                        }),
                    },
                },
            ],
        }),
        [],
    )

    if (dashboardHeatMapData.length === 0) {
        return (
            <section className={styles.card}>
                <h2 className={styles.title}>
                    Активность тренировок
                </h2>

                <p className={styles.empty}>
                    Данных об активности пока нет
                </p>
            </section>
        )
    }

    return (
        <section className={styles.card}>
            <header className={styles.header}>
                <div>
                    <h2 className={styles.title}>
                        Активность тренировок
                    </h2>

                    <p className={styles.description}>
                        Ваша активность по дням
                    </p>
                </div>

                <div className={styles.statistics}>
                    <div className={styles.statistic}>
            <span className={styles.statisticValue}>
              {statistics.completedScenarios}
            </span>

                        <span className={styles.statisticLabel}>
              сценариев
            </span>
                    </div>

                    <div className={styles.statistic}>
            <span className={styles.statisticValue}>
              {statistics.activeDays}
            </span>

                        <span className={styles.statisticLabel}>
              активных дней
            </span>
                    </div>
                </div>
            </header>

            <div className={styles.chartViewport}>
                <div className={styles.chart}>
                    <AgCharts options={options}/>
                </div>
            </div>

        </section>
    )
}