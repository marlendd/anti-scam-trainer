import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { AppProviders } from './app/providers/AppProviders'
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

GridModuleRegistry.registerModules([
  GridAllCommunityModule,
])

ChartsModuleRegistry.registerModules([
  ChartsAllCommunityModule,
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <AppProviders>
        <App />
      </AppProviders>
    </BrowserRouter>
  </StrictMode>,
)