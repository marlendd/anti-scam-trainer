import { apiRequest } from '@/shared/api'

export type SubmitAnswerResult = {
    attemptId: string
    nodeId: string
    choiceId: string
    consequence: string
    nextNodeId?: string
    endingId?: string
    completed: boolean
    score?: number
}

type SubmitAnswerRequestDto = {
    node_id: string
    choice_id: string
    idempotency_key: string
}

type SubmitAnswerResponseDto = {
    attempt_id: string
    node_id: string
    choice_id: string
    consequence: string
    next_node_id?: string
    ending_id?: string
    completed: boolean
    score?: number
}

type SubmitAnswerParams = {
    attemptId: string
    nodeId: string
    choiceId: string
}

export async function submitAnswer({
    attemptId,
    nodeId,
    choiceId,
}: SubmitAnswerParams): Promise<SubmitAnswerResult> {
    const body: SubmitAnswerRequestDto = {
        node_id: nodeId,
        choice_id: choiceId,
        idempotency_key: crypto.randomUUID(),
    }

    const response = await apiRequest<SubmitAnswerResponseDto>(
        `/attempts/${attemptId}/answers`,
        {
            method: 'POST',
            body,
        },
    )

    return {
        attemptId: response.attempt_id,
        nodeId: response.node_id,
        choiceId: response.choice_id,
        consequence: response.consequence,
        nextNodeId: response.next_node_id,
        endingId: response.ending_id,
        completed: response.completed,
        score: response.score,
    }
}