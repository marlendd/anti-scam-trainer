import type { ImgHTMLAttributes } from 'react'

import logoLarge from '@/shared/assets/images/LogoScam.svg'
import logoDefault from '@/shared/assets/images/Logo.svg'
import logoSmall from '@/shared/assets/images/LogoSmall.svg'

import styles from './Logo.module.scss'

type LogoProps = Omit<
  ImgHTMLAttributes<HTMLImageElement>,
  'src'
> & {
  alt?: string
}

export function Logo({
  alt = 'Антискам-тренажёр',
  className,
  ...props
}: LogoProps) {
  const classes = [styles.image, className]
    .filter(Boolean)
    .join(' ')

  return (
    <picture className={styles.logo}>
      <source
        media="(max-width: 860px)"
        srcSet={logoSmall}
      />

      <source
        media="(max-width: 1100px)"
        srcSet={logoDefault}
      />

      <img
        {...props}
        className={classes}
        src={logoLarge}
        alt={alt}
      />
    </picture>
  )
}