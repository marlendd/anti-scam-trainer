import { useMemo, useState } from 'react'

import {
  PuzzleCollectionCard,
  puzzleCollectionsMock,
} from '@/entities/puzzle'

import { useDocumentTitle } from '@/shared/lib/use-document-title'

import styles from './PuzzlePage.module.scss'

export function PuzzlePage() {
  useDocumentTitle('Коллекция')

  const [claimedRewards, setClaimedRewards] =
    useState<string[]>([])

  const collections = useMemo(
    () =>
      puzzleCollectionsMock.map(
        (collection) => ({
          ...collection,
          rewardClaimed:
            claimedRewards.includes(
              collection.id,
            ),
        }),
      ),
    [claimedRewards],
  )

  const totalPieces =
    collections.length * 9

  const collectedPieces = collections.reduce(
    (total, collection) =>
      total +
      collection.unlockedPieces.length,
    0,
  )

  function handleClaimReward(
    collectionId: string,
  ) {
    setClaimedRewards((current) => {
      if (current.includes(collectionId)) {
        return current
      }

      return [
        ...current,
        collectionId,
      ]
    })
  }

  return (
    <main className={styles.page}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>
            Коллекция
          </h1>

          <p className={styles.description}>
            Проходите сценарии, собирайте
            пазлы и открывайте бонусы.
          </p>
        </div>

        <div className={styles.summary}>
          <span className={styles.summaryValue}>
            {collectedPieces}
          </span>

          <span className={styles.summaryLabel}>
            из {totalPieces} фрагментов
          </span>
        </div>
      </header>

      <div className={styles.grid}>
        {collections.map((collection) => (
          <PuzzleCollectionCard
            key={collection.id}
            collection={collection}
            onClaimReward={
              handleClaimReward
            }
          />
        ))}
      </div>
    </main>
  )
}