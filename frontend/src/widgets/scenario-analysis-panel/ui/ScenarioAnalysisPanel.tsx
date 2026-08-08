import type { ScenarioAnalysis, ScenarioRedFlag } from '@/entities/training-scenario'

import styles from './ScenarioAnalysisPanel.module.scss'

type ScenarioAnalysisPanelProps = {
    analysis: ScenarioAnalysis
    revealed?: boolean
    className?: string
}

function HighlightedDescription({ redFlag }: { redFlag: ScenarioRedFlag }) {
    if (!redFlag.accent || !redFlag.description.includes(redFlag.accent)) {
        return <>{redFlag.description}</>
    }

    const [before, after] = redFlag.description.split(redFlag.accent, 2)

    return (
        <>
            {before}

            <strong className={styles.accent}>{redFlag.accent}</strong>

            {after}
        </>
    )
}

export function ScenarioAnalysisPanel({
    analysis,
    revealed = true,
    className,
}: ScenarioAnalysisPanelProps) {
    const rootClassName = [styles.panel, className].filter(Boolean).join(' ')

    if (!revealed) {
        return (
            <aside className={rootClassName}>
                <div className={styles.locked}>
                    <span className={styles.lockedLabel}>Разбор сценария</span>

                    <p className={styles.lockedDescription}>
                        Разбор появится после завершения переписки
                    </p>
                </div>
            </aside>
        )
    }

    return (
        <aside className={rootClassName}>
            <div className={styles.redFlags}>
                {analysis.redFlags.map((redFlag) => (
                    <section key={redFlag.id} className={styles.redFlag}>
                        <h3 className={styles.redFlagTitle}>{redFlag.title}</h3>

                        <p className={styles.redFlagDescription}>
                            <HighlightedDescription redFlag={redFlag} />
                        </p>
                    </section>
                ))}
            </div>
        </aside>
    )
}
