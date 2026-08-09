import { useMutation, useQueryClient } from '@tanstack/react-query'

import {
    submitAnswer,
    type SubmitAnswerResult,
} from './submit-answer'

type SubmitAnswerVariables = {
    attemptId: string
    nodeId: string
    choiceId: string
}

export function useSubmitAnswer() {
    const queryClient = useQueryClient()

    return useMutation<
        SubmitAnswerResult,
        Error,
        SubmitAnswerVariables
    >({
        mutationFn: submitAnswer,

        onSuccess: async (_, variables) => {
            await queryClient.invalidateQueries({
                queryKey: [
                    'training-attempt',
                    'state',
                    variables.attemptId,
                ],
            })
        },
    })
}