import { getProfileSummary } from './get-profile-summary'
import {useQuery} from "@tanstack/react-query";

export function useProfileSummary() {
    return useQuery({
        queryKey: ['profile-progress', 'summary'],
        queryFn: getProfileSummary,
    })
}