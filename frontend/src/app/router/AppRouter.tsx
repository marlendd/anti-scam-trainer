import {Navigate, Route, Routes} from 'react-router-dom'

import {AuthLayout} from '@/app/layouts/auth-layout'
import {DefaultLayout} from '@/app/layouts/default'

import {DashboardPage} from '@/pages/dashboard'
import {GlossaryPage} from '@/pages/glossary'
// import {HomePage} from '@/pages/home'
import {LeaderboardPage} from '@/pages/leaderboard'
import {LoginPage} from '@/pages/login'
import {NotFoundPage} from '@/pages/not-found'
import {RegisterPage} from '@/pages/register'
import {ScamSchemePage} from '@/pages/scam-scheme'
import {WelcomePage} from '@/pages/welcome'
import {RoleSelectionPage} from "@/pages/role-selection";
import {GameLayout} from "@/app/layouts/game-layout";
import {BuyerPathPage} from "../../pages/buyer-path/ui/BuyerPathPage.tsx";
import {ForgotPasswordPage} from "@/pages/forgot-password";

export function AppRouter() {
    return (
        <Routes>
            <Route path="/" element={<Navigate to="/home" replace/>}/>

            <Route element={<AuthLayout/>}>
                <Route path="/login" element={<LoginPage/>}/>
                <Route path="/register" element={<RegisterPage/>}/>
                <Route path="/forgot-password" element={<ForgotPasswordPage/>}/>
            </Route>

            <Route element={<GameLayout/>}>
                <Route path="/training">
                    <Route path="role-selection" element={<RoleSelectionPage/>}/>


                    <Route path="path/:pathId">
                        <Route index element={<BuyerPathPage/>}/>
                        <Route path=":schemeId" element={<BuyerPathPage/>}/>
                    </Route>

                </Route>

            </Route>

            <Route element={<DefaultLayout/>}>
                <Route path="/home" element={<WelcomePage/>}/>
                <Route path="/welcome" element={<WelcomePage/>}/>
                <Route path="/dashboard" element={<DashboardPage/>}/>
                <Route path="/leaderboard" element={<LeaderboardPage/>}/>

                <Route path="/glossary">
                    <Route index element={<GlossaryPage/>}/>
                    <Route path=":schemeId" element={<ScamSchemePage/>}/>
                </Route>

                <Route path="*" element={<NotFoundPage/>}/>
            </Route>

        </Routes>
    )
}