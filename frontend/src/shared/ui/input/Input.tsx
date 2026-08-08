// src/shared/ui/input/Input.tsx

import { useId, type ComponentPropsWithRef, type ReactNode } from 'react'

import styles from './Input.module.scss'

export type InputProps = ComponentPropsWithRef<'input'> & {
    label: string
    error?: string
    hint?: string
    endAdornment?: ReactNode
    containerClassName?: string
}

export function Input({
    ref,
    id,
    label,
    error,
    hint,
    endAdornment,
    className,
    containerClassName,
    required,
    'aria-describedby': ariaDescribedBy,
    'aria-invalid': ariaInvalid,
    ...props
}: InputProps) {
    const generatedId = useId()
    const inputId = id ?? generatedId

    const hintId = `${inputId}-hint`
    const errorId = `${inputId}-error`

    const describedBy =
        [ariaDescribedBy, hint ? hintId : undefined, error ? errorId : undefined]
            .filter(Boolean)
            .join(' ') || undefined

    const inputClassName = [
        styles.input,
        endAdornment ? styles.inputWithAdornment : undefined,
        error ? styles.inputError : undefined,
        className,
    ]
        .filter(Boolean)
        .join(' ')

    const fieldClassName = [styles.field, containerClassName].filter(Boolean).join(' ')

    return (
        <div className={fieldClassName}>
            <label className={styles.label} htmlFor={inputId}>
                {label}

                {required && (
                    <span className={styles.required} aria-hidden="true">
                        *
                    </span>
                )}
            </label>

            <div className={styles.control}>
                <input
                    {...props}
                    ref={ref}
                    id={inputId}
                    required={required}
                    className={inputClassName}
                    aria-describedby={describedBy}
                    aria-invalid={error ? true : ariaInvalid}
                />

                {endAdornment && <div className={styles.endAdornment}>{endAdornment}</div>}
            </div>

            {error ? (
                <span id={errorId} className={styles.error} role="alert">
                    {error}
                </span>
            ) : hint ? (
                <span id={hintId} className={styles.hint}>
                    {hint}
                </span>
            ) : null}
        </div>
    )
}
