import type { ComponentProps } from 'react'
import type { IconDefinition } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import styles from './Icon.module.scss'

type FontAwesomeIconProps = ComponentProps<typeof FontAwesomeIcon>

type IconProps = {
    icon: IconDefinition
    label?: string
    decorative?: boolean
    className?: string
    style?: FontAwesomeIconProps['style']
}

export function Icon({ icon, label, decorative = true, className, style }: IconProps) {
    const classes = [styles.icon, className].filter(Boolean).join(' ')

    return (
        <FontAwesomeIcon
            className={classes}
            icon={icon}
            aria-hidden={decorative ? true : undefined}
            aria-label={label}
            role={label ? 'img' : undefined}
            focusable="false"
            style={style}
        />
    )
}
