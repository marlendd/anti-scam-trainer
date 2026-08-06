import { Outlet } from 'react-router-dom'

import {
  GameHeader,
  type GameHeaderProps,
} from '@/widgets/game-header'

import styles from './GameLayout.module.scss'

type GameLayoutProps = {
  headerProps?: GameHeaderProps
}

export function GameLayout({
  headerProps,
}: GameLayoutProps) {
  return (
    <div className={styles.layout}>
      <GameHeader {...headerProps} />

      <main className={styles.content}>
        <Outlet />
      </main>
    </div>
  )
}