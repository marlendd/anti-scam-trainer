import {
    useMutation,
    useQueryClient,
} from '@tanstack/react-query'

import {
    activeAttemptQueryKey,
} from './use-training-session'
import { restartAttempt } from './restart-attempt'

export function useRestartAttempt() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: restartAttempt,

        onSuccess: async (
            _,
            scenarioId,
        ) => {
            /*
             * Очищаем старое состояние attempt.
             *
             * Важно: restart может создать новый attempt,
             * а может сбросить существующий с тем же id.
             * Поэтому инвалидируем и active, и state.
             */
            await queryClient.invalidateQueries({
                queryKey: [
                    'training-attempt',
                    'state',
                ],
            })

            await queryClient.invalidateQueries({
                queryKey:
                    activeAttemptQueryKey(
                        scenarioId,
                    ),
            })
        },
    })
}