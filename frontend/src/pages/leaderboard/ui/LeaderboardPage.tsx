import { useState } from 'react'
import type { ColDef } from 'ag-grid-community'
import { AgGridReact } from 'ag-grid-react'

import { appGridTheme } from '@/shared/config/ad-grid/theme'
import { useDocumentTitle } from '@/shared/lib/use-document-title'

import styles from './LeaderboardPage.module.scss'

type LeaderboardRow = {
    rank: number
    player: string
    fragments: number
    points: number
}

type StatisticsPeriod = 'day' | 'week' | 'month' | 'all'

type PeriodOption = {
    value: StatisticsPeriod
    label: string
}

const periodOptions: PeriodOption[] = [
    { value: 'day', label: 'За день' },
    { value: 'week', label: 'За неделю' },
    { value: 'month', label: 'За месяц' },
    { value: 'all', label: 'За всё время' },
]

const rowData: LeaderboardRow[] = [
    { rank: 1, player: 'Mia', fragments: 42, points: 980 },
    { rank: 2, player: 'Alex', fragments: 38, points: 910 },
    { rank: 3, player: 'Nora', fragments: 31, points: 845 },
    { rank: 4, player: 'Leo', fragments: 29, points: 790 },
    { rank: 5, player: 'Zoe', fragments: 24, points: 720 },
    { rank: 6, player: 'Oliver', fragments: 23, points: 698 },
    { rank: 7, player: 'Emma', fragments: 22, points: 674 },
    { rank: 8, player: 'Noah', fragments: 21, points: 651 },
    { rank: 9, player: 'Sophia', fragments: 20, points: 632 },
    { rank: 10, player: 'Liam', fragments: 19, points: 608 },
    { rank: 11, player: 'Ava', fragments: 18, points: 587 },
    { rank: 12, player: 'Lucas', fragments: 18, points: 566 },
    { rank: 13, player: 'Amelia', fragments: 17, points: 548 },
    { rank: 14, player: 'Ethan', fragments: 16, points: 529 },
    { rank: 15, player: 'Isabella', fragments: 16, points: 511 },
    { rank: 16, player: 'Mason', fragments: 15, points: 493 },
    { rank: 17, player: 'Harper', fragments: 15, points: 478 },
    { rank: 18, player: 'James', fragments: 14, points: 462 },
    { rank: 19, player: 'Evelyn', fragments: 14, points: 447 },
    { rank: 20, player: 'Henry', fragments: 13, points: 431 },
    { rank: 21, player: 'Ella', fragments: 13, points: 415 },
    { rank: 22, player: 'Benjamin', fragments: 12, points: 398 },
    { rank: 23, player: 'Scarlett', fragments: 12, points: 382 },
    { rank: 24, player: 'Jack', fragments: 11, points: 367 },
    { rank: 25, player: 'Grace', fragments: 11, points: 351 },
    { rank: 26, player: 'Daniel', fragments: 10, points: 336 },
    { rank: 27, player: 'Chloe', fragments: 10, points: 321 },
    { rank: 28, player: 'Michael', fragments: 9, points: 305 },
    { rank: 29, player: 'Victoria', fragments: 9, points: 289 },
    { rank: 30, player: 'Sebastian', fragments: 8, points: 274 },
    { rank: 31, player: 'Lily', fragments: 8, points: 259 },
    { rank: 32, player: 'Samuel', fragments: 7, points: 243 },
    { rank: 33, player: 'Hannah', fragments: 7, points: 228 },
    { rank: 34, player: 'David', fragments: 6, points: 212 },
    { rank: 35, player: 'Aria', fragments: 6, points: 197 },
    { rank: 36, player: 'Joseph', fragments: 5, points: 181 },
    { rank: 37, player: 'Layla', fragments: 5, points: 166 },
    { rank: 38, player: 'Matthew', fragments: 4, points: 151 },
    { rank: 39, player: 'Penelope', fragments: 4, points: 135 },
    { rank: 40, player: 'Andrew', fragments: 3, points: 120 },
]

const columnDefs: ColDef<LeaderboardRow>[] = [
    {
        field: 'rank',
        headerName: 'Место',
        width: 110,
    },
    {
        field: 'player',
        headerName: 'Игрок',
        flex: 1,
        minWidth: 180,
    },
    {
        field: 'fragments',
        headerName: 'Фрагменты',
        width: 160,
    },
    {
        field: 'points',
        headerName: 'Баллы',
        width: 140,
    },
]

const defaultColDef: ColDef<LeaderboardRow> = {
    resizable: false,
    sortable: true,
    suppressMovable: true,
}

export function LeaderboardPage() {
    const [period, setPeriod] = useState<StatisticsPeriod>('week')

    useDocumentTitle('Таблица лидеров')

    return (
        <main className={styles.page}>
            <section className={styles.card}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Таблица лидеров</h2>
                    <span className={styles.counter}>40</span>

                    <div className={styles.periodSelect}>
                        <select
                            className={styles.periodSelectControl}
                            value={period}
                            aria-label="Период статистики"
                            onChange={(event) => {
                                setPeriod(event.target.value as StatisticsPeriod)
                            }}
                        >
                            {periodOptions.map((option) => (
                                <option key={option.value} value={option.value}>
                                    {option.label}
                                </option>
                            ))}
                        </select>

                        <span className={styles.periodSelectChevron} aria-hidden="true" />
                    </div>
                </div>

                <div className={`ag-theme-quartz ${styles.grid}`}>
                    <AgGridReact<LeaderboardRow>
                        theme={appGridTheme}
                        rowData={rowData}
                        columnDefs={columnDefs}
                        defaultColDef={defaultColDef}
                        domLayout="autoHeight"
                        headerHeight={48}
                        rowHeight={56}
                        suppressMovableColumns
                        suppressCellFocus
                        pagination
                        paginationPageSize={9}
                        paginationPageSizeSelector={[10, 25, 50]}
                    />
                </div>
            </section>
        </main>
    )
}
