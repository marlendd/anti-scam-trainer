import styles from './ScamSchemeGrid.module.scss'
import { ScamSchemeCard } from '@/entities/scam-scheme-card'
import { scamSchemes } from '../model/mock.ts'

export function ScamSchemeGrid() {
    return (
        <section className={styles.grid}>
            {scamSchemes.map((scheme) => (
                <ScamSchemeCard
                    key={scheme.id}
                    id={scheme.id}
                    title={scheme.title}
                    description={scheme.description}
                />
            ))}
        </section>
    )
}
