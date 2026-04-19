import type { FilterState, Game, Language, ViewMode } from "../types/app";
import { normalizeCharacterKey } from "../utils/characterIds";

interface CommonParams {
    lang: Language;
    game: Game;
    episode: string;
    truth: string;
    arc: string;
    chapter: string;
    interactionA: string;
    interactionB: string;
    offset: number;
}

export type ParsedRoute = CommonParams &
    (
        | { viewMode: "search"; query: string; character: string; exact: boolean }
        | { viewMode: "browse"; character: string }
        | { viewMode: "stats" }
        | { viewMode: "quoteLookup"; audioId: string }
        | { viewMode: "voiceBuilder"; segments: string }
        | { viewMode: "featured" }
    );

interface SerializeContext {
    viewMode: ViewMode;
    filters: FilterState;
    currentAudioId: string | null;
    searchOffset: number;
    browseOffset: number;
    searchQuery: string;
    searchExact: boolean;
}

interface RouteRule {
    viewMode: ViewMode;
    match: (params: URLSearchParams) => boolean;
    parse: (params: URLSearchParams, game: Game) => Record<string, string | boolean>;
    serialize: (params: URLSearchParams, ctx: SerializeContext) => void;
}

const ROUTES: RouteRule[] = [
    {
        viewMode: "stats",
        match: p => p.get("stats") === "1",
        parse: (_p, _g) => ({}),
        serialize: p => {
            p.set("stats", "1");
        },
    },
    {
        viewMode: "voiceBuilder",
        match: p => !!p.get("builder"),
        parse: (p, _g) => ({ segments: p.get("builder")! }),
        serialize: p => {
            p.set("builder", "1");
        },
    },
    {
        viewMode: "quoteLookup",
        match: p => !!p.get("quote"),
        parse: (p, _g) => ({ audioId: p.get("quote")! }),
        serialize: (p, ctx) => {
            if (ctx.currentAudioId) {
                p.set("quote", ctx.currentAudioId);
            }
        },
    },
    {
        viewMode: "browse",
        match: p => !!p.get("browse"),
        parse: (p, g) => {
            const browse = p.get("browse")!;
            return { character: browse !== "1" ? normalizeCharacterKey(browse, g) : "" };
        },
        serialize: (p, ctx) => {
            p.set("browse", ctx.filters.character || "1");
            if (ctx.filters.interactionA && ctx.filters.interactionB) {
                p.set("interactionA", ctx.filters.interactionA);
                p.set("interactionB", ctx.filters.interactionB);
            }
        },
    },
    {
        viewMode: "search",
        match: p => !!p.get("q"),
        parse: (p, g) => ({
            query: p.get("q")!,
            character: normalizeCharacterKey(p.get("character") || "", g),
            exact: p.get("exact") === "1",
        }),
        serialize: (p, ctx) => {
            if (ctx.searchQuery.trim()) {
                p.set("q", ctx.searchQuery.trim());
            }
            if (ctx.filters.character) {
                p.set("character", ctx.filters.character);
            }
            if (ctx.filters.interactionA && ctx.filters.interactionB) {
                p.set("interactionA", ctx.filters.interactionA);
                p.set("interactionB", ctx.filters.interactionB);
            }
            if (ctx.searchExact) {
                p.set("exact", "1");
            }
        },
    },
];

export function parseRoute(search: string): ParsedRoute {
    const params = new URLSearchParams(search);

    const lang = (params.get("lang") || "auto") as Language;
    const gameParam = params.get("game");
    const game: Game = gameParam === "higurashi" ? "higurashi" : gameParam === "ciconia" ? "ciconia" : "umineko";
    const episode = params.get("episode") || "0";
    const truth = params.get("truth") || "";
    const arc = params.get("arc") || "";
    const chapter = params.get("chapter") || "";
    const interactionA = normalizeCharacterKey(params.get("interactionA") || "", game);
    const interactionB = normalizeCharacterKey(params.get("interactionB") || "", game);
    const offset = parseInt(params.get("offset") || "0") || 0;
    const common = { lang, game, episode, truth, arc, chapter, interactionA, interactionB, offset };

    for (const route of ROUTES) {
        if (route.match(params)) {
            return { ...common, viewMode: route.viewMode, ...route.parse(params, game) } as ParsedRoute;
        }
    }

    return { ...common, viewMode: "featured" };
}

export function buildUrl(
    state: {
        viewMode: ViewMode;
        filters: FilterState;
        currentAudioId: string | null;
        searchOffset: number;
        browseOffset: number;
    },
    language: Language,
    game: Game,
    searchQuery: string,
    searchExact: boolean,
): string {
    const params = new URLSearchParams();

    if (game !== "umineko") {
        params.set("game", game);
    }

    const route = ROUTES.find(r => r.viewMode === state.viewMode);
    if (route) {
        route.serialize(params, { ...state, searchQuery, searchExact });
    }

    if (state.filters.episode && state.filters.episode !== "0") {
        params.set("episode", state.filters.episode);
    }
    if (state.filters.truth) {
        params.set("truth", state.filters.truth);
    }
    if (state.filters.arc) {
        params.set("arc", state.filters.arc);
    }
    if (state.filters.chapter) {
        params.set("chapter", state.filters.chapter);
    }

    const offset = state.viewMode === "browse" ? state.browseOffset : state.searchOffset;
    if (offset > 0) {
        params.set("offset", String(offset));
    }
    if (language !== "auto") {
        params.set("lang", language);
    }

    const qs = params.toString();
    return qs ? `?${qs}` : window.location.pathname;
}
