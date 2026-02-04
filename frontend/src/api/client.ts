const API_BASE = "/api/v1";

export async function apiFetch<T>(path: string): Promise<T> {
    const response = await fetch(`${API_BASE}${path}`);
    if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
    }
    return response.json();
}

export function buildQueryString(params: Record<string, string | number | undefined>): string {
    const search = new URLSearchParams();
    for (const [key, value] of Object.entries(params)) {
        if (value !== undefined && value !== "" && value !== 0) {
            search.set(key, String(value));
        }
    }
    const qs = search.toString();
    return qs ? `?${qs}` : "";
}

export function audioUrl(charId: string, audioId: string): string {
    return `${API_BASE}/audio/${charId}/${audioId}`;
}

export function combinedAudioUrl(charId: string, ids: string): string {
    return `${API_BASE}/audio/${charId}/combined?ids=${ids}`;
}

export function ogImageUrl(audioId: string, lang: string): string {
    const langParam = lang && lang !== "en" ? `?lang=${lang}` : "";
    return `${API_BASE}/og/${audioId}.png${langParam}`;
}
