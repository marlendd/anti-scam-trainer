import type { ImgHTMLAttributes } from 'react'
import logo from '@/shared/assets/images/LogoScam.svg'
import logoInverse from '@/shared/assets/images/Logo_inverse.svg'
import styles from './Logo.module.scss'

type LogoVariant = 'default' | 'inverse'

type LogoProps = Omit<ImgHTMLAttributes<HTMLImageElement>, 'src' | 'alt'> & {
  variant?: LogoVariant
  alt?: string
}

export function Logo({
  variant = 'default',
  alt = 'Logo',
  className,
  ...props
}: LogoProps) {
  const image = variant === 'inverse' ? logoInverse : logo
  const classes = [styles.logo, className].filter(Boolean).join(' ')

  return (
      <div className={styles.logo}>
        <img className={classes} src={image} alt={alt} {...props} />
        <span className={styles.info}>

        </span>
      </div>
  )
}
