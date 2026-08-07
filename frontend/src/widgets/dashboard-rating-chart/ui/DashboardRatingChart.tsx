import { useMemo } from 'react'
import type { AgChartOptions } from 'ag-charts-community'
import { AgCharts } from 'ag-charts-react'

import { rankingHistory } from '../model/rankingHistory'

import styles from './DashboardRatingChart.module.scss'

const IMPROVEMENT_COLOR = '#04e061'
const DECLINE_COLOR = '#ff4053'
const NEUTRAL_COLOR = '#858585'

function createOptions(chartColor: string): AgChartOptions {
  return {
    data: rankingHistory,

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
  const currentRank =
    rankingHistory[rankingHistory.length - 1]?.rank

  const previousRank =
    rankingHistory[rankingHistory.length - 2]?.rank

  const rankDifference =
    previousRank !== undefined && currentRank !== undefined
      ? previousRank - currentRank
      : 0

  const chartColor =
    rankDifference > 0
      ? IMPROVEMENT_COLOR
      : rankDifference < 0
        ? DECLINE_COLOR
        : NEUTRAL_COLOR

  const options = useMemo(
    () => createOptions(chartColor),
    [chartColor],
  )

  return (
    <section className={styles.card}>
      <header className={styles.header}>
        <div className={styles.text}>
          <h2 className={styles.title}>
            Позиция в рейтинге
          </h2>

          <p className={styles.description}>
            Чем выше линия, тем лучше результат
          </p>
        </div>

        <div className={styles.currentRank}>
          <span className={styles.rank}>
            {currentRank}
          </span>

          {rankDifference !== 0 && (
            <span
              className={
                rankDifference > 0
                  ? styles.improvement
                  : styles.decline
              }
            >
              {rankDifference > 0 ? '↑' : '↓'}{' '}
              {Math.abs(rankDifference)}
            </span>
          )}
        </div>
      </header>

      <div className={styles.chart}>
        <AgCharts options={options} />
      </div>
    </section>
  )
}