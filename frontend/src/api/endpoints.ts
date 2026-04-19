import { apiFetch, buildQueryString, gameApiFetch } from "./client";
import type {
    BrowseResponse,
    CharactersResponse,
    CiconiaStatsResponse,
    ConfigResponse,
    ContextResponse,
    HigurashiStatsResponse,
    NearestVoicedResponse,
    Quote,
    SearchResponse,
    StatsResponse,
} from "../types/api";
import type { Game, Language } from "../types/app";
import { resolveLanguage } from "../types/app";

const PAGE_SIZE = 30;

export { PAGE_SIZE };

export async function searchQuotes(
    game: Game,
    query: string,
    lang: Language,
    offset: number = 0,
    characterId?: string,
    interactionA?: string,
    interactionB?: string,
    episode?: string,
    truth?: string,
    arc?: string,
    exact?: boolean,
    chapter?: string,
): Promise<SearchResponse> {
    const qs = buildQueryString({
        q: query,
        limit: PAGE_SIZE,
        offset,
        lang,
        character: characterId || undefined,
        interactionA: interactionA || undefined,
        interactionB: interactionB || undefined,
        episode: episode && episode !== "0" ? episode : undefined,
        truth: truth || undefined,
        arc: arc || undefined,
        chapter: chapter || undefined,
        exact: exact ? "true" : undefined,
    });
    return gameApiFetch<SearchResponse>(game, `/search${qs}`);
}

export async function getRandomQuote(
    game: Game,
    lang: Language,
    characterId?: string,
    episode?: string,
    truth?: string,
    arc?: string,
    chapter?: string,
): Promise<Quote> {
    const qs = buildQueryString({
        lang: resolveLanguage(lang),
        character: characterId || undefined,
        episode: episode && episode !== "0" ? episode : undefined,
        truth: truth || undefined,
        arc: arc || undefined,
        chapter: chapter || undefined,
    });
    return gameApiFetch<Quote>(game, `/random${qs}`);
}

export async function getQuoteByAudioId(game: Game, audioId: string, lang: Language): Promise<Quote> {
    return gameApiFetch<Quote>(game, `/quote/${audioId}?lang=${resolveLanguage(lang)}`);
}

export async function browseDialogue(
    game: Game,
    lang: Language,
    offset: number = 0,
    characterId?: string,
    interactionA?: string,
    interactionB?: string,
    episode?: string,
    truth?: string,
    arc?: string,
    chapter?: string,
): Promise<BrowseResponse> {
    const qs = buildQueryString({
        limit: PAGE_SIZE,
        offset,
        lang: resolveLanguage(lang),
        character: characterId || undefined,
        interactionA: interactionA || undefined,
        interactionB: interactionB || undefined,
        episode: episode && episode !== "0" ? episode : undefined,
        truth: truth || undefined,
        arc: arc || undefined,
        chapter: chapter || undefined,
    });
    return gameApiFetch<BrowseResponse>(game, `/browse${qs}`);
}

export async function getContext(
    game: Game,
    audioId: string,
    lang: Language,
    lines: number = 5,
): Promise<ContextResponse> {
    return gameApiFetch<ContextResponse>(game, `/context/${audioId}?lang=${resolveLanguage(lang)}&lines=${lines}`);
}

export async function getNearestVoiced(
    game: Game,
    audioId: string,
    lang: Language,
    direction: "next" | "prev",
): Promise<NearestVoicedResponse> {
    return gameApiFetch<NearestVoicedResponse>(
        game,
        `/nearest-voiced/${audioId}?lang=${resolveLanguage(lang)}&direction=${direction}`,
    );
}

export async function getStats(
    game: Game,
    episode?: string,
): Promise<StatsResponse | HigurashiStatsResponse | CiconiaStatsResponse> {
    const qs = episode && episode !== "0" ? `?episode=${episode}` : "";
    return gameApiFetch<StatsResponse | HigurashiStatsResponse | CiconiaStatsResponse>(game, `/stats${qs}`);
}

export async function getConfig(): Promise<ConfigResponse> {
    return apiFetch<ConfigResponse>("/config");
}

export async function getCharacters(game: Game): Promise<CharactersResponse> {
    return gameApiFetch<CharactersResponse>(game, "/characters");
}
