import {Navigate, Route, Routes} from 'react-router-dom'

import {AuthLayout} from '@/app/layouts/auth-layout'
import {DefaultLayout} from '@/app/layouts/default'
import {GameLayout} from '@/app/layouts/game-layout'

import {DashboardPage} from '@/pages/dashboard'
import {ForgotPasswordPage} from '@/pages/forgot-password'
import {GlossaryPage} from '@/pages/glossary'
import {LeaderboardPage} from '@/pages/leaderboard'
import {LoginPage} from '@/pages/login'
import {NotFoundPage} from '@/pages/not-found'
import {RegisterPage} from '@/pages/register'
import {RoleSelectionPage} from '@/pages/role-selection'
import {ScamSchemePage} from '@/pages/scam-scheme'
import {ScamOrNotPage} from '@/pages/scam-or-not'
import {WelcomePage} from '@/pages/welcome'

import {BuyerPathPage} from '../../pages/buyer-path/ui/BuyerPathPage.tsx'
import {PuzzlePage} from '@/pages/puzzle/ui/PuzzlePage.tsx'

import {RequireAuth} from './guards/RequireAuth'

export function AppRouter() {
    return (
        <Routes>
            <Route path="/" element={<Navigate to="/home" replace/>}/>

            <Route element={<AuthLayout/>}>
                <Route path="/login" element={<LoginPage/>}/>
                <Route path="/register" element={<RegisterPage/>}/>
                <Route path="/forgot-password" element={<ForgotPasswordPage/>}/>
            </Route>

            <Route element={<RequireAuth/>}>
                <Route element={<GameLayout/>}>
                    <Route path="/training">
                        <Route path="role-selection" element={<RoleSelectionPage/>}/>
                        <Route path="scam-or-not" element={<ScamOrNotPage/>}/>
                        <Route
                            path="scam-or-not/:logicalScenarioId"
                            element={<ScamOrNotPage/>}
                        />

                        <Route path="path/:pathId">
                            <Route index element={<BuyerPathPage/>}/>
                            <Route path=":schemeId" element={<BuyerPathPage/>}/>
                        </Route>
                    </Route>
                </Route>
            </Route>

            <Route element={<DefaultLayout/>}>
                <Route path="/home" element={<WelcomePage/>}/>
                <Route path="/welcome" element={<WelcomePage/>}/>
                <Route path="/leaderboard" element={<LeaderboardPage/>}/>

                <Route element={<RequireAuth/>}>
                    <Route path="/dashboard" element={<DashboardPage/>}/>
                    <Route path="/puzzle" element={<PuzzlePage/>}/>
                </Route>

                <Route path="/glossary">
                    <Route index element={<GlossaryPage/>}/>
                    <Route path=":schemeId" element={<ScamSchemePage/>}/>
                </Route>

                <Route path="*" element={<NotFoundPage/>}/>
            </Route>
        </Routes>
    )
}
