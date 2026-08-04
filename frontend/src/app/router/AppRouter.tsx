import { Navigate, Route, Routes } from 'react-router-dom'
import { DashboardPage } from '@/pages/dashboard'
import { GlossaryPage } from '@/pages/glossary'
import { HomePage } from '@/pages/home'
import { LeaderboardPage } from '@/pages/leaderboard'
import { WelcomePage } from '@/pages/welcome'

export function AppRouter() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/home" replace />} />
      <Route path="/home" element={<HomePage />} />
      <Route path="/welcome" element={<WelcomePage />} />
      <Route path="/dashboard" element={<DashboardPage />} />
      <Route path="/glossary" element={<GlossaryPage />} />
      <Route path="/leaderboard" element={<LeaderboardPage />} />
    </Routes>
  )
}
