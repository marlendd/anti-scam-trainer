import { apiRequest } from '@/shared/api'

export type RiskProfile = {
    dominantRisk?: string
    riskCount?: number
    description?: string
}

export type AttemptFeedback = {
    strengths: string[]
    weaknesses: string[]
    riskProfile?: RiskProfile
    recommendations: string[]
    learningTips: string[]
    motivation?: string
    source: 'fallback' | 'ai'
}

type AttemptFeedbackResponseDto = {
    strengths?: string[]
    weaknesses?: string[]
    risk_profile?: {
        dominant_risk?: string
        risk_count?: number
        description?: string
    }
    recommendations?: string[]
    learning_tips?: string[]
    motivation?: string
    source: 'fallback' | 'ai'
}

export async function getAttemptFeedback(attemptId: string): Promise<AttemptFeedback> {
    const response = await apiRequest<AttemptFeedbackResponseDto>(`/attempts/${attemptId}/feedback`)

    return {
        strengths: response.strengths ?? [],
        weaknesses: response.weaknesses ?? [],
        riskProfile: response.risk_profile
            ? {
                  dominantRisk: response.risk_profile.dominant_risk,
                  riskCount: response.risk_profile.risk_count,
                  description: response.risk_profile.description,
              }
            : undefined,
        recommendations: response.recommendations ?? [],
        learningTips: response.learning_tips ?? [],
        motivation: response.motivation,
        source: response.source ?? 'fallback'
    }
}
