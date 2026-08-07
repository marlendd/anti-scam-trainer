// model/scenarioCategoryStats.ts

export type ScenarioCategoryStat = {
  category: string
  completed: number
}

export const scenarioCategoryStats: ScenarioCategoryStat[] = [
  {
    category: 'Фишинг',
    completed: 4,
  },
  {
    category: 'Предоплата',
    completed: 3,
  },
  {
    category: 'Доставка',
    completed: 2,
  },
  {
    category: 'СМС-коды',
    completed: 1,
  },
  {
    category: 'Поддержка',
    completed: 3,
  },
]