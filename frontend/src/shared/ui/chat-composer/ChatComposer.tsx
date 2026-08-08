// src/shared/ui/chat-composer/ChatComposer.tsx

import {
    useLayoutEffect,
    useRef,
    type ChangeEvent,
    type FormEvent,
    type KeyboardEvent,
} from 'react'
import { faArrowUp } from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import styles from './ChatComposer.module.scss'

type ChatComposerProps = {
    value: string
    onChange: (value: string) => void
    onSend: () => void
    placeholder?: string
    disabled?: boolean
    maxLength?: number
    className?: string
}

export function ChatComposer({
    value,
    onChange,
    onSend,
    placeholder = 'Сообщение',
    disabled = false,
    maxLength = 2000,
    className,
}: ChatComposerProps) {
    const textareaRef = useRef<HTMLTextAreaElement>(null)

    useLayoutEffect(() => {
        const textarea = textareaRef.current

        if (!textarea) {
            return
        }

        textarea.style.height = '0px'
        textarea.style.height = `${Math.min(textarea.scrollHeight, 160)}px`
    }, [value])

    function sendMessage() {
        if (!value.trim() || disabled) {
            return
        }

        onSend()
    }

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        sendMessage()
    }

    function handleChange(event: ChangeEvent<HTMLTextAreaElement>) {
        onChange(event.target.value)
    }

    function handleKeyDown(event: KeyboardEvent<HTMLTextAreaElement>) {
        if (event.key !== 'Enter' || event.shiftKey || event.nativeEvent.isComposing) {
            return
        }

        event.preventDefault()
        sendMessage()
    }

    const rootClassName = [styles.composer, className].filter(Boolean).join(' ')

    return (
        <form className={rootClassName} onSubmit={handleSubmit}>
            <textarea
                ref={textareaRef}
                className={styles.input}
                value={value}
                placeholder={placeholder}
                disabled={disabled}
                maxLength={maxLength}
                rows={1}
                aria-label="Сообщение"
                onChange={handleChange}
                onKeyDown={handleKeyDown}
            />

            <button
                className={styles.sendButton}
                type="submit"
                disabled={disabled || !value.trim()}
                aria-label="Отправить сообщение"
            >
                <Icon icon={faArrowUp} />
            </button>
        </form>
    )
}
