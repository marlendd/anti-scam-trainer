import {DashboardRatingChart} from "@/widgets/dashboard-rating-chart";
import {DashboardOverview} from "@/widgets/dashboard-overview";

import styles from './DashboardPage.module.scss';
// import {DashboardSchemesRadar} from "@/widgets/dashboard-schemes-radar";
import {DashboardHeatMap} from "@/widgets/dashboard-heat-map";

export function DashboardPage() {
  return (
    <main className={styles.page}>
        <DashboardOverview/>
        <DashboardRatingChart/>
        {/*<DashboardSchemesRadar/>*/}
        <DashboardHeatMap/>
    </main>
  )
}
