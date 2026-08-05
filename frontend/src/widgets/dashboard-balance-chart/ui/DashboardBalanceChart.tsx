import type { AgChartOptions } from 'ag-charts-community'
import { AgCharts } from 'ag-charts-react'

import { balanceHistory } from '../model/balanceHistory'

import styles from './DashboardBalanceChart.module.scss'

const INCREASE_COLOR = '#04e061'
const DECREASE_COLOR = '#ff4053'
const NEUTRAL_COLOR = '#858585'

export function DashboardBalanceChart() {
  const currentBalance =
    balanceHistory[balanceHistory.length - 1]?.balance

  const previousBalance =
    balanceHistory[balanceHistory.length - 2]?.balance

  const balanceDifference =
    previousBalance !== undefined && currentBalance !== undefined
      ? currentBalance - previousBalance
      : 0

  const chartColor =
    balanceDifference > 0
      ? INCREASE_COLOR
      : balanceDifference < 0
        ? DECREASE_COLOR
        : NEUTRAL_COLOR

  const options: AgChartOptions = {
    data: balanceHistory,

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
        nice: false,

        label: {
          formatter: ({ value }) => `${value} ₽`,
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
        yKey: 'balance',
        yName: 'Баланс',

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
                label: 'Остаток',
                value: `${datum.balance} ₽`,
              },
            ],
          }),
        },
      },
    ],
  }

  return (
    <section className={styles.card}>
      <header className={styles.header}>
        <div className={styles.text}>
          <h2 className={styles.title}>Ваш баланс</h2>

          <p className={styles.description}>
            История изменения баланса
          </p>
        </div>

        <div className={styles.currentBalance}>
          <span className={styles.balance}>
            {currentBalance} ₽
          </span>

          {balanceDifference !== 0 && (
            <span
              className={
                balanceDifference > 0
                  ? styles.improvement
                  : styles.decline
              }
            >
              {balanceDifference > 0 ? '↑' : '↓'}{' '}
              {Math.abs(balanceDifference)} ₽
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