import { useCallback, useState } from "react";
import * as api from "../api/endpoints";
import type { BrowseResponse } from "../types/api";
import type { Language } from "../types/app";

export function useBrowse() {
    const [data, setData] = useState<BrowseResponse | null>(null);
    const [offset, setOffset] = useState(0);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const browse = useCallback(
        async (
            characterId: string,
            language: Language,
            off: number = 0,
            episode?: string,
            truth?: string,
        ): Promise<{ offset: number; total: number } | undefined> => {
            setLoading(true);
            setError(null);
            try {
                const result = await api.browseDialogue(language, off, characterId, episode, truth);
                setData(result);
                setOffset(result.offset);
                setTotal(result.total);
                return { offset: result.offset, total: result.total };
            } catch {
                setError("Failed to load dialogue.");
                return undefined;
            } finally {
                setLoading(false);
            }
        },
        [],
    );

    const clear = useCallback(() => {
        setData(null);
        setOffset(0);
        setTotal(0);
        setError(null);
    }, []);

    return { data, offset, total, loading, error, browse, clear };
}
