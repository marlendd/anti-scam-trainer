import { Navigate, Outlet } from 'react-router-dom'

import { useCurrentUser } from '@/entities/user'
import {
    clearAuthSession,
    hasAuthSession,
} from '@/shared/lib/auth-session'

export function RequireAuth() {
    const hasSession = hasAuthSession()

    const {
        isError,
    } = useCurrentUser()

    if (!hasSession) {
        return <Navigate to="/login" replace />
    }

    if (isError) {
        clearAuthSession()

        return <Navigate to="/login" replace />
    }

    return <Outlet />
}