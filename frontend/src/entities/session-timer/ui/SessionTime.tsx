import { useEffect, useState } from 'react'

import { faStopwatch } from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import styles from './SessionTime.module.scss'

type SessionTimeProps = {
    isRunning?: boolean
    initialSeconds?: number
    storageKey?: string
}

function formatTime(totalSeconds: number) {
    const minutes = Math.floor(totalSeconds / 60)
    const seconds = totalSeconds % 60

    return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
}

export function SessionTime({
    isRunning = true,
    initialSeconds = 0,
    storageKey = 'session-time',
}: SessionTimeProps) {
    const [elapsedSeconds, setElapsedSeconds] = useState(() => {
        const savedValue = localStorage.getItem(storageKey)

        if (savedValue === null) {
            return initialSeconds
        }

        const parsedValue = Number(savedValue)

        return Number.isFinite(parsedValue) ? parsedValue : initialSeconds
    })

    useEffect(() => {
        if (!isRunning) {
            return
        }

        const interval = window.setInterval(() => {
            setElapsedSeconds((current) => current + 1)
        }, 1000)

        return () => {
            window.clearInterval(interval)
        }
    }, [isRunning])

    useEffect(() => {
        localStorage.setItem(storageKey, String(elapsedSeconds))
    }, [elapsedSeconds, storageKey])

    return (
        <div className={styles.time} aria-label={`Время сессии: ${formatTime(elapsedSeconds)}`}>
            <Icon icon={faStopwatch} />

            <span className={styles.value}>{formatTime(elapsedSeconds)}</span>
        </div>
    )
}
