import { faPuzzlePiece } from '@fortawesome/free-solid-svg-icons'
import { Icon } from '@/shared/ui/icon'
import styles from './counter.module.scss'

type FragmentCounterProps = {
  value: number
}

export function FragmentCounter({ value }: FragmentCounterProps) {
  return (
    <div className={styles.counter}>
      <Icon icon={faPuzzlePiece} style={{ color: '#965eeb' }}/>
      <span>{value}</span>
    </div>
  )
}
