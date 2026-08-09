import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { AppProviders } from './app/providers/AppProviders'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import './app/styles/global.scss'
import App from './app/App'

import {
    AllCommunityModule as GridAllCommunityModule,
    ModuleRegistry as GridModuleRegistry,
} from 'ag-grid-community'

import {
    AllCommunityModule as ChartsAllCommunityModule,
    ModuleRegistry as ChartsModuleRegistry,
} from 'ag-charts-community'

GridModuleRegistry.registerModules([GridAllCommunityModule])

ChartsModuleRegistry.registerModules([ChartsAllCommunityModule])

const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <BrowserRouter>
                <AppProviders>
                    <App />
                </AppProviders>
            </BrowserRouter>
        </QueryClientProvider>
    </StrictMode>,
)
