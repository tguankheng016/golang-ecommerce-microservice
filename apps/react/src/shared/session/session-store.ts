import { AppAuthService } from "@shared/auth/app-auth-service";
import APIClient from "@shared/service-proxies/api-client";
import { UserLoginInfoDto } from "@shared/service-proxies/identity-service-proxies";
import { create } from "zustand";

interface SessionState {
    user?: UserLoginInfoDto;
    grantedPermissions: { [key: string]: boolean };
    loading: boolean;
    fetchCurrentUser: (signal: AbortSignal) => Promise<void>;
    isGranted(permissionName: string): boolean;
}

const useSessionStore = create<SessionState>((set, get) => ({
    user: undefined,
    grantedPermissions: {},
    loading: true,
    fetchCurrentUser: (signal: AbortSignal) => {
        const identityService = APIClient.getIdentityService();
        set({ loading: true });
        return identityService.getCurrentSession(signal)
            .then((res) => {
                if (!res.user || !res.user.id) {
                    // Try refresh
                    const authService = new AppAuthService();
                    authService.refreshToken()
                        .then((refreshRes) => {
                            if (refreshRes) {
                                identityService.getCurrentSession(signal)
                                    .then((secondRes) => {
                                        set({ user: secondRes.user });
                                        set({ loading: false });
                                    })
                            } else {
                                set({ user: res.user });
                                set({ loading: false });
                            }
                        });
                } else {
                    set({ user: res.user });
                    set({ grantedPermissions: res.grantedPermissions });
                    set({ loading: false });
                }
            });
    },
    isGranted: (permissionName: string) => {
        const { grantedPermissions } = get();
        return permissionName in grantedPermissions;
    }
}));

export default useSessionStore;