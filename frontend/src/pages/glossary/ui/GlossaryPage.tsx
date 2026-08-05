import styles from './GlossaryPage.module.scss';
import GlossaryImage from '@/shared/assets/images/glossary.png';
import {ScamSchemeGrid} from "@/widgets/scam-scheme-grid";

export function GlossaryPage() {
    return (
        <main>
            <section className={styles.header}>
                <div className={styles.text}>
                    <h2 className={styles.title}>Глоссарий</h2>
                    <span className={styles.description}>Каждая тема — реальный сценарий мошенничества, составленный по мотивам актуальных уловок</span>
                </div>
                <img src={GlossaryImage} alt="глоссарий" className={styles.illustration}/>
            </section>

            <ScamSchemeGrid/>
        </main>
    )
}
