import { Outlet } from 'react-router-dom'

import { Header } from '@/widgets/header'

import styles from './DefaultLayout.module.scss'

export const DefaultLayout = () => {
    return (
        <div className={styles.layout}>
            <Header />

            <div className={styles.body}>
                <main className={styles.content}>
                    <Outlet />
                </main>
            </div>
        </div>
    )
}
