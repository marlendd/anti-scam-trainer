import styles from './DashboardOverview.module.scss';
import {Avatar} from "@/shared/ui/avatar";
import {LogoutButton} from "@/features/logout-button";

export function DashboardOverview() {
    return <section className={styles.overview}>

        <div className={styles.user}>
            <Avatar size={120}/>

            <div className={styles.text}>
                <p className={styles.name}>
                    Unknown User
                </p>
                <p className={styles.rank}>
                    <LogoutButton/>
                </p>
            </div>
        </div>
    </section>
}