import { faChevronLeft, faPause } from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import styles from './GameHeader.module.scss'
import { FragmentCounter } from '@/widgets/header/ui/FragmentCounter.tsx'
import { PointsCounter } from '@/widgets/header/ui/PointsCounter.tsx'
import { Button } from '@/shared/ui/button'
import { useLocation, useNavigate } from 'react-router-dom'
import { faHome } from '@fortawesome/free-regular-svg-icons'
import { SessionTime } from '@/entities/session-timer/ui/SessionTime.tsx'

export type GameHeaderVariant = 'setup' | 'session'

export type GameHeaderProps = {
    variant?: GameHeaderVariant
    timerSeconds?: number
    score?: number
    lives?: number
    maxLives?: number
    currentQuestion?: number
    totalQuestions?: number
}

export function GameHeader({
    variant = 'setup',
    // currentQuestion = 1,
    // totalQuestions = 10,
}: GameHeaderProps) {
    const isSession = variant === 'session'

    const navigate = useNavigate()
    const { pathname } = useLocation()

    function handlePause() {
        const parentPath = pathname.replace(/\/[^/]+$/, '')

        navigate(parentPath)
    }

    const handleBack = () => {
        navigate(-1)
    }

    const handleHome = () => {
        navigate('/')
    }

    return (
        <header className={styles.header}>
            <div className={styles.inner}>
                <div className={styles.left}>
                    <div className={styles.controls}>
                        <Button variant="ghost" onClick={handleBack}>
                            <Icon icon={faChevronLeft} />
                        </Button>
                        {isSession ? (
                            <Button variant="ghost" onClick={handlePause}>
                                <Icon icon={faPause} />
                            </Button>
                        ) : (
                            <Button variant="ghost" onClick={handleHome}>
                                <Icon icon={faHome} />
                            </Button>
                        )}
                    </div>

                    {isSession && <SessionTime />}

                    <span className={styles.separator} aria-hidden="true" />

                    <div className={styles.stats}>
                        <FragmentCounter value={0} />
                        <PointsCounter value={0} />
                    </div>
                </div>

                {/*<div className={styles.right}>*/}
                {/*    {isSession ? (*/}
                {/*        <>*/}
                {/*            <span className={styles.question}>*/}
                {/*                {currentQuestion} из {totalQuestions}*/}
                {/*            </span>*/}

                {/*            <div*/}
                {/*                className={styles.progress}*/}
                {/*                aria-label={`Вопрос ${currentQuestion} из ${totalQuestions}`}*/}
                {/*            >*/}
                {/*                {Array.from({*/}
                {/*                    length: totalQuestions,*/}
                {/*                }).map((_, index) => {*/}
                {/*                    const questionNumber = index + 1*/}

                {/*                    const state =*/}
                {/*                        questionNumber < currentQuestion*/}
                {/*                            ? 'completed'*/}
                {/*                            : questionNumber === currentQuestion*/}
                {/*                              ? 'current'*/}
                {/*                              : 'pending'*/}

                {/*                    return (*/}
                {/*                        <span*/}
                {/*                            key={questionNumber}*/}
                {/*                            className={styles.progressDot}*/}
                {/*                            data-state={state}*/}
                {/*                            aria-hidden="true"*/}
                {/*                        />*/}
                {/*                    )*/}
                {/*                })}*/}
                {/*            </div>*/}
                {/*        </>*/}
                {/*    ) : null}*/}
                {/*</div>*/}
            </div>
        </header>
    )
}
