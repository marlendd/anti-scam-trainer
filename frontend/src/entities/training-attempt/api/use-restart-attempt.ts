import { useMutation } from '@tanstack/react-query'

import { restartAttempt } from './restart-attempt'

export function useRestartAttempt() {
    return useMutation({
        mutationFn: restartAttempt,
    })
}