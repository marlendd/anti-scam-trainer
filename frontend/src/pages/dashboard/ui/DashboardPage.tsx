import {DashboardRatingChart} from "@/widgets/dashboard-rating-chart";
import {DashboardOverview} from "@/widgets/dashboard-overview";

import styles from './DashboardPage.module.scss';
import {DashboardHeatMap} from "@/widgets/dashboard-heat-map";
import {DashboardBalanceChart} from "@/widgets/dashboard-balance-chart";
import {DashboardScenariosBarChart} from "@/widgets/dashboard-scenarios-bar-chart";

export function DashboardPage() {
  return (
    <main className={styles.page}>
        <DashboardOverview/>
        <DashboardRatingChart/>
        <DashboardHeatMap/>
        <DashboardBalanceChart/>
        <DashboardScenariosBarChart/>
    </main>
  )
}
