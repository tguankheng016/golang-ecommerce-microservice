import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import router from './routes'
import { AppSessionProvider } from '@shared/session'
import { AppThemeProvider } from '@shared/theme'
import { PrimeReactProvider } from 'primereact/api';

import '@assets/css/style.bundle.css'
import '@assets/plugins/global/plugins.bundle.css'
import 'sweetalert2/dist/sweetalert2.min.css'
import './index.css'

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <AppThemeProvider>
            <PrimeReactProvider>
                <AppSessionProvider>
                    <RouterProvider router={router} />
                </AppSessionProvider>
            </PrimeReactProvider>
        </AppThemeProvider>
    </StrictMode>
)
