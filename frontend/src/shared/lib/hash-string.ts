export function hashString(value: string) {
    let hash = 2166136261

    for (let index = 0; index < value.length; index += 1) {
        hash ^= value.charCodeAt(index)
        hash = Math.imul(hash, 16777619)
    }

    return hash >>> 0
}

export function createSeededRandom(seed: number) {
    let state = seed

    return () => {
        state += 0x6d2b79f5

        let value = state

        value = Math.imul(
            value ^ (value >>> 15),
            value | 1,
        )

        value ^=
            value +
            Math.imul(
                value ^ (value >>> 7),
                value | 61,
            )

        return (
            ((value ^ (value >>> 14)) >>> 0) /
            4294967296
        )
    }
}

export function shuffleDeterministic<T>(
    items: T[],
    seed: string,
) {
    const result = [...items]

    const random = createSeededRandom(
        hashString(seed),
    )

    for (
        let index = result.length - 1;
        index > 0;
        index -= 1
    ) {
        const targetIndex = Math.floor(
            random() * (index + 1),
        )

        ;[
            result[index],
            result[targetIndex],
        ] = [
            result[targetIndex],
            result[index],
        ]
    }

    return result
}