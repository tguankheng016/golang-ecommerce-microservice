import { AccountLayout } from "@account/components/layout"
import { CallbackLoginPage } from "@account/pages/callback-login";
import { LoginPage } from "@account/pages/login";
import { AppLayout } from "@app/components/layout";
import { RolePage } from "@app/pages/admin/roles";
import { UserPage } from "@app/pages/admin/users";
import { FlightPage } from "@app/pages/flights";
import { HomePage } from "@app/pages/home";
import { AccountRoute, AppProtectedRoute, Error404Page } from "@shared/components/routing";
import { createBrowserRouter, Navigate, RouteObject } from "react-router-dom";

const accountRoutes: RouteObject[] = [
    {
        path: '',
        element: <AccountLayout />,
        children: [
            { index: true, element: <Navigate to="login" /> },
            { path: 'login', element: <AccountRoute><LoginPage /></AccountRoute> },
            { path: 'callback/login', element: <CallbackLoginPage /> }
        ]
    }
];

const appRoutes: RouteObject[] = [
    {
        path: '',
        element: <AppLayout />,
        children: [
            { index: true, element: <Navigate to="home" /> },
            { path: 'home', element: <AppProtectedRoute><HomePage /></AppProtectedRoute> },
            { path: 'flights', element: <AppProtectedRoute><FlightPage /></AppProtectedRoute> },
            { path: 'admin/users', element: <AppProtectedRoute><UserPage /></AppProtectedRoute> },
            { path: 'admin/roles', element: <AppProtectedRoute><RolePage /></AppProtectedRoute> }
        ]
    }
];

const routes: RouteObject[] = [
    {
        path: '/',
        errorElement: <Error404Page />,
        children: [
            { index: true, element: <Navigate to="app" /> },
            { path: 'app', children: appRoutes },
            { path: 'account', children: accountRoutes },
        ]
    }
];

const router = createBrowserRouter(routes);

export default router