import { useCallback, useState } from "react";
import * as api from "../api/endpoints";
import type { Quote } from "../types/api";
import type { FilterState, Language } from "../types/app";

export function useFeaturedQuote() {
    const [quote, setQuote] = useState<Quote | null>(null);
    const [currentAudioId, setCurrentAudioId] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const randomQuote = useCallback(
        async (language: Language, filters: FilterState): Promise<{ audioId: string | null } | undefined> => {
            setLoading(true);
            setError(null);
            try {
                const q = await api.getRandomQuote(language, filters.character, filters.episode, filters.truth);
                if ("error" in q) {
                    setError("No quotes found for this character.");
                    return undefined;
                }
                setQuote(q);
                const firstId = q.audioId ? q.audioId.split(", ")[0] : null;
                setCurrentAudioId(firstId);
                return { audioId: firstId };
            } catch {
                setError("Failed to retrieve a quote.");
                return undefined;
            } finally {
                setLoading(false);
            }
        },
        [],
    );

    const lookupByAudioId = useCallback(
        async (audioId: string, language: Language): Promise<{ audioId: string } | undefined> => {
            setLoading(true);
            setError(null);
            try {
                const q = await api.getQuoteByAudioId(audioId, language);
                if ("error" in q) {
                    setError(`No quote found for audio ID "${audioId}".`);
                    return undefined;
                }
                setQuote(q);
                setCurrentAudioId(audioId);
                return { audioId };
            } catch {
                setError("Failed to look up audio ID. Please try again.");
                return undefined;
            } finally {
                setLoading(false);
            }
        },
        [],
    );

    const clear = useCallback(() => {
        setQuote(null);
        setCurrentAudioId(null);
        setError(null);
    }, []);

    return { quote, currentAudioId, loading, error, randomQuote, lookupByAudioId, clear };
}
