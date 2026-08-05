// model/scamSchemeChartData.ts

export type ScamSchemeChartDatum = {
  scheme: string
  fullTitle: string
  accuracy: number
  color: string
}

export const scamSchemeChartData: ScamSchemeChartDatum[] = [
  {
    scheme: 'Фишинговая ссылка',
    fullTitle: 'Фишинговая ссылка',
    accuracy: 82,
    color: '#965eeb',
  },
  {
    scheme: 'Предоплата',
    fullTitle: 'Требование предоплаты',
    accuracy: 74,
    color: '#00aaff',
  },
  {
    scheme: 'Поддельная доставка',
    fullTitle: 'Поддельная доставка',
    accuracy: 68,
    color: '#ff6163',
  },
  {
    scheme: 'Код из СМС',
    fullTitle: 'Запрос кода из СМС',
    accuracy: 91,
    color: '#00c27a',
  },
  {
    scheme: 'Ложная поддержка',
    fullTitle: 'Ложная служба поддержки',
    accuracy: 63,
    color: '#ffb020',
  },
  {
    scheme: 'Другой мессенджер',
    fullTitle: 'Переход в другой мессенджер',
    accuracy: 78,
    color: '#8c6cff',
  },
  {
    scheme: 'Срочность',
    fullTitle: 'Искусственная срочность',
    accuracy: 71,
    color: '#00b8d9',
  },
  {
    scheme: 'Лишний перевод',
    fullTitle: 'Лишний перевод',
    accuracy: 58,
    color: '#f05fa6',
  },
  {
    scheme: 'Подмена товара',
    fullTitle: 'Подмена товара',
    accuracy: 66,
    color: '#7bbf44',
  },
]