import { ForgotPasswordForm } from '@/features/auth/ui/ForgotPasswordForm'

import styles from './ForgotPasswordPage.module.scss'

export function ForgotPasswordPage() {
  return (
    <main className={styles.page}>
      <section className={styles.card}>
        <header className={styles.header}>
          {/*<span className={styles.eyebrow}>*/}
          {/*  Восстановление доступа*/}
          {/*</span>*/}

          {/*<h1 className={styles.title}>*/}
          {/*  Забыли пароль?*/}
          {/*</h1>*/}

          {/*<p className={styles.description}>*/}
          {/*  Введите электронную почту, которую использовали*/}
          {/*  при регистрации. Мы отправим ссылку для создания*/}
          {/*  нового пароля.*/}
          {/*</p>*/}
        </header>

        <ForgotPasswordForm />
      </section>
    </main>
  )
}