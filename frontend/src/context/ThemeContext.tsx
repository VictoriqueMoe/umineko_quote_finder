import { type ReactNode, useCallback, useLayoutEffect, useState } from "react";
import { Chart as ChartJS } from "chart.js";
import type { CiconiaTheme, Game, HigurashiTheme, ThemeType, UminekoTheme } from "../types/app";
import { ThemeContext } from "./themeContextValue";

const UMINEKO_STORAGE_KEY = "uq-theme-umineko";
const HIGURASHI_STORAGE_KEY = "uq-theme-higurashi";
const CICONIA_STORAGE_KEY = "uq-theme-ciconia";
const PARTICLES_KEY = "uq-particles";

const UMINEKO_THEMES: UminekoTheme[] = ["featherine", "bernkastel", "lambdadelta"];
const HIGURASHI_THEMES: HigurashiTheme[] = ["rika", "mion", "satoko"];
const CICONIA_THEMES: CiconiaTheme[] = ["miyao", "lingji", "stanislaw"];
const DEFAULT_UMINEKO: UminekoTheme = "featherine";
const DEFAULT_HIGURASHI: HigurashiTheme = "rika";
const DEFAULT_CICONIA: CiconiaTheme = "miyao";

const THEME_CHART_COLOURS: Record<ThemeType, string> = {
    featherine: "#a89bb8",
    bernkastel: "#9bb5d0",
    lambdadelta: "#d09bb8",
    rika: "#a8a8cc",
    mion: "#9bb8ad",
    satoko: "#b8a890",
    miyao: "#3b8fd6",
    lingji: "#e5a900",
    stanislaw: "#e0e0e0",
};

function isUminekoTheme(t: string): t is UminekoTheme {
    return UMINEKO_THEMES.includes(t as UminekoTheme);
}

function isHigurashiTheme(t: string): t is HigurashiTheme {
    return HIGURASHI_THEMES.includes(t as HigurashiTheme);
}

function isCiconiaTheme(t: string): t is CiconiaTheme {
    return CICONIA_THEMES.includes(t as CiconiaTheme);
}

function storageKeyForGame(game: Game): string {
    if (game === "higurashi") {
        return HIGURASHI_STORAGE_KEY;
    }
    if (game === "ciconia") {
        return CICONIA_STORAGE_KEY;
    }
    return UMINEKO_STORAGE_KEY;
}

function defaultThemeForGame(game: Game): ThemeType {
    if (game === "higurashi") {
        return DEFAULT_HIGURASHI;
    }
    if (game === "ciconia") {
        return DEFAULT_CICONIA;
    }
    return DEFAULT_UMINEKO;
}

function themeMatchesGame(theme: string, game: Game): boolean {
    if (game === "umineko") {
        return isUminekoTheme(theme);
    }
    if (game === "higurashi") {
        return isHigurashiTheme(theme);
    }
    return isCiconiaTheme(theme);
}

function getStoredThemeForGame(game: Game): ThemeType {
    try {
        const stored = localStorage.getItem(storageKeyForGame(game));
        if (stored && themeMatchesGame(stored, game)) {
            return stored as ThemeType;
        }
    } catch {
        // localStorage unavailable
    }
    return defaultThemeForGame(game);
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
    const game = params.get("game");
    if (game === "higurashi") {
        return "higurashi";
    }
    if (game === "ciconia") {
        return "ciconia";
    }
    return "umineko";
}

function storageKeyForTheme(theme: ThemeType): string {
    if (isUminekoTheme(theme)) {
        return UMINEKO_STORAGE_KEY;
    }
    if (isHigurashiTheme(theme)) {
        return HIGURASHI_STORAGE_KEY;
    }
    return CICONIA_STORAGE_KEY;
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
            localStorage.setItem(storageKeyForTheme(newTheme), newTheme);
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
