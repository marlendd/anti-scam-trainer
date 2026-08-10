import { apiRequest } from '@/shared/api'

import type { ProfileSummary } from '../model/types'
import type { ProfileSummaryResponseDto } from './types'

export async function getProfileSummary(): Promise<ProfileSummary> {
    const response = await apiRequest<ProfileSummaryResponseDto>(
        '/profile/summary',
    )

    return {
        totalScore: response.total_score,
        totalFragments: response.total_fragments,
    }
}