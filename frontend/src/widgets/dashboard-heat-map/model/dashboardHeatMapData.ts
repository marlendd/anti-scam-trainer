// src/widgets/dashboard-heatmap/model/dashboardHeatMapData.ts

export type DashboardHeatMapDatum = {
  date: string
  dateLabel: string
  week: string
  weekday: string
  activity: number
}

const WEEKDAYS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const weekFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'short',
})

const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})

function formatDateKey(date: Date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  return `${year}-${month}-${day}`
}

function startOfWeek(date: Date) {
  const result = new Date(
    date.getFullYear(),
    date.getMonth(),
    date.getDate(),
  )

  const weekday = result.getDay()
  const daysSinceMonday = weekday === 0 ? 6 : weekday - 1

  result.setDate(result.getDate() - daysSinceMonday)

  return result
}

function addDays(date: Date, days: number) {
  const result = new Date(date)
  result.setDate(result.getDate() + days)

  return result
}

function seededRandom(seed: number) {
  const value = Math.sin(seed) * 10_000

  return value - Math.floor(value)
}

/**
 * Псевдослучайный детерминированный mock.
 * Для одной даты всегда возвращается одинаковое значение.
 */
function getMockActivity(date: Date) {
  const seed =
    date.getFullYear() * 10_000 +
    (date.getMonth() + 1) * 100 +
    date.getDate()

  const activityChance = seededRandom(seed)
  const activityValue = seededRandom(seed + 42)

  const isWeekend = date.getDay() === 0 || date.getDay() === 6
  const emptyDayChance = isWeekend ? 0.5 : 0.3

  if (activityChance < emptyDayChance) {
    return 0
  }

  return Math.floor(activityValue * 8) + 1
}

export function createdashboardHeatMapData(
  weeks = 20,
  today = new Date(),
): DashboardHeatMapDatum[] {
  const currentDate = new Date(
    today.getFullYear(),
    today.getMonth(),
    today.getDate(),
  )

  const currentWeekStart = startOfWeek(currentDate)

  const rangeStart = addDays(
    currentWeekStart,
    -(weeks - 1) * 7,
  )

  const data: DashboardHeatMapDatum[] = []

  for (
    let date = rangeStart;
    date <= currentDate;
    date = addDays(date, 1)
  ) {
    const monday = startOfWeek(date)

    const weekdayIndex =
      date.getDay() === 0 ? 6 : date.getDay() - 1

    data.push({
      date: formatDateKey(date),
      dateLabel: dateFormatter.format(date),
      week: weekFormatter.format(monday),
      weekday: WEEKDAYS[weekdayIndex],
      activity: getMockActivity(date),
    })
  }

  return data
}

export const dashboardHeatMapData =
  createdashboardHeatMapData()