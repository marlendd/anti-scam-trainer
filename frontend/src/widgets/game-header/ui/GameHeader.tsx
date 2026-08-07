import {
    faChevronLeft, faPause,
    // faHeart,
    faStopwatch,
} from '@fortawesome/free-solid-svg-icons'

import {Icon} from '@/shared/ui/icon'

import styles from './GameHeader.module.scss'
import {FragmentCounter} from "@/widgets/header/ui/FragmentCounter.tsx";
import {PointsCounter} from "@/widgets/header/ui/PointsCounter.tsx";
import {Button} from "@/shared/ui/button";
import {useLocation, useNavigate} from "react-router-dom";
import {faHome} from "@fortawesome/free-regular-svg-icons";

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

function formatTimer(totalSeconds: number) {
    const safeSeconds = Math.max(0, totalSeconds)

    const minutes = Math.floor(safeSeconds / 60)
    const seconds = safeSeconds % 60

    return `${String(minutes).padStart(2, '0')}:${String(
        seconds,
    ).padStart(2, '0')}`
}

export function GameHeader({
                               variant = 'setup',
                               timerSeconds = 0,
                               // lives = 3,
                               // maxLives = 3,
                               currentQuestion = 1,
                               totalQuestions = 10,
                           }: GameHeaderProps) {
    const isSession = variant === 'session'

    const navigate = useNavigate()
    const {pathname} = useLocation()

    function handlePause() {
        const parentPath = pathname.replace(/\/[^/]+$/, '')

        navigate(parentPath)
    }

    const handleBack = () => {
        navigate(-1);
    }

    const handleHome = () => {
        navigate('/');
    }

    return (
        <header className={styles.header}>
            <div className={styles.inner}>
                <div className={styles.left}>

                    <div className={styles.controls}>
                        <Button variant='ghost' onClick={handleBack}>
                            <Icon icon={faChevronLeft}/>
                        </Button>
                        {
                            isSession ?
                                <Button variant='ghost' onClick={handlePause}>
                                    <Icon icon={faPause}/>
                                </Button>
                                :
                                <Button variant='ghost' onClick={handleHome}>
                                    <Icon icon={faHome}/>
                                </Button>
                        }
                    </div>

                    {isSession ? (
                        <>
                            <div
                                className={styles.timer}
                                aria-label={`Осталось времени: ${formatTimer(
                                    timerSeconds,
                                )}`}
                            >
                <span
                    className={styles.timerIcon}
                    aria-hidden="true"
                >
                  <Icon icon={faStopwatch} style={{color: '#111'}}/>
                </span>

                                <span className={styles.timerValue}>
                  {formatTimer(timerSeconds)}
                </span>
                            </div>


                        </>
                    ) : null}
                    <span
                        className={styles.separator}
                        aria-hidden="true"
                    />

                    <div className={styles.stats}>
                        <FragmentCounter value={0}/>
                        <PointsCounter value={0}/>
                    </div>
                </div>

                {/*<div*/}
                {/*  className={styles.lives}*/}
                {/*  aria-label={`Жизни: ${lives} из ${maxLives}`}*/}
                {/*>*/}
                {/*  <span className={styles.livesLabel}>*/}
                {/*    Жизни:*/}
                {/*  </span>*/}

                {/*  <span className={styles.hearts}>*/}
                {/*    {Array.from({ length: maxLives }).map(*/}
                {/*      (_, index) => (*/}
                {/*        <span*/}
                {/*          key={index}*/}
                {/*          className={styles.heart}*/}
                {/*          data-active={index < lives}*/}
                {/*          aria-hidden="true"*/}
                {/*        >*/}
                {/*          <Icon icon={faHeart} />*/}
                {/*        </span>*/}
                {/*      ),*/}
                {/*    )}*/}
                {/*  </span>*/}
                {/*</div>*/}

                <div className={styles.right}>
                    {isSession ? (
                        <>
              <span className={styles.question}>
                 {currentQuestion} из {totalQuestions}
              </span>

                            <div
                                className={styles.progress}
                                aria-label={`Вопрос ${currentQuestion} из ${totalQuestions}`}
                            >
                                {Array.from({
                                    length: totalQuestions,
                                }).map((_, index) => {
                                    const questionNumber = index + 1

                                    const state =
                                        questionNumber < currentQuestion
                                            ? 'completed'
                                            : questionNumber === currentQuestion
                                                ? 'current'
                                                : 'pending'

                                    return (
                                        <span
                                            key={questionNumber}
                                            className={styles.progressDot}
                                            data-state={state}
                                            aria-hidden="true"
                                        />
                                    )
                                })}
                            </div>
                        </>
                    ) : null}
                </div>
            </div>
        </header>
    )
}