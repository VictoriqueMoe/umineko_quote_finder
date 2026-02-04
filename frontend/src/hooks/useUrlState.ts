import { useCallback, useEffect, useRef } from "react";
import type { FilterState, Language, PushUrlParams } from "../types/app";

interface UrlStateCallbacks {
    onSearch: (query: string, offset: number, character: string, lang: Language) => void;
    onBrowse: (character: string, offset: number, episode: string, lang: Language) => void;
    onStats: (lang: Language) => void;
    onQuoteLookup: (audioId: string, lang: Language) => void;
    onDefault: (lang: Language) => void;
    setLanguage: (lang: Language) => void;
    setFilters: (filters: Partial<FilterState>) => void;
}

export function useUrlState(callbacks: UrlStateCallbacks) {
    const callbacksRef = useRef(callbacks);
    useEffect(() => {
        callbacksRef.current = callbacks;
    });

    const loadFromURL = useCallback(() => {
        const params = new URLSearchParams(window.location.search);
        const cb = callbacksRef.current;

        const lang = (params.get("lang") || "en") as Language;
        cb.setLanguage(lang);

        const episode = params.get("episode") || "0";
        const truth = params.get("truth") || "";
        cb.setFilters({ episode, truth });

        const offset = parseInt(params.get("offset") || "0") || 0;

        if (params.get("stats") === "1") {
            cb.onStats(lang);
            return;
        }

        const quoteId = params.get("quote");
        if (quoteId) {
            cb.onQuoteLookup(quoteId, lang);
            return;
        }

        const browse = params.get("browse");
        if (browse) {
            const character = browse !== "1" ? browse : "";
            cb.setFilters({ character, episode, truth });
            cb.onBrowse(character, offset, episode, lang);
            return;
        }

        const q = params.get("q");
        if (q) {
            const character = params.get("character") || "";
            cb.setFilters({ character, episode, truth });
            cb.onSearch(q, offset, character, lang);
            return;
        }

        cb.onDefault(lang);
    }, []);

    useEffect(() => {
        loadFromURL();
        window.addEventListener("popstate", loadFromURL);
        return () => {
            window.removeEventListener("popstate", loadFromURL);
        };
    }, [loadFromURL]);
}

export function pushUrl(state: PushUrlParams, language: Language, searchQuery: string) {
    const params = new URLSearchParams();

    if (state.viewMode === "stats") {
        params.set("stats", "1");
    } else if (state.viewMode === "browse") {
        params.set("browse", state.filters.character || "1");
    } else if (state.viewMode === "search" && searchQuery.trim()) {
        params.set("q", searchQuery.trim());
        if (state.filters.character) {
            params.set("character", state.filters.character);
        }
    } else if (state.viewMode === "quoteLookup" && state.currentAudioId) {
        params.set("quote", state.currentAudioId);
    }

    if (state.filters.episode && state.filters.episode !== "0") {
        params.set("episode", state.filters.episode);
    }

    if (state.filters.truth) {
        params.set("truth", state.filters.truth);
    }

    const offset = state.viewMode === "browse" ? state.browseOffset : state.searchOffset;
    if (offset > 0) {
        params.set("offset", String(offset));
    }

    if (language !== "en") {
        params.set("lang", language);
    }

    const qs = params.toString();
    const newURL = qs ? `?${qs}` : window.location.pathname;
    history.pushState(null, "", newURL);
}
