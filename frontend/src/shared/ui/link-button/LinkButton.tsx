import type { ReactNode } from 'react'
import type { To } from 'react-router-dom'
import { Link } from 'react-router-dom'
import type { IconDefinition } from '@fortawesome/fontawesome-svg-core'
import { Icon } from '@/shared/ui/icon'
import styles from './LinkButton.module.scss'

type LinkButtonProps = {
  to: To
  icon: IconDefinition
  label: string
  className?: string
  children?: ReactNode
}

export function LinkButton({
  to,
  icon,
  label,
  className,
  children,
}: LinkButtonProps) {
  const classes = [styles.linkButton, className].filter(Boolean).join(' ')

  return (
    <Link className={classes} to={to} aria-label={label}>
      <Icon icon={icon} />
      <span className={styles.label}>{children ?? label}</span>
    </Link>
  )
}
