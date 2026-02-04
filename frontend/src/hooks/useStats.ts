import { useCallback, useRef, useState } from "react";
import * as api from "../api/endpoints";
import type { StatsResponse } from "../types/api";

export function useStats() {
    const [data, setData] = useState<StatsResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const cache = useRef<Record<string, StatsResponse>>({});

    const loadStats = useCallback(async (episode: string): Promise<void> => {
        setLoading(true);
        setError(null);
        try {
            const cacheKey = `ep${episode || "0"}`;
            if (!cache.current[cacheKey]) {
                cache.current[cacheKey] = await api.getStats(episode);
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
