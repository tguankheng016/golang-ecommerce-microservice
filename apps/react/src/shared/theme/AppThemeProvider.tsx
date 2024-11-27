import { ReactNode, useEffect } from "react";
import { useThemeStore, Theme } from '@shared/theme';
import { AppConsts } from "@shared/app-consts";
import useLayoutStore from "./layout-store";

interface Props {
    children: ReactNode;
}

const AppThemeProvider = ({ children }: Props) => {
    const { setTheme } = useThemeStore();
    const { setIsMobile } = useLayoutStore();

    useEffect(() => {
        let savedTheme = localStorage.getItem(AppConsts.localStorage.theme) as Theme;

        if (!savedTheme) {
            savedTheme = 'dark';
        }

        if (savedTheme) {
            setTheme(savedTheme);
        }
    }, [setTheme]);

    useEffect(() => {
        const handleResize = () => {
            setIsMobile(window.innerWidth < 768);
        };

        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, []);

    return (
        <>{children}</>
    )
}

export default AppThemeProvider