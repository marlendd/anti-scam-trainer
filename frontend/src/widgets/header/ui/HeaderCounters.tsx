import { useProfileSummary } from '@/entities/profile-progress'
import { FragmentCounter } from './FragmentCounter.tsx'
import { PointsCounter } from './PointsCounter.tsx'

export function HeaderCounters() {
    const {
        data,
        isPending,
    } = useProfileSummary()

    return (
        <>
            <PointsCounter
                value={isPending ? 0 : data?.totalScore ?? 0}
            />

            <FragmentCounter
                value={isPending ? 0 : data?.totalFragments ?? 0}
            />
        </>
    )
}