import APIClient from "@shared/service-proxies/api-client";
import { UserLoginInfoDto } from "@shared/service-proxies/identity-service-proxies";
import { create } from "zustand";

interface SessionState {
    user?: UserLoginInfoDto;
    loading: boolean;
    fetchCurrentUser: (signal: AbortSignal) => Promise<void>;
}

const useSessionStore = create<SessionState>((set) => ({
    user: undefined,
    loading: true,
    fetchCurrentUser: (signal: AbortSignal) => {
        const identityService = APIClient.getIdentityService();
        set({ loading: true });
        return identityService.getCurrentSession(signal)
            .then((res) => {
                set({ user: res.user });
                set({ loading: false });
            });
    }
}));

export default useSessionStore;