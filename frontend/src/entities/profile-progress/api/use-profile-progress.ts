import { useQuery } from '@tanstack/react-query'

import {
    getCategoriesProgress,
    getPuzzleProgress,
    getRankHistory,
    getRoleProgress,
} from './get-profile-progress'

export function useRoleProgress() {
    return useQuery({
        queryKey: ['profile-progress', 'roles'],
        queryFn: getRoleProgress,
    })
}

export function useCategoriesProgress() {
    return useQuery({
        queryKey: ['profile-progress', 'categories'],
        queryFn: getCategoriesProgress,
    })
}

export function usePuzzleProgress() {
    return useQuery({
        queryKey: ['profile-progress', 'puzzle'],
        queryFn: getPuzzleProgress,
    })
}

export function useRankHistory() {
    return useQuery({
        queryKey: ['profile-progress', 'rank-history'],
        queryFn: getRankHistory,
    })
}
