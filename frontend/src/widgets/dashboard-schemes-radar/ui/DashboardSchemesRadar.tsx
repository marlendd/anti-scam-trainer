import { useMemo } from 'react'
import type { AgChartOptions } from 'ag-charts-enterprise'
import { AgCharts } from 'ag-charts-react'

import {
  scamSchemeChartData,
  type ScamSchemeChartDatum,
} from '../model/scamSchemeRadarData.ts'

import styles from './DashboardSchemesRadar.module.scss'

export function DashboardSchemesRadar() {
  const averageAccuracy = useMemo(() => {
    if (scamSchemeChartData.length === 0) {
      return null
    }

    const totalAccuracy = scamSchemeChartData.reduce(
      (sum, item) => sum + item.accuracy,
      0,
    )

    return Math.round(totalAccuracy / scamSchemeChartData.length)
  }, [])

  const options = useMemo<AgChartOptions<ScamSchemeChartDatum>>(
    () => ({
      data: scamSchemeChartData,

      background: {
        fill: 'transparent',
      },

      padding: {
        top: 48,
        right: 100,
        bottom: 48,
        left: 100,
      },

      legend: {
        enabled: false,
      },

      axes: {
        angle: {
          type: 'angle-category',

          paddingInner: 0.12,

          label: {
            orientation: 'fixed',
            fontSize: 12,
            color: '#5c5c5c',
            padding: 14,
          },

          gridLine: {
            enabled: false,
          },
        },

        radius: {
          type: 'radius-number',

          min: 0,
          max: 100,
          nice: false,

          label: {
            fontSize: 11,
            color: '#8c8c8c',
            formatter: ({ value }) => `${value}%`,
          },

          gridLine: {
            enabled: true,
            style: [
              {
                stroke: '#e5e7eb',
                lineDash: [4, 4],
              },
            ],
          },
        },
      },

      series: [
        {
          type: 'nightingale',

          angleKey: 'scheme',
          radiusKey: 'accuracy',

          angleName: 'Мошенническая схема',
          radiusName: 'Точность',

          cornerRadius: 8,

          fillOpacity: 0.86,
          strokeWidth: 1,

          itemStyler: ({ datum }) => ({
            fill: datum.color,
            stroke: datum.color,
          }),

          tooltip: {
            renderer: ({ datum }) => ({
              heading: datum.fullTitle,
              data: [
                {
                  label: 'Точность',
                  value: `${datum.accuracy}%`,
                },
              ],
            }),
          },
        },
      ],
    }),
    [],
  )

  if (scamSchemeChartData.length === 0) {
    return (
      <section className={styles.card}>
        <h2 className={styles.title}>
          Распознавание мошеннических схем
        </h2>

        <p className={styles.empty}>
          Недостаточно данных для построения диаграммы
        </p>
      </section>
    )
  }

  return (
    <section className={styles.card}>
      <header className={styles.header}>
        <div>
          <h2 className={styles.title}>
            Распознавание схем
          </h2>

          <p className={styles.description}>
            Показывает точность распознавания
            конкретной схемы
          </p>
        </div>

        {averageAccuracy !== null && (
          <div className={styles.average}>
            <span className={styles.averageLabel}>
              Средняя точность
            </span>

            <span className={styles.averageValue}>
              {averageAccuracy}%
            </span>
          </div>
        )}
      </header>

      <div className={styles.chart}>
        <AgCharts options={options} />
      </div>
    </section>
  )
}