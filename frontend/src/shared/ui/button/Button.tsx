import type { ComponentPropsWithRef, ReactNode } from 'react'

import styles from './Button.module.scss'

export type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'danger'

export type ButtonSize = 'small' | 'medium' | 'large'

export type ButtonProps = ComponentPropsWithRef<'button'> & {
    variant?: ButtonVariant
    size?: ButtonSize
    fullWidth?: boolean
    isLoading?: boolean
    startIcon?: ReactNode
    endIcon?: ReactNode
}

export function Button({
    ref,
    children,
    className,
    variant = 'primary',
    size = 'medium',
    fullWidth = false,
    isLoading = false,
    startIcon,
    endIcon,
    disabled,
    type = 'button',
    ...props
}: ButtonProps) {
    const rootClassName = [styles.button, fullWidth ? styles.fullWidth : undefined, className]
        .filter(Boolean)
        .join(' ')

    return (
        <button
            {...props}
            ref={ref}
            type={type}
            className={rootClassName}
            data-variant={variant}
            data-size={size}
            disabled={disabled || isLoading}
            aria-busy={isLoading || undefined}
        >
            {isLoading ? (
                <span className={styles.spinner} aria-hidden="true" />
            ) : (
                startIcon && (
                    <span className={styles.icon} aria-hidden="true">
                        {startIcon}
                    </span>
                )
            )}

            <span className={styles.content}>{children}</span>

            {!isLoading && endIcon && (
                <span className={styles.icon} aria-hidden="true">
                    {endIcon}
                </span>
            )}
        </button>
    )
}
