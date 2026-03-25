import type { FilterState, Language, ViewMode } from "../types/app";
import { normalizeCharacterKey } from "../utils/characterIds";

interface CommonParams {
    lang: Language;
    episode: string;
    truth: string;
    interactionA: string;
    interactionB: string;
    offset: number;
}

export type ParsedRoute = CommonParams &
    (
        | { viewMode: "search"; query: string; character: string }
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
}

interface RouteRule {
    viewMode: ViewMode;
    match: (params: URLSearchParams) => boolean;
    parse: (params: URLSearchParams) => Record<string, string>;
    serialize: (params: URLSearchParams, ctx: SerializeContext) => void;
}

const ROUTES: RouteRule[] = [
    {
        viewMode: "stats",
        match: p => p.get("stats") === "1",
        parse: () => ({}),
        serialize: p => {
            p.set("stats", "1");
        },
    },
    {
        viewMode: "voiceBuilder",
        match: p => !!p.get("builder"),
        parse: p => ({ segments: p.get("builder")! }),
        serialize: p => {
            p.set("builder", "1");
        },
    },
    {
        viewMode: "quoteLookup",
        match: p => !!p.get("quote"),
        parse: p => ({ audioId: p.get("quote")! }),
        serialize: (p, ctx) => {
            if (ctx.currentAudioId) {
                p.set("quote", ctx.currentAudioId);
            }
        },
    },
    {
        viewMode: "browse",
        match: p => !!p.get("browse"),
        parse: p => {
            const browse = p.get("browse")!;
            return { character: browse !== "1" ? normalizeCharacterKey(browse) : "" };
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
        parse: p => ({
            query: p.get("q")!,
            character: normalizeCharacterKey(p.get("character") || ""),
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
        },
    },
];

export function parseRoute(search: string): ParsedRoute {
    const params = new URLSearchParams(search);

    const lang = (params.get("lang") || "auto") as Language;
    const episode = params.get("episode") || "0";
    const truth = params.get("truth") || "";
    const interactionA = normalizeCharacterKey(params.get("interactionA") || "");
    const interactionB = normalizeCharacterKey(params.get("interactionB") || "");
    const offset = parseInt(params.get("offset") || "0") || 0;
    const common = { lang, episode, truth, interactionA, interactionB, offset };

    for (const route of ROUTES) {
        if (route.match(params)) {
            return { ...common, viewMode: route.viewMode, ...route.parse(params) } as ParsedRoute;
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
    searchQuery: string,
): string {
    const params = new URLSearchParams();

    const route = ROUTES.find(r => r.viewMode === state.viewMode);
    if (route) {
        route.serialize(params, { ...state, searchQuery });
    }

    if (state.filters.episode && state.filters.episode !== "0") {
        params.set("episode", state.filters.episode);
    }
    if (state.filters.truth) {
        params.set("truth", state.filters.truth);
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
