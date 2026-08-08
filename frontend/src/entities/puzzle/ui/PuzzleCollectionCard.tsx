import { faCheck, faGift } from '@fortawesome/free-solid-svg-icons'

import type { PuzzleCollection } from '../model/types'

import { Button } from '@/shared/ui/button'
import { Icon } from '@/shared/ui/icon'

import { PuzzleBoard } from './PuzzleBoard'

import styles from './PuzzleCollectionCard.module.scss'

type PuzzleCollectionCardProps = {
    collection: PuzzleCollection
    onClaimReward?: (collectionId: string) => void
}

const TOTAL_PIECES = 9

export function PuzzleCollectionCard({ collection, onClaimReward }: PuzzleCollectionCardProps) {
    const unlockedCount = collection.unlockedPieces.length

    const isCompleted = unlockedCount === TOTAL_PIECES

    const piecesLeft = TOTAL_PIECES - unlockedCount

    return (
        <article className={styles.card} data-completed={isCompleted}>
            <div className={styles.puzzle}>
                <PuzzleBoard
                    imageSrc={collection.imageSrc}
                    unlockedPieces={collection.unlockedPieces}
                />

                {isCompleted && (
                    <div className={styles.completedBadge}>
                        <Icon icon={faCheck} />
                        Собрано
                    </div>
                )}
            </div>

            <div className={styles.content}>
                <div>
                    <h2 className={styles.title}>{collection.title}</h2>

                    <p className={styles.description}>{collection.description}</p>
                </div>

                <div className={styles.reward}>
                    <div className={styles.rewardIcon}>
                        <Icon icon={faGift} />
                    </div>

                    <div className={styles.rewardContent}>
                        <strong className={styles.rewardTitle}>{collection.reward.title}</strong>

                        {collection.reward.description && (
                            <span className={styles.rewardDescription}>
                                {collection.reward.description}
                            </span>
                        )}
                    </div>
                </div>

                <div className={styles.footer}>
                    {collection.rewardClaimed ? (
                        <div className={styles.claimed}>
                            <Icon icon={faCheck} />
                            Награда получена
                        </div>
                    ) : isCompleted ? (
                        <Button
                            fullWidth
                            onClick={() => {
                                onClaimReward?.(collection.id)
                            }}
                        >
                            Получить награду
                        </Button>
                    ) : (
                        <span className={styles.remaining}>
                            {piecesLeft === 1
                                ? 'Остался 1 фрагмент'
                                : `Осталось ${piecesLeft} фрагмента`}
                        </span>
                    )}
                </div>
            </div>
        </article>
    )
}
