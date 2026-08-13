import {
    Navigate,
    Outlet,
    useLocation,
} from 'react-router-dom'

import { useCurrentUser } from '@/entities/user'
import {
    clearAuthSession,
    hasAuthSession,
} from '@/shared/lib/auth-session'

export function RequireAuth() {
    const location = useLocation()

    const hasSession = hasAuthSession()

    const {
        isError,
        isPending,
    } = useCurrentUser()

    const returnTo =
        location.pathname +
        location.search +
        location.hash

    const loginUrl = `/login?returnTo=${encodeURIComponent(
        returnTo,
    )}`

    if (!hasSession) {
        return (
            <Navigate
                to={loginUrl}
                replace
            />
        )
    }

    if (isPending) {
        return null
    }

    if (isError) {
        clearAuthSession()

        return (
            <Navigate
                to={loginUrl}
                replace
            />
        )
    }

    return <Outlet />
}