import { type ReactNode, useCallback, useEffect, useRef, useState } from "react";
import type { Game, Language } from "../types/app";
import type { CharactersResponse } from "../types/api";
import { getCharacters, getConfig } from "../api/endpoints";
import { AppContext } from "./appContextValue";

function parseInitialGame(): Game {
    const params = new URLSearchParams(window.location.search);
    const g = params.get("game");
    if (g === "higurashi") {
        return "higurashi";
    }
    return "umineko";
}

const EMPTY_CHARS: CharactersResponse = { characters: {}, additional: {} };

export function AppProvider({ children }: { children: ReactNode }) {
    const [language, setLanguage] = useState<Language>("auto");
    const [game, setGameState] = useState<Game>(parseInitialGame);
    const [hasAudio, setHasAudio] = useState(true);
    const [characters, setCharacters] = useState<CharactersResponse>(EMPTY_CHARS);
    const [sortedCharacters, setSortedCharacters] = useState<[string, string][]>([]);
    const [sortedAdditionalCharacters, setSortedAdditionalCharacters] = useState<[string, string][]>([]);

    const initialGameRef = useRef(game);

    const loadCharacters = useCallback((g: Game) => {
        getCharacters(g)
            .then(resp => {
                setCharacters(resp);
                const sorted = Object.entries(resp.characters || {}).sort((a, b) => a[1].localeCompare(b[1]));
                setSortedCharacters(sorted);
                const sortedExtra = Object.entries(resp.additional || {}).sort((a, b) => a[1].localeCompare(b[1]));
                setSortedAdditionalCharacters(sortedExtra);
            })
            .catch(err => {
                console.error("Failed to load characters:", err);
            });
    }, []);

    useEffect(() => {
        getConfig()
            .then(config => {
                setHasAudio(config.hasAudio);
            })
            .catch(() => {
                console.warn("Failed to load config");
            });

        loadCharacters(initialGameRef.current);
    }, [loadCharacters]);

    const setGame = useCallback(
        (g: Game) => {
            setGameState(g);
            loadCharacters(g);
        },
        [loadCharacters],
    );

    return (
        <AppContext.Provider
            value={{
                language,
                setLanguage,
                game,
                setGame,
                hasAudio,
                characters,
                sortedCharacters,
                sortedAdditionalCharacters,
            }}
        >
            {children}
        </AppContext.Provider>
    );
}
