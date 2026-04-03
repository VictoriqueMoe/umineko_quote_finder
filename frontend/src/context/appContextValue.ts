import { createContext } from "react";
import type { Game, Language } from "../types/app";
import type { CharactersResponse } from "../types/api";

export interface AppContextValue {
    language: Language;
    setLanguage: (lang: Language) => void;
    game: Game;
    setGame: (game: Game) => void;
    hasAudio: boolean;
    characters: CharactersResponse;
    sortedCharacters: [string, string][];
    sortedAdditionalCharacters: [string, string][];
}

export const AppContext = createContext<AppContextValue | null>(null);
