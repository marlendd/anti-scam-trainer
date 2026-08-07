import {buyerSchemesMock} from "@/entities/training-path";
import {TrainingSchemeList} from "@/widgets/training-scheme-list";
import styles from './BuyerPathPage.module.scss';

import {useState} from 'react'

import {
    fakeDeliveryScenario,
    type TrainingScenario,
} from '@/entities/training-scenario'
import {useDocumentTitle} from '@/shared/lib/use-document-title'
import {TrainingChat} from '@/widgets/training-chat'
import {useMatch} from "react-router-dom";


export function BuyerPathPage() {
    const [scenario] = useState<TrainingScenario>(fakeDeliveryScenario)

    const sessionMatch = useMatch({
        path: '/training/path/:pathId/:schemeId',
        end: true,
    })

    useDocumentTitle(scenario.title)

    return (
        <main className={styles.page}>
            <TrainingSchemeList
                schemes={buyerSchemesMock}
            />

            <TrainingChat
                mode='interactive'
                scenario={sessionMatch ? scenario : null}
                onScenarioComplete={() => {

                }}
            />
        </main>
    )
}