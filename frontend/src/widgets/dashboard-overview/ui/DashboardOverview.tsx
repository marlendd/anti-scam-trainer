import styles from './DashboardOverview.module.scss'
import { Avatar } from '@/shared/ui/avatar'
import { LogoutButton } from '@/features/logout-button'
import {useCurrentUser} from "@/entities/user";

export function DashboardOverview() {
    const {
        data: user,
        isLoading,
        isError,
    } = useCurrentUser()

    if (isLoading) {
        return null
    }

    if (isError || !user) {
        return null
    }

    return (
        <section className={styles.overview}>
            <div className={styles.user}>
                <Avatar size={120} />

                <div className={styles.text}>
                    <p className={styles.name}>
                        {user.name || user.email}
                    </p>

                    <p className={styles.rank}>
                        <LogoutButton />
                    </p>
                </div>
            </div>
        </section>
    )
}