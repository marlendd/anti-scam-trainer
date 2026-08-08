import type { ReactNode } from 'react'
import { QueryProvider } from './QueryProvider'
import { ReduxProvider } from './ReduxProvider'

type AppProvidersProps = {
    children: ReactNode
}

export function AppProviders({ children }: AppProvidersProps) {
    return (
        <ReduxProvider>
            <QueryProvider>{children}</QueryProvider>
        </ReduxProvider>
    )
}
