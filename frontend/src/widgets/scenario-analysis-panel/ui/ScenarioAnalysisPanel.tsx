// src/widgets/scenario-analysis-panel/ui/ScenarioAnalysisPanel.tsx

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
            {/*<h2 className={styles.title}>{analysis.title}</h2>*/}

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

            {/*<div className={styles.divider} />*/}

            {/*<section className={styles.safeSection}>*/}
            {/*  <h3 className={styles.safeTitle}>*/}
            {/*    Как действовать безопасно*/}
            {/*  </h3>*/}

            {/*  <ul className={styles.safeActions}>*/}
            {/*    {analysis.safeActions.map((action) => (*/}
            {/*      <li key={action} className={styles.safeAction}>*/}
            {/*        <span className={styles.check} aria-hidden="true">*/}
            {/*          <Icon icon={faCheck} />*/}
            {/*        </span>*/}

            {/*        <span>{action}</span>*/}
            {/*      </li>*/}
            {/*    ))}*/}
            {/*  </ul>*/}
            {/*</section>*/}

            {/*<div className={styles.goldenRule}>*/}
            {/*  <strong className={styles.goldenRuleTitle}>*/}
            {/*    Золотое правило:*/}
            {/*  </strong>{' '}*/}
            {/*  {analysis.goldenRule}*/}
            {/*</div>*/}
        </aside>
    )
}
