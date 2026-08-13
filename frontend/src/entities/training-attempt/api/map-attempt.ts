import type {
    TrainingAttempt,
    TrainingAttemptState,
} from '../model/types'
import type {
    AttemptResponseDto,
    AttemptStateResponseDto,
} from './types'

export function mapAttempt(
    dto: AttemptResponseDto,
): TrainingAttempt {
    return {
        id: dto.id,
        scenarioId: dto.scenario_id,
        status: dto.status,
        currentNodeId: dto.current_node_id,
        endingId: dto.ending_id,
        ending: dto.ending,
        score: dto.score,
        startedAt: dto.started_at,
        updatedAt: dto.updated_at,
        completedAt: dto.completed_at,
    }
}

export function mapAttemptState(
    dto: AttemptStateResponseDto,
): TrainingAttemptState {
    return {
        ...mapAttempt(dto),

        scenario: {
            title: dto.scenario.title,
            description: dto.scenario.description,
            role: dto.scenario.role,
            product: {
                title: dto.scenario.product.title,
                price: dto.scenario.product.price,
            },
        },

        currentNode: dto.current_node
            ? {
                  id: dto.current_node.id,
                  author: dto.current_node.author,
                  text: dto.current_node.text,

                  messages: dto.current_node.messages.map((message) => ({
                      author: message.author,
                      text: message.text,
                  })),

                  choices: dto.current_node.choices.map((choice) => ({
                      id: choice.id,
                      text: choice.text,
                  })),
              }
            : undefined,

        history: dto.history.map((item) => ({
            node: {
                id: item.node.id,
                author: item.node.author,
                text: item.node.text,

                messages: item.node.messages.map((message) => ({
                    author: message.author,
                    text: message.text,
                })),
            },

            selectedChoice: {
                id: item.selected_choice.id,
                text: item.selected_choice.text,
            },

            consequence: item.consequence,
            answeredAt: item.answered_at,
        })),
    }
}