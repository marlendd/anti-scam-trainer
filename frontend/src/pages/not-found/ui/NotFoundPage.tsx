import NotFoundImage from '@/shared/assets/images/404.webp'
import { useDocumentTitle } from '@/shared/lib/use-document-title'

import styles from './NotFoundPage.module.scss'

export function NotFoundPage() {
    useDocumentTitle('Страница не найдена')

    return (
        <main className={styles.page}>
            <div className={styles.content}>
                <img src={NotFoundImage} alt="404" className={styles.image} />

                <h1 className={styles.title}>Страница не найдена</h1>

                <p className={styles.description}>
                    Возможно, она была удалена или адрес указан неверно.
                </p>
            </div>
        </main>
    )
}
