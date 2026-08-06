import {Link} from 'react-router-dom'
import {Logo} from '@/shared/ui/logo'
import styles from './Header.module.scss'
import {FragmentCounter} from './FragmentCounter'
import {PointsCounter} from './PointsCounter'
import {Avatar} from "@/shared/ui/avatar";
import {LoginButton} from "@/features/login-button";
// import {PlayButton} from "@/features/play-button";

export const Header = () => {
    return (
        <header className={styles.header}>
            <Link to="/">
                <Logo/>
            </Link>

            <div className={styles.left}>
                <nav className={styles.navigation}>
                    <Link to='/'>Сценарии</Link>
                    <Link to='/glossary'>Глоссарий</Link>
                    <Link to='/leaderboard'>Лидеры</Link>
                    <Link to='/dashboard'>Статистика</Link>
                </nav>

                <div className={styles.info}>
                    <div className={styles.counters}>
                        <FragmentCounter value={0}/>
                        <PointsCounter value={0}/>
                    </div>

                    {/*<Link to={'/dashboard'}>*/}
                    {/*    <Avatar/>*/}
                    {/*</Link>*/}

                    <LoginButton/>

                    {/*<PlayButton/>*/}
                </div>

            </div>
        </header>
    )
}
