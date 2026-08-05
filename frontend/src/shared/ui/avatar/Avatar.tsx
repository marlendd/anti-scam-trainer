import type { HTMLAttributes } from 'react'
import styles from './Avatar.module.scss'
import {Icon} from "@/shared/ui/icon";
import {faUser} from "@fortawesome/free-solid-svg-icons";

type AvatarProps = Omit<HTMLAttributes<HTMLDivElement>, 'children'> & {
  src?: string
  alt?: string
  name?: string
  size?: number
}

function getInitials(value: string) {
  return value
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join('')
}

export function Avatar({
  src,
  alt,
  name,
  size = 40,
  className,
  style,
  ...props
}: AvatarProps) {
  const classes = [styles.avatar, className].filter(Boolean).join(' ')
  const initials = name ? getInitials(name) : <Icon icon={faUser} style={{ color: '#9c9c9c', width: size/2, height: size/2 }}/>

  return (
    <div
      className={classes}
      style={{ width: size, height: size, ...style }}
      {...props}
    >
      {src ? (
        <img className={styles.image} src={src} alt={alt ?? name ?? 'Avatar'} />
      ) : (
        initials
      )}
    </div>
  )
}
