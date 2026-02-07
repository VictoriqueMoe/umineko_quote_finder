import { useCallback, useState } from "react";
import * as api from "../api/endpoints";
import type { SearchResult } from "../types/api";
import type { FilterState, Language } from "../types/app";

export function useSearch() {
    const [results, setResults] = useState<SearchResult[]>([]);
    const [query, setQuery] = useState("");
    const [offset, setOffset] = useState(0);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const search = useCallback(
        async (
            q: string,
            language: Language,
            off: number = 0,
            filters?: FilterState,
        ): Promise<{ offset: number; total: number } | undefined> => {
            if (!q.trim()) {
                setError("Enter a search term to find quotes.");
                return undefined;
            }
            setLoading(true);
            setError(null);
            try {
                const data = await api.searchQuotes(
                    q,
                    language,
                    off,
                    filters?.character,
                    filters?.episode,
                    filters?.truth,
                );
                setResults(data.results || []);
                setQuery(q);
                setOffset(data.offset);
                setTotal(data.total);
                return { offset: data.offset, total: data.total };
            } catch {
                setError("Failed to search. Please try again.");
                return undefined;
            } finally {
                setLoading(false);
            }
        },
        [],
    );

    const clear = useCallback(() => {
        setResults([]);
        setQuery("");
        setOffset(0);
        setTotal(0);
        setError(null);
    }, []);

    return { results, query, offset, total, loading, error, search, clear };
}
