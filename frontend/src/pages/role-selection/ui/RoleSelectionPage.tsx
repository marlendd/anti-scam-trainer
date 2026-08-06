import {useNavigate} from 'react-router-dom'
// import {Button} from '@/shared/ui/button'

import buyerImage from '@/shared/assets/images/buyer-path.webp'
import sellerImage from '@/shared/assets/images/seller-path.webp'

import styles from './RoleSelectionPage.module.scss'

type RoleCardProps = {
    title: string
    description: string
    image: string
    imageAlt: string
    accent: 'blue' | 'red' | 'green'
    onSelect: () => void
}

function RoleCard({
                      title,
                      description,
                      image,
                      imageAlt,
                      accent,
                      onSelect,
                  }: RoleCardProps) {
    return (
        <article
            onClick={onSelect}
            className={styles.card}
            data-accent={accent}
        >
            <div className={styles.imageWrapper}>
                <img
                    className={styles.image}
                    src={image}
                    alt={imageAlt}
                    width={420}
                    height={420}
                />
            </div>

            <div className={styles.content}>
                <h2 className={styles.cardTitle}>{title}</h2>

                <p className={styles.cardDescription}>
                    {description}
                </p>

                {/*<Button*/}
                {/*  fullWidth*/}
                {/*  onClick={onSelect}*/}
                {/*>*/}
                {/*  Выбрать*/}
                {/*</Button>*/}
            </div>
        </article>
    )
}

export function RoleSelectionPage() {
    const navigate = useNavigate()

    function handleSelectBuyer() {
        navigate('/training/path/buyer')
    }

    function handleSelectSeller() {
        navigate('/training/path/seller')
    }

    return (
        <main className={styles.page}>
            <section className={styles.hero}>

                {/*<h1 className={styles.title}>*/}
                {/*    С какой стороны начнём?*/}
                {/*</h1>*/}

                <p className={styles.description}>
                    Выберите путь, который хотите пройти первым.
                    У покупателя и продавца свои сценарии, риски и
                    мошеннические схемы.
                </p>
            </section>

            <section className={styles.grid}>
                <RoleCard
                    title="Покупатель"
                    description="Научитесь безопасно покупать товары и проверять продавцов."
                    image={buyerImage}
                    imageAlt="Иллюстрация пути покупателя"
                    accent="blue"
                    onSelect={handleSelectBuyer}
                />

                <RoleCard
                    title="Продавец"
                    description="Научитесь безопасно продавать товары и проверять клиентов."
                    image={sellerImage}
                    imageAlt="Иллюстрация пути продавца"
                    accent="green"
                    onSelect={handleSelectSeller}
                />
            </section>
        </main>
    )
}