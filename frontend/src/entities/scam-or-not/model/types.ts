export type ScamOrNotCard = {
    id: string
    topicId: string
    situation: string
    isScam: boolean
    explanation: string
    riskSigns: string[]
}

export type ScamOrNotTopic = {
    id: string
    title: string
}
