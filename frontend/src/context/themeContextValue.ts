import { createContext } from "react";
import type { Game, ThemeType } from "../types/app";

export interface ThemeContextValue {
    theme: ThemeType;
    setTheme: (theme: ThemeType) => void;
    switchThemeForGame: (game: Game) => void;
    particlesEnabled: boolean;
    setParticlesEnabled: (enabled: boolean) => void;
}

export const ThemeContext = createContext<ThemeContextValue | null>(null);
