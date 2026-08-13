import { faChevronLeft, faPause } from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import styles from './GameHeader.module.scss'
import { Button } from '@/shared/ui/button'
import { useLocation, useNavigate } from 'react-router-dom'
import { faHome } from '@fortawesome/free-regular-svg-icons'
import { SessionTime } from '@/entities/session-timer/ui/SessionTime.tsx'
import {HeaderCounters} from "@/widgets/header";

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
}: GameHeaderProps) {
    const isSession = variant === 'session'

    const navigate = useNavigate()
    const { pathname } = useLocation()

    function handlePause() {
        const parentPath = pathname.replace(/\/[^/]+$/, '')

        navigate(parentPath)
    }

    const handleBack = () => {
        navigate('/training/role-selection')
    }

    const handleHome = () => {
        navigate('/')
    }

    return (
        <header className={styles.header}>
            <div className={styles.inner}>
                <div className={styles.left}>
                    <div className={styles.controls}>
                        {isSession ? (
                            <>
                                <Button variant="ghost" onClick={handlePause}>
                                    <Icon icon={faChevronLeft} />
                                </Button>
                                <Button variant="ghost" onClick={handlePause}>
                                    <Icon icon={faPause} />
                                </Button>
                            </>
                        ) : (
                            <>
                                <Button variant="ghost" onClick={handleBack}>
                                    <Icon icon={faChevronLeft} />
                                </Button>
                                <Button variant="ghost" onClick={handleHome}>
                                    <Icon icon={faHome} />
                                </Button>
                            </>
                        )}
                    </div>

                    {isSession && <SessionTime />}

                    <span className={styles.separator} aria-hidden="true" />

                    <div className={styles.stats}>
                        <HeaderCounters/>
                    </div>
                </div>
            </div>
        </header>
    )
}
