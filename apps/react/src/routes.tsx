import { AccountLayout } from "@account/components/layout"
import { CallbackLoginPage } from "@account/pages/callback-login";
import { LoginPage } from "@account/pages/login";
import { AppLayout } from "@app/components/layout";
import { RolePage } from "@app/pages/admin/roles";
import { UserPage } from "@app/pages/admin/users";
import { HomePage } from "@app/pages/home";
import { CategoryPage } from "@app/pages/main/categories";
import { ProductPage } from "@app/pages/main/products";
import { AccountRoute, AppProtectedRoute, Error404Page } from "@shared/components/routing";
import { createBrowserRouter, Navigate, RouteObject } from "react-router-dom";

const adminRoutes: RouteObject[] = [
    {
        path: '',
        children: [
            { path: 'users', element: <AppProtectedRoute><UserPage /></AppProtectedRoute> },
            { path: 'roles', element: <AppProtectedRoute><RolePage /></AppProtectedRoute> },
        ]
    }
]

const mainRoutes: RouteObject[] = [
    {
        path: '',
        children: [
            { path: 'categories', element: <AppProtectedRoute><CategoryPage /></AppProtectedRoute> },
            { path: 'products', element: <AppProtectedRoute><ProductPage /></AppProtectedRoute> },
        ]
    }
]

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
            { path: 'admin', children: adminRoutes },
            { path: 'main', children: mainRoutes }
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