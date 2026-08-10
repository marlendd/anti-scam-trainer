import { Outlet, useMatch } from 'react-router-dom'

import { GameHeader } from '@/widgets/game-header'

import styles from './GameLayout.module.scss'

export function GameLayout() {
    const sessionMatch = useMatch({
        path: '/training/path/:pathId/:schemeId',
        end: true,
    })

    const headerVariant = sessionMatch ? 'session' : 'setup'

    return (
        <div className={styles.layout}>
            <GameHeader variant={headerVariant} />

            <main className={styles.content}>
                <Outlet />
            </main>
        </div>
    )
}
