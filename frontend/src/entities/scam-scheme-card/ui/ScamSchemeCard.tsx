import styles from './ScamSchemeCard.module.scss'
import { Link } from 'react-router-dom'
import { Icon } from '@/shared/ui/icon'
import { faChevronRight } from '@fortawesome/free-solid-svg-icons'

type ScamSchemeCardProps = {
    id: string
    title: string
    description: string
}

export function ScamSchemeCard({ id, title, description }: ScamSchemeCardProps) {
    return (
        <Link to={`/glossary/${id}`} className={styles.card}>
            <div className={styles.header}>
                <h3 className={styles.title}>{title}</h3>

                <span className={styles.chevron} aria-hidden="true">
                    <Icon icon={faChevronRight} />
                </span>
            </div>
            <p className={styles.description}>{description}</p>
        </Link>
    )
}
