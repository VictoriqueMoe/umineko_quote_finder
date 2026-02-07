import {type ReactNode, useCallback, useLayoutEffect, useState} from "react";
import {Chart as ChartJS} from "chart.js";
import type {ThemeType} from "../types/app";
import {ThemeContext} from "./themeContextValue";

const STORAGE_KEY = "uq-theme";
const PARTICLES_KEY = "uq-particles";
const DEFAULT_THEME: ThemeType = "featherine";

const THEME_CHART_COLOURS: Record<ThemeType, string> = {
    featherine: "#a89bb8",
    bernkastel: "#9bb5d0",
    lambdadelta: "#d09bb8",
};

function getStoredTheme(): ThemeType {
    try {
        const stored = localStorage.getItem(STORAGE_KEY);
        if (stored === "featherine" || stored === "bernkastel" || stored === "lambdadelta") {
            return stored;
        }
    } catch {
        // localStorage unavailable
    }
    return DEFAULT_THEME;
}

function getStoredParticles(): boolean {
    try {
        const stored = localStorage.getItem(PARTICLES_KEY);
        if (stored !== null) {
            return stored === "true";
        }
    } catch {
        // localStorage unavailable
    }
    return true;
}

export function ThemeProvider({ children }: { children: ReactNode }) {
    const [theme, setThemeState] = useState<ThemeType>(getStoredTheme);
    const [particlesEnabled, setParticlesEnabledState] = useState(getStoredParticles);

    useLayoutEffect(() => {
        if (theme === DEFAULT_THEME) {
            document.documentElement.removeAttribute("data-theme");
        } else {
            document.documentElement.setAttribute("data-theme", theme);
        }
        ChartJS.defaults.color = THEME_CHART_COLOURS[theme];
    }, [theme]);

    const setTheme = useCallback((newTheme: ThemeType) => {
        setThemeState(newTheme);
        try {
            localStorage.setItem(STORAGE_KEY, newTheme);
        } catch {
            // localStorage unavailable
        }
    }, []);

    const setParticlesEnabled = useCallback((enabled: boolean) => {
        setParticlesEnabledState(enabled);
        try {
            localStorage.setItem(PARTICLES_KEY, String(enabled));
        } catch {
            // localStorage unavailable
        }
    }, []);

    return (
        <ThemeContext.Provider value={{ theme, setTheme, particlesEnabled, setParticlesEnabled }}>
            {children}
        </ThemeContext.Provider>
    );
}
