import { apiFetch, buildQueryString } from "./client";
import type {
    BrowseResponse,
    CharactersResponse,
    ConfigResponse,
    ContextResponse,
    Quote,
    SearchResponse,
    StatsResponse,
} from "../types/api";
import type { Language } from "../types/app";

const PAGE_SIZE = 30;

export { PAGE_SIZE };

export async function searchQuotes(
    query: string,
    lang: Language,
    offset: number = 0,
    characterId?: string,
    episode?: string,
    truth?: string,
): Promise<SearchResponse> {
    const qs = buildQueryString({
        q: query,
        limit: PAGE_SIZE,
        offset,
        lang,
        character: characterId || undefined,
        episode: episode && episode !== "0" ? episode : undefined,
        truth: truth || undefined,
    });
    return apiFetch<SearchResponse>(`/search${qs}`);
}

export async function getRandomQuote(
    lang: Language,
    characterId?: string,
    episode?: string,
    truth?: string,
): Promise<Quote> {
    const qs = buildQueryString({
        lang,
        character: characterId || undefined,
        episode: episode && episode !== "0" ? episode : undefined,
        truth: truth || undefined,
    });
    return apiFetch<Quote>(`/random${qs}`);
}

export async function getQuoteByAudioId(audioId: string, lang: Language): Promise<Quote> {
    return apiFetch<Quote>(`/quote/${audioId}?lang=${lang}`);
}

export async function browseDialogue(
    lang: Language,
    offset: number = 0,
    characterId?: string,
    episode?: string,
    truth?: string,
): Promise<BrowseResponse> {
    const qs = buildQueryString({
        limit: PAGE_SIZE,
        offset,
        lang,
        character: characterId || undefined,
        episode: episode && episode !== "0" ? episode : undefined,
        truth: truth || undefined,
    });
    return apiFetch<BrowseResponse>(`/browse${qs}`);
}

export async function getContext(audioId: string, lang: Language, lines: number = 5): Promise<ContextResponse> {
    return apiFetch<ContextResponse>(`/context/${audioId}?lang=${lang}&lines=${lines}`);
}

export async function getStats(episode?: string): Promise<StatsResponse> {
    const qs = episode && episode !== "0" ? `?episode=${episode}` : "";
    return apiFetch<StatsResponse>(`/stats${qs}`);
}

export async function getConfig(): Promise<ConfigResponse> {
    return apiFetch<ConfigResponse>("/config");
}

export async function getCharacters(): Promise<CharactersResponse> {
    return apiFetch<CharactersResponse>("/characters");
}
