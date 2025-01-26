import axios, { AxiosError, AxiosRequestConfig, CanceledError } from "axios";
import { IdentitiesServiceProxy } from "./identity-service-proxies";
import { AppConsts } from "@shared/app-consts";
import { CookieService } from "@shared/cookies/cookie-service";
import { SwalMessageService } from '@shared/sweetalert2';
import { AppAuthService } from "@shared/auth/app-auth-service";
import { CategoriesServiceProxy, ProductsServiceProxy } from "./product-service-proxies";
import { CartsServiceProxy } from "./cart-service-proxies";

export interface ApiException {
    status: number;
    detail: string;
}

interface RetryQueueItem {
    resolve: (value?: any) => void;
    reject: (error?: never) => void;
    config: AxiosRequestConfig;
}

let isRefreshing = false;

const refreshAndRetryQueue: RetryQueueItem[] = [];

const baseUrl = import.meta.env.VITE_REMOTE_SERVICE_BASE_URL;

const axiosInstance = axios.create();

axiosInstance.interceptors.request.use(function (config) {
    const userToken = CookieService.getCookie(AppConsts.cookieName.accessToken);

    if (userToken) {
        config.headers.Authorization = `Bearer ${userToken}`;
    }

    return config;
}, function (error) {
    // Do something with request error
    return Promise.reject(error);
});

axiosInstance.interceptors.response.use(function (response) {
    return response;
}, async function (error: AxiosError) {
    const originalRequest = error.config;

    // Check for 401 error
    if (error.response && error.response.status === 401 && originalRequest) {
        if (!isRefreshing) {
            isRefreshing = true;
            try {
                // Refresh the access token
                const authService = new AppAuthService();
                await authService.refreshToken();

                // Update the request headers with the new access token
                const newAccessToken = CookieService.getCookie(AppConsts.cookieName.accessToken);
                originalRequest.headers['Authorization'] = `Bearer ${newAccessToken}`;

                // Retry all requests in the queue with the new token
                refreshAndRetryQueue.forEach(({ config, resolve, reject }) => {
                    axiosInstance
                        .request(config)
                        .then((response) => resolve(response))
                        .catch((err) => reject(err));
                });

                // Clear the queue
                refreshAndRetryQueue.length = 0;

                // Retry the original request
                return axiosInstance(originalRequest);
            } finally {
                isRefreshing = false;
            }
        }

        // Add the original request to the queue
        return new Promise<void>((resolve, reject) => {
            refreshAndRetryQueue.push({ config: originalRequest, resolve, reject });
        });
    }

    return handleErrorResponse(error);
});

const handleErrorResponse = (error: AxiosError): Promise<never> => {
    if (error instanceof CanceledError) {
        return Promise.reject(error);
    }

    const apiError = error?.response?.data as ApiException;

    SwalMessageService.showError(apiError?.detail ?? error.message);

    return Promise.reject(error);
}

class APIClient {
    static getIdentityService(): IdentitiesServiceProxy {
        const service = new IdentitiesServiceProxy(baseUrl, axiosInstance);
        return service;
    }

    static getProductService(): ProductsServiceProxy {
        const service = new ProductsServiceProxy(baseUrl, axiosInstance);
        return service;
    }

    static getCartService(): CartsServiceProxy {
        const service = new CartsServiceProxy(baseUrl, axiosInstance);
        return service;
    }

    static getCategoryService(): CategoriesServiceProxy {
        const service = new CategoriesServiceProxy(baseUrl, axiosInstance);
        return service;
    }
}

export default APIClient