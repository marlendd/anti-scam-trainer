import { faPlay } from '@fortawesome/free-solid-svg-icons'

import { Icon } from '@/shared/ui/icon'

import styles from './PlayButton.module.scss'

export function PlayButton() {
  const rootClassName = [styles.button]
    .filter(Boolean)
    .join(' ')

  return (
    <button
      className={rootClassName}
      aria-label='Запустить'
    >
      <span className={styles.icon} aria-hidden="true">
        <Icon icon={faPlay} style={{ color: 'white', height: '16px' }} />
      </span>
    </button>
  )
}