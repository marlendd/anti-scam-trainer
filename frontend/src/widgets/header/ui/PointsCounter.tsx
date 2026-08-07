import { faGem } from '@fortawesome/free-solid-svg-icons'
import { Icon } from '@/shared/ui/icon'
import styles from './counter.module.scss'

type PointsCounterProps = {
  value: number
}

export function PointsCounter({ value }: PointsCounterProps) {
  return (
    <div className={styles.counter}>
      <Icon icon={faGem} style={{ color: '#00aaff' }}/>
      <span>{value}</span>
    </div>
  )
}
