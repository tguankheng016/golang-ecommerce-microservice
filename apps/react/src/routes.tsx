/* eslint-disable react-refresh/only-export-components */
import { LoadingScreen } from "@shared/components/loading-screen";
import { AccountRoute, AppProtectedRoute, Error404Page } from "@shared/components/routing";
import { lazy, Suspense } from "react";
import { createBrowserRouter, Navigate, RouteObject } from "react-router-dom";

// Lazy load the components
const UserPage = lazy(() => import('@app/pages/admin/users').then(module => ({ default: module.UserPage })));
const RolePage = lazy(() => import('@app/pages/admin/roles').then(module => ({ default: module.RolePage })));

const adminRoutes: RouteObject[] = [
    {
        path: '',
        children: [
            {
                path: 'users',
                element: (
                    <Suspense fallback={<LoadingScreen />}>
                        <AppProtectedRoute><UserPage /></AppProtectedRoute>
                    </Suspense>
                )
            },
            {
                path: 'roles',
                element: (
                    <Suspense fallback={<LoadingScreen />}>
                        <AppProtectedRoute><RolePage /></AppProtectedRoute>
                    </Suspense>
                )
            },
        ]
    }
]

const CategoryPage = lazy(() => import('@app/pages/main/categories').then(module => ({ default: module.CategoryPage })));
const ProductPage = lazy(() => import('@app/pages/main/products').then(module => ({ default: module.ProductPage })));

const mainRoutes: RouteObject[] = [
    {
        path: '',
        children: [
            {
                path: 'categories',
                element: (
                    <Suspense fallback={<LoadingScreen />}>
                        <AppProtectedRoute><CategoryPage /></AppProtectedRoute>
                    </Suspense>
                )
            },
            {
                path: 'products',
                element: (
                    <Suspense fallback={<LoadingScreen />}>
                        <AppProtectedRoute><ProductPage /></AppProtectedRoute>
                    </Suspense>
                )
            },
        ]
    }
]

const AccountLayout = lazy(() => import('@account/components/layout').then(module => ({ default: module.AccountLayout })));
const LoginPage = lazy(() => import('@account/pages/login').then(module => ({ default: module.LoginPage })));
const CallbackLoginPage = lazy(() => import('@account/pages/callback-login').then(module => ({ default: module.CallbackLoginPage })));

const accountRoutes: RouteObject[] = [
    {
        path: '',
        element: (
            <Suspense fallback={<LoadingScreen />}>
                <AccountLayout />
            </Suspense>
        ),
        children: [
            { index: true, element: <Navigate to="login" /> },
            {
                path: 'login',
                element: (
                    <Suspense fallback={<LoadingScreen />}>
                        <AccountRoute><LoginPage /></AccountRoute>
                    </Suspense>
                )
            },
            {
                path: 'callback/login',
                element: (
                    <Suspense fallback={<LoadingScreen />}>
                        <CallbackLoginPage />
                    </Suspense>
                )
            },
        ]
    }
];

const AppLayout = lazy(() => import('@app/components/layout').then(module => ({ default: module.AppLayout })));
const HomePage = lazy(() => import('@app/pages/home').then(module => ({ default: module.HomePage })));

const appRoutes: RouteObject[] = [
    {
        path: '',
        element: (
            <Suspense fallback={<LoadingScreen />}>
                <AppLayout />
            </Suspense>
        ),
        children: [
            { index: true, element: <Navigate to="home" /> },
            {
                path: 'home',
                element: (
                    <Suspense fallback={<LoadingScreen />}>
                        <AppProtectedRoute><HomePage /></AppProtectedRoute>
                    </Suspense>
                )
            },
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