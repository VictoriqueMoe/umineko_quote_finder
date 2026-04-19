import type { Game } from "../types/app";

const API_BASE = "/api/v1";

export async function apiFetch<T>(path: string): Promise<T> {
    const response = await fetch(`${API_BASE}${path}`);
    if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
    }
    return response.json();
}

export async function gameApiFetch<T>(game: Game, path: string): Promise<T> {
    return apiFetch<T>(`/${game}${path}`);
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

export function audioUrl(game: Game, charId: string, audioId: string): string {
    return `${API_BASE}/${game}/audio/voice/${charId}/${audioId}`;
}

export function combinedAudioUrl(
    game: Game,
    segments: Array<{ charId: string; audioId: string }>,
    delay?: boolean,
): string {
    const param = segments.map(s => `${s.charId}:${s.audioId}`).join(",");
    const delayParam = delay ? "&delay=true" : "";
    return `${API_BASE}/${game}/audio/voice/combined?segments=${param}${delayParam}`;
}

export function resolveCharId(audioId: string, defaultCharId: string, audioCharMap?: Record<string, string>): string {
    return audioCharMap?.[audioId] ?? defaultCharId;
}

export function seAudioUrl(game: Game, filename: string): string {
    return `${API_BASE}/${game}/audio/se/${filename}`;
}

export function ogImageUrl(game: Game, audioId: string, lang: string, full?: boolean): string {
    const params = new URLSearchParams({ lang: lang || "en" });
    if (full) {
        params.set("full", "true");
    }
    if (game === "higurashi" || game === "ciconia") {
        params.set("audioId", audioId);
        return `${API_BASE}/${game}/og/quote.png?${params.toString()}`;
    }
    return `${API_BASE}/${game}/og/${audioId}.png?${params.toString()}`;
}
