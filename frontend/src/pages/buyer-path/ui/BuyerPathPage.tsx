import { useMatch, useParams } from 'react-router-dom'

import { useTrainingSession } from '@/entities/training-attempt'
import { type TrainingRole, useScenarios } from '@/entities/training-scenario'
import { useDocumentTitle } from '@/shared/lib/use-document-title'
import { TrainingChat } from '@/widgets/training-chat'
import { TrainingSchemeList } from '@/widgets/training-scheme-list'

import styles from './BuyerPathPage.module.scss'

function isTrainingRole(value: string | undefined): value is TrainingRole {
    return value === 'buyer' || value === 'seller'
}

export function BuyerPathPage() {
    const { pathId } = useParams()

    const role = isTrainingRole(pathId) ? pathId : null

    const sessionMatch = useMatch({
        path: '/training/path/:pathId/:schemeId',
        end: true,
    })

    const scenarioId = sessionMatch?.params.schemeId ?? null

    const isSession = scenarioId !== null

    const {
        data: scenarios = [],
        isPending: isScenariosPending,
        isError: isScenariosError,
    } = useScenarios(role)

    const scenario = scenarios.find((scenario) => scenario.id === scenarioId) ?? null

    const {
        data: attempt,
        isPending: isSessionPending,
        isError: isSessionError,
    } = useTrainingSession(scenarioId)

    useDocumentTitle(scenario?.title ?? 'Тренировка')

    if (!role) {
        return <div>Неизвестная роль</div>
    }

    if (isScenariosPending) {
        return <div>Загрузка...</div>
    }

    if (isScenariosError) {
        return <div>Не удалось загрузить сценарии</div>
    }

    if (isSession && !scenario) {
        return <div>Сценарий не найден</div>
    }

    if (isSessionPending) {
        return <div>Загрузка сценария...</div>
    }

    if (isSessionError) {
        return <div>Не удалось запустить сценарий</div>
    }

    return (
        <main className={styles.page} data-session={isSession}>
            <div className={styles.pathColumn}>
                <TrainingSchemeList scenarios={scenarios} />
            </div>

            <div className={styles.chatColumn}>
                {/*
                    Здесь единственное место, которое я сейчас
                    не могу корректно дописать вслепую.

                    Нужные реальные данные чата теперь лежат в:

                    attempt?.currentNode

                    а твой TrainingChat сейчас всё ещё ожидает
                    старый prop `scenario`.
                */}

                <TrainingChat mode="interactive" attempt={attempt} onScenarioComplete={() => {}} />
            </div>
        </main>
    )
}
