import { useCallback, useState } from "react";
import * as api from "../api/endpoints";
import type { BrowseResponse } from "../types/api";
import type { Game, Language } from "../types/app";

export function useBrowse() {
    const [data, setData] = useState<BrowseResponse | null>(null);
    const [offset, setOffset] = useState(0);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const browse = useCallback(
        async (
            game: Game,
            characterId: string,
            language: Language,
            off: number = 0,
            interactionA?: string,
            interactionB?: string,
            episode?: string,
            truth?: string,
            arc?: string,
            chapter?: string,
        ): Promise<{ offset: number; total: number } | undefined> => {
            setLoading(true);
            setError(null);
            try {
                if (!!interactionA !== !!interactionB) {
                    setError("Select both interaction characters or clear both.");
                    return undefined;
                }
                const result = await api.browseDialogue(
                    game,
                    language,
                    off,
                    characterId,
                    interactionA,
                    interactionB,
                    episode,
                    truth,
                    arc,
                    chapter,
                );
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
