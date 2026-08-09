import { useId } from 'react'

import type { PuzzlePieceId } from '../model/types'

import styles from './PuzzleBoard.module.scss'

type PuzzleBoardProps = {
    imageSrc: string
    unlockedPieces: PuzzlePieceId[]
    highlightedPiece?: PuzzlePieceId | null
    onPieceClick?: (pieceId: PuzzlePieceId) => void
}

const BOARD_SIZE = 300
const GRID_SIZE = 3
const CELL_SIZE = BOARD_SIZE / GRID_SIZE
const TAB_SIZE = 18

const pieces: PuzzlePieceId[] = [1, 2, 3, 4, 5, 6, 7, 8, 9]

/**
 * Вертикальные соединения.
 *
 *  1 -> выступ вправо
 * -1 -> выступ влево
 */
const verticalConnections = [
    [1, -1],
    [-1, 1],
    [1, -1],
]

/**
 * Горизонтальные соединения.
 *
 *  1 -> выступ вниз
 * -1 -> выступ вверх
 */
const horizontalConnections = [
    [-1, 1, -1],
    [1, -1, 1],
]

function horizontalEdge(startX: number, endX: number, y: number, connection = 0) {
    if (connection === 0) {
        return `L ${endX} ${y}`
    }

    const direction = Math.sign(endX - startX)
    const length = Math.abs(endX - startX)

    const x1 = startX + direction * length * 0.38
    const x2 = startX + direction * length * 0.42
    const x3 = startX + direction * length * 0.4

    const center = startX + direction * length * 0.5

    const x4 = startX + direction * length * 0.6
    const x5 = startX + direction * length * 0.58
    const x6 = startX + direction * length * 0.62

    const tabY = y + connection * TAB_SIZE

    return [
        `L ${x1} ${y}`,
        `C ${x2} ${y} ${x3} ${tabY} ${center} ${tabY}`,
        `C ${x4} ${tabY} ${x5} ${y} ${x6} ${y}`,
        `L ${endX} ${y}`,
    ].join(' ')
}

function verticalEdge(x: number, startY: number, endY: number, connection = 0) {
    if (connection === 0) {
        return `L ${x} ${endY}`
    }

    const direction = Math.sign(endY - startY)
    const length = Math.abs(endY - startY)

    const y1 = startY + direction * length * 0.38
    const y2 = startY + direction * length * 0.42
    const y3 = startY + direction * length * 0.4

    const center = startY + direction * length * 0.5

    const y4 = startY + direction * length * 0.6
    const y5 = startY + direction * length * 0.58
    const y6 = startY + direction * length * 0.62

    const tabX = x + connection * TAB_SIZE

    return [
        `L ${x} ${y1}`,
        `C ${x} ${y2} ${tabX} ${y3} ${tabX} ${center}`,
        `C ${tabX} ${y4} ${x} ${y5} ${x} ${y6}`,
        `L ${x} ${endY}`,
    ].join(' ')
}

function getPiecePath(row: number, column: number) {
    const x = column * CELL_SIZE
    const y = row * CELL_SIZE

    const right = x + CELL_SIZE
    const bottom = y + CELL_SIZE

    const topConnection = row === 0 ? 0 : horizontalConnections[row - 1][column]

    const rightConnection = column === GRID_SIZE - 1 ? 0 : verticalConnections[row][column]

    const bottomConnection = row === GRID_SIZE - 1 ? 0 : horizontalConnections[row][column]

    const leftConnection = column === 0 ? 0 : verticalConnections[row][column - 1]

    return [
        `M ${x} ${y}`,

        horizontalEdge(x, right, y, topConnection),

        verticalEdge(right, y, bottom, rightConnection),

        horizontalEdge(right, x, bottom, bottomConnection),

        verticalEdge(x, bottom, y, leftConnection),

        'Z',
    ].join(' ')
}

function getPieceCoordinates(pieceId: PuzzlePieceId) {
    const index = pieceId - 1

    return {
        row: Math.floor(index / GRID_SIZE),
        column: index % GRID_SIZE,
    }
}

export function PuzzleBoard({
    imageSrc,
    unlockedPieces,
    highlightedPiece,
    onPieceClick,
}: PuzzleBoardProps) {
    const instanceId = useId().replace(/:/g, '')

    return (
        <div className={styles.board}>
            <svg
                className={styles.svg}
                viewBox={`0 0 ${BOARD_SIZE} ${BOARD_SIZE}`}
                role="img"
                aria-label={`Собрано ${unlockedPieces.length} из ${pieces.length} фрагментов`}
            >
                <defs>
                    {pieces.map((pieceId) => {
                        const { row, column } = getPieceCoordinates(pieceId)

                        return (
                            <clipPath
                                key={pieceId}
                                id={`${instanceId}-piece-${pieceId}`}
                                clipPathUnits="userSpaceOnUse"
                            >
                                <path d={getPiecePath(row, column)} />
                            </clipPath>
                        )
                    })}
                </defs>

                <rect
                    className={styles.background}
                    x="0"
                    y="0"
                    width={BOARD_SIZE}
                    height={BOARD_SIZE}
                    rx="20"
                />

                {pieces.map((pieceId) => {
                    const { row, column } = getPieceCoordinates(pieceId)

                    const path = getPiecePath(row, column)

                    const isUnlocked = unlockedPieces.includes(pieceId)

                    const isHighlighted = highlightedPiece === pieceId

                    const centerX = column * CELL_SIZE + CELL_SIZE / 2

                    const centerY = row * CELL_SIZE + CELL_SIZE / 2

                    return (
                        <g
                            key={pieceId}
                            className={styles.piece}
                            data-unlocked={isUnlocked}
                            data-highlighted={isHighlighted}
                            onClick={() => {
                                if (isUnlocked) {
                                    onPieceClick?.(pieceId)
                                }
                            }}
                        >
                            {isUnlocked ? (
                                <>
                                    <image
                                        href={imageSrc}
                                        x="0"
                                        y="0"
                                        width={BOARD_SIZE}
                                        height={BOARD_SIZE}
                                        preserveAspectRatio="xMidYMid slice"
                                        clipPath={`url(#${instanceId}-piece-${pieceId})`}
                                        className={styles.image}
                                    />

                                    <path d={path} className={styles.unlockedOutline} />
                                </>
                            ) : (
                                <>
                                    <path d={path} className={styles.lockedPiece} />

                                    <text
                                        x={centerX}
                                        y={centerY}
                                        className={styles.question}
                                        textAnchor="middle"
                                        dominantBaseline="central"
                                    >
                                        ?
                                    </text>
                                </>
                            )}
                        </g>
                    )
                })}
            </svg>
        </div>
    )
}
