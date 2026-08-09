import { useMutation, useQueryClient } from '@tanstack/react-query'

import { submitAnswer } from './submit-answer'

export function useSubmitAnswer() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: submitAnswer,

        onSuccess: async (result, variables) => {
            await queryClient.invalidateQueries({
                queryKey: [
                    'training-attempt',
                    'state',
                    variables.attemptId,
                ],
            })

            if (!result.completed) {
                return
            }

            await Promise.all([
                queryClient.invalidateQueries({
                    queryKey: ['profile-progress'],
                }),

                queryClient.invalidateQueries({
                    queryKey: ['leaderboard'],
                }),
            ])
        },
    })
}