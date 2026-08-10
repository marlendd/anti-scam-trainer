import type { ColDef } from 'ag-grid-community'
import { AgGridReact } from 'ag-grid-react'

import { type LeaderboardEntry, useLeaderboard } from '@/entities/leaderboard'
import { appGridTheme } from '@/shared/config/ad-grid/theme'
import { useDocumentTitle } from '@/shared/lib/use-document-title'

import styles from './LeaderboardPage.module.scss'

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

const columnDefs: ColDef<LeaderboardEntry>[] = [
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
        field: 'score',
        headerName: 'Баллы',
        width: 140,
    },
]

const defaultColDef: ColDef<LeaderboardEntry> = {
    resizable: false,
    sortable: true,
    suppressMovable: true,
}

export function LeaderboardPage() {
    const { data, isPending, isError } = useLeaderboard({
        limit: 50,
        offset: 0,
    })

    useDocumentTitle('Таблица лидеров')

    const rowData = data?.entries ?? []

    return (
        <main className={styles.page}>
            <section className={styles.card}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Таблица лидеров</h2>

                    {!isPending && !isError && (
                        <span className={styles.counter}>{rowData.length}</span>
                    )}

                    <div className={styles.periodSelect}>
                        <select
                            className={styles.periodSelectControl}
                            value="all"
                            aria-label="Период статистики"
                            disabled
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

                {isPending && <div className={styles.state}>Загружаем рейтинг...</div>}

                {isError && (
                    <div className={styles.state}>Не удалось загрузить таблицу лидеров.</div>
                )}

                {!isPending && !isError && rowData.length === 0 && (
                    <div className={styles.state}>В таблице лидеров пока никого нет.</div>
                )}

                {!isPending && !isError && rowData.length > 0 && (
                    <div className={`ag-theme-quartz ${styles.grid}`}>
                        <AgGridReact<LeaderboardEntry>
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
                            paginationPageSize={10}
                            paginationPageSizeSelector={[10, 25, 50]}
                        />
                    </div>
                )}
            </section>
        </main>
    )
}
