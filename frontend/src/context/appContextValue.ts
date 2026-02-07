import { createContext } from "react";
import type { Language } from "../types/app";
import type { CharactersResponse } from "../types/api";

export interface AppContextValue {
    language: Language;
    setLanguage: (lang: Language) => void;
    hasAudio: boolean;
    characters: CharactersResponse;
    sortedCharacters: [string, string][];
}

export const AppContext = createContext<AppContextValue | null>(null);
