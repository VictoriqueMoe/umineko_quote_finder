import { useCallback, useRef, useState } from "react";
import * as api from "../api/endpoints";
import type { HigurashiStatsResponse, StatsResponse } from "../types/api";
import type { Game } from "../types/app";

export function useStats() {
    const [data, setData] = useState<StatsResponse | HigurashiStatsResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const cache = useRef<Record<string, StatsResponse | HigurashiStatsResponse>>({});

    const loadStats = useCallback(async (game: Game, episode: string): Promise<void> => {
        setLoading(true);
        setError(null);
        try {
            const cacheKey = `${game}:ep${episode || "0"}`;
            if (!cache.current[cacheKey]) {
                cache.current[cacheKey] = await api.getStats(game, episode);
            }
            setData(cache.current[cacheKey]);
        } catch {
            setError("Failed to load statistics.");
        } finally {
            setLoading(false);
        }
    }, []);

    const clear = useCallback(() => {
        setData(null);
        setError(null);
    }, []);

    return { data, loading, error, loadStats, clear };
}
