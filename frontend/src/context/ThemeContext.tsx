import { type ReactNode, useCallback, useLayoutEffect, useState } from "react";
import { Chart as ChartJS } from "chart.js";
import type { Game, HigurashiTheme, ThemeType, UminekoTheme } from "../types/app";
import { ThemeContext } from "./themeContextValue";

const UMINEKO_STORAGE_KEY = "uq-theme-umineko";
const HIGURASHI_STORAGE_KEY = "uq-theme-higurashi";
const PARTICLES_KEY = "uq-particles";

const UMINEKO_THEMES: UminekoTheme[] = ["featherine", "bernkastel", "lambdadelta"];
const HIGURASHI_THEMES: HigurashiTheme[] = ["rika", "mion", "satoko"];
const DEFAULT_UMINEKO: UminekoTheme = "featherine";
const DEFAULT_HIGURASHI: HigurashiTheme = "rika";

const THEME_CHART_COLOURS: Record<ThemeType, string> = {
    featherine: "#a89bb8",
    bernkastel: "#9bb5d0",
    lambdadelta: "#d09bb8",
    rika: "#a8a8cc",
    mion: "#9bb8ad",
    satoko: "#b8a890",
};

function isUminekoTheme(t: string): t is UminekoTheme {
    return UMINEKO_THEMES.includes(t as UminekoTheme);
}

function isHigurashiTheme(t: string): t is HigurashiTheme {
    return HIGURASHI_THEMES.includes(t as HigurashiTheme);
}

function getStoredThemeForGame(game: Game): ThemeType {
    try {
        const key = game === "umineko" ? UMINEKO_STORAGE_KEY : HIGURASHI_STORAGE_KEY;
        const stored = localStorage.getItem(key);
        if (stored) {
            if (game === "umineko" && isUminekoTheme(stored)) {
                return stored;
            }
            if (game === "higurashi" && isHigurashiTheme(stored)) {
                return stored;
            }
        }
    } catch {
        // localStorage unavailable
    }
    return game === "umineko" ? DEFAULT_UMINEKO : DEFAULT_HIGURASHI;
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

function parseInitialGame(): Game {
    const params = new URLSearchParams(window.location.search);
    return params.get("game") === "higurashi" ? "higurashi" : "umineko";
}

export function ThemeProvider({ children }: { children: ReactNode }) {
    const [theme, setThemeState] = useState<ThemeType>(() => getStoredThemeForGame(parseInitialGame()));
    const [particlesEnabled, setParticlesEnabledState] = useState(getStoredParticles);

    useLayoutEffect(() => {
        if (theme === "featherine") {
            document.documentElement.removeAttribute("data-theme");
        } else {
            document.documentElement.setAttribute("data-theme", theme);
        }
        ChartJS.defaults.color = THEME_CHART_COLOURS[theme];
    }, [theme]);

    const setTheme = useCallback((newTheme: ThemeType) => {
        setThemeState(newTheme);
        try {
            const key = isUminekoTheme(newTheme) ? UMINEKO_STORAGE_KEY : HIGURASHI_STORAGE_KEY;
            localStorage.setItem(key, newTheme);
        } catch {
            // localStorage unavailable
        }
    }, []);

    const switchThemeForGame = useCallback((game: Game) => {
        const stored = getStoredThemeForGame(game);
        setThemeState(stored);
        if (stored === "featherine") {
            document.documentElement.removeAttribute("data-theme");
        } else {
            document.documentElement.setAttribute("data-theme", stored);
        }
        ChartJS.defaults.color = THEME_CHART_COLOURS[stored];
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
        <ThemeContext.Provider value={{ theme, setTheme, switchThemeForGame, particlesEnabled, setParticlesEnabled }}>
            {children}
        </ThemeContext.Provider>
    );
}
