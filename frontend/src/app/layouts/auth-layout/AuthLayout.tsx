// src/app/layouts/auth-layout/AuthLayout.tsx

import { Outlet } from 'react-router-dom'
import styles from './AuthLayout.module.scss'


export function AuthLayout() {
  return (
    <main className={styles.page}>
      <section className={styles.card}>
        <Outlet />
      </section>
    </main>
  )
}