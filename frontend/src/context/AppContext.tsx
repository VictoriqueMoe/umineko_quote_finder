import { type ReactNode, useEffect, useState } from "react";
import type { Language } from "../types/app";
import type { CharactersResponse } from "../types/api";
import { getCharacters, getConfig } from "../api/endpoints";
import { AppContext } from "./appContextValue";

export function AppProvider({ children }: { children: ReactNode }) {
    const [language, setLanguage] = useState<Language>("en");
    const [hasAudio, setHasAudio] = useState(true);
    const [characters, setCharacters] = useState<CharactersResponse>({});
    const [sortedCharacters, setSortedCharacters] = useState<[string, string][]>([]);

    useEffect(() => {
        getConfig()
            .then(config => {
                setHasAudio(config.hasAudio);
            })
            .catch(() => {
                console.warn("Failed to load config");
            });

        getCharacters()
            .then(chars => {
                setCharacters(chars);
                const sorted = Object.entries(chars).sort((a, b) => a[1].localeCompare(b[1]));
                setSortedCharacters(sorted);
            })
            .catch(err => {
                console.error("Failed to load characters:", err);
            });
    }, []);

    return (
        <AppContext.Provider value={{ language, setLanguage, hasAudio, characters, sortedCharacters }}>
            {children}
        </AppContext.Provider>
    );
}
