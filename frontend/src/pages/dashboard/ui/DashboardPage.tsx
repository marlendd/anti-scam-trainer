import { useDocumentTitle } from '@/shared/lib/use-document-title'
import { DashboardBalanceChart } from '@/widgets/dashboard-balance-chart'
import { DashboardOverview } from '@/widgets/dashboard-overview'
import { DashboardRatingChart } from '@/widgets/dashboard-rating-chart'
import { DashboardScenariosBarChart } from '@/widgets/dashboard-scenarios-bar-chart'
import { DashboardPathProgress } from '@/widgets/dashboard-path-progress'
import styles from './DashboardPage.module.scss'

export function DashboardPage() {
    useDocumentTitle('Личный кабинет')

    return (
        <main className={styles.page}>
            <DashboardOverview />
            <DashboardRatingChart />
            <DashboardBalanceChart />
            <DashboardScenariosBarChart />
            <DashboardPathProgress />
        </main>
    )
}
