import { buyerSchemesMock } from "@/entities/training-path";
import {TrainingSchemeList} from "@/widgets/training-scheme-list";
import styles from './BuyerPathPage.module.scss';

export function BuyerPathPage() {
  return (
    <main className={styles.page}>
      <TrainingSchemeList
        schemes={buyerSchemesMock}
      />
    </main>
  )
}