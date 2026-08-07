import styles from './DashboardOverview.module.scss';
import {Avatar} from "@/shared/ui/avatar";
import {Icon} from "@/shared/ui/icon";
import {faArrowRightFromBracket} from "@fortawesome/free-solid-svg-icons";
import {Link} from "react-router-dom";

export function DashboardOverview() {
    return <section className={styles.overview}>

        <div className={styles.user}>
            <Avatar size={120}/>

            <div className={styles.text}>
                <p className={styles.name}>
                    Unknown User
                </p>
                <p className={styles.rank}>
                    <Link to="/">
                        <span>Выйти</span>
                        <Icon icon={faArrowRightFromBracket} style={{ color: 'inherit', height: '14px' }} />
                    </Link>
                </p>
            </div>
        </div>
    </section>
}