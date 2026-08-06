import GlossaryImage from '@/shared/assets/images/glossary.webp'
import {useDocumentTitle} from '@/shared/lib/use-document-title'
import {ScamSchemeGrid} from '@/widgets/scam-scheme-grid'
import styles from './GlossaryPage.module.scss'

export function GlossaryPage() {
    useDocumentTitle('Глоссарий')

    return (
        <main>
            <section className={styles.header}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Глоссарий</h2>
                    <span className={styles.description}>
            Каждая тема — реальный сценарий мошенничества, составленный по мотивам
            актуальных уловок
          </span>
                </div>
                <img src={GlossaryImage} alt="глоссарий" className={styles.illustration} loading="eager"
                     fetchPriority="high"/>
            </section>

            <ScamSchemeGrid/>
        </main>
    )
}