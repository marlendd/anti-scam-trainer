import {Link} from 'react-router-dom'
import {Logo} from '@/shared/ui/logo'
import styles from './Header.module.scss'
import {LoginButton} from '@/features/login-button'
import {BurgerNavigation} from '@/widgets/header'
import {useCurrentUser} from '@/entities/user'
import {Avatar} from '@/shared/ui/avatar'
import {HeaderCounters} from "@/widgets/header/ui/HeaderCounters.tsx";

export const Header = () => {
    const {data: user, isPending, isError} = useCurrentUser()

    return (
        <header className={styles.header}>
            <Link to="/">
                <Logo/>
            </Link>

            <div className={styles.left}>
                <nav className={styles.navigation}>
                    <Link to="/training/role-selection">Сценарии</Link>
                    <Link to="/glossary">Глоссарий</Link>
                    <Link to="/leaderboard">Лидеры</Link>
                    {!isPending && !isError && user !== null && (
                        <Link to="/dashboard">Статистика</Link>
                    )}
                </nav>

                <div className={styles.info}>
                    {!isPending && !isError && user !== null && (
                        <div className={styles.counters}>
                            <HeaderCounters/>
                        </div>)
                    }

                    {isPending && <LoginButton/>}

                    {!isPending && !isError && user === null && <LoginButton/>}

                    {!isPending && !isError && user !== null && (
                        <Link to={'/dashboard'}>
                            <Avatar/>
                        </Link>
                    )}

                    {isError && <LoginButton/>}

                    <BurgerNavigation/>
                </div>
            </div>
        </header>
    )
}
