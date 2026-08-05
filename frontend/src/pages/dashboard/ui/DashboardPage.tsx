import { useDocumentTitle } from '@/shared/lib/use-document-title'
import { DashboardBalanceChart } from '@/widgets/dashboard-balance-chart'
import { DashboardHeatMap } from '@/widgets/dashboard-heat-map'
import { DashboardOverview } from '@/widgets/dashboard-overview'
import { DashboardRatingChart } from '@/widgets/dashboard-rating-chart'
import { DashboardScenariosBarChart } from '@/widgets/dashboard-scenarios-bar-chart'
import styles from './DashboardPage.module.scss'

export function DashboardPage() {
  useDocumentTitle('Личный кабинет')

  return (
    <main className={styles.page}>
      <DashboardOverview />
      <DashboardRatingChart />
      <DashboardScenariosBarChart />
      <DashboardBalanceChart />
      <DashboardHeatMap />
    </main>
  )
}