import {useEffect, useState} from 'react'
import {
    faArrowLeft,
    faArrowRight,
    faPlay,
} from '@fortawesome/free-solid-svg-icons'
import {useNavigate} from 'react-router-dom'

import {Button} from '@/shared/ui/button'
import {Icon} from '@/shared/ui/icon'

import {welcomeSlides} from '../model/welcomeSlides'

import styles from './WelcomePage.module.scss'

import type {CSSProperties} from 'react'

type WelcomePageStyle = CSSProperties & {
    '--slide-accent': string
}

export function WelcomePage() {
    const navigate = useNavigate()

    const [activeSlideIndex, setActiveSlideIndex] = useState(0)

    const activeSlide = welcomeSlides[activeSlideIndex]
    const isFirstSlide = activeSlideIndex === 0
    const isLastSlide =
        activeSlideIndex === welcomeSlides.length - 1

    const pageStyle: WelcomePageStyle = {
        '--slide-accent': activeSlide.accent,
    }

    function handleNext() {
        if (isLastSlide) {
            navigate('/training/role-selection')
            return
        }

        setActiveSlideIndex((currentIndex) => currentIndex + 1)
    }

    function handlePrevious() {
        setActiveSlideIndex((currentIndex) =>
            Math.max(currentIndex - 1, 0),
        )
    }

    function handleSkip() {
        navigate('/training/role-selection')
    }

    useEffect(() => {
        const nextSlide = welcomeSlides[activeSlideIndex + 1]

        if (!nextSlide) {
            return
        }

        const image = new Image()
        image.src = nextSlide.illustration
    }, [activeSlideIndex])

    return (
        <main className={styles.page} style={pageStyle}>
            <button
                type="button"
                className={styles.skipButton}
                onClick={handleSkip}
            >
                Пропустить
            </button>

            <section
                key={activeSlide.id}
                className={styles.slide}
                aria-live="polite"
            >
                <div className={styles.content}>

                    <div className={styles.top}>

                          <span className={styles.eyebrow}>
                            {activeSlide.eyebrow}
                          </span>

                        <div className={styles.text}>
                            <h1 className={styles.title}>
                                {activeSlide.title}
                            </h1>

                            <p className={styles.description}>
                                {activeSlide.description}
                            </p>
                        </div>
                    </div>

                    <div className={styles.controls}>
                        {!isFirstSlide && (
                            <Button
                                variant="secondary"
                                aria-label="Предыдущий слайд"
                                startIcon={<Icon icon={faArrowLeft}/>}
                                onClick={handlePrevious}
                            >
                                Назад
                            </Button>
                        )}

                        <Button
                            endIcon={
                                <Icon
                                    icon={isLastSlide ? faPlay : faArrowRight}
                                />
                            }
                            onClick={handleNext}
                        >
                            {isLastSlide ? 'Начать игру' : 'Далее'}
                        </Button>
                    </div>
                </div>

                <div className={styles.illustrationWrapper}>
                    <div
                        className={styles.illustrationBackground}
                        aria-hidden="true"
                    />

                    <img
                        className={styles.illustration}
                        src={activeSlide.illustration}
                        alt={activeSlide.illustrationAlt}
                        width={720}
                        height={720}
                        loading="eager"
                        decoding="async"
                    />
                </div>
            </section>

            <nav
                className={styles.pagination}
                aria-label="Слайды знакомства"
            >
                {welcomeSlides.map((slide, index) => (
                    <button
                        key={slide.id}
                        type="button"
                        className={styles.paginationButton}
                        data-active={index === activeSlideIndex}
                        aria-label={`Перейти к слайду ${index + 1}`}
                        aria-current={
                            index === activeSlideIndex ? 'step' : undefined
                        }
                        onClick={() => setActiveSlideIndex(index)}
                    />
                ))}
            </nav>
        </main>
    )
}