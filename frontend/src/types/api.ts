export interface Quote {
    text: string;
    textHtml?: string;
    character: string;
    characterId?: string;
    episode?: number;
    contentType?: string;
    audioId?: string;
}

export interface SearchResult {
    quote: Quote;
}

export interface SearchResponse {
    results: SearchResult[];
    total: number;
    offset: number;
}

export interface BrowseResponse {
    quotes: Quote[];
    total: number;
    offset: number;
    character?: string;
}

export interface ContextLine {
    text: string;
    textHtml?: string;
    character: string;
    audioId?: string;
}

export interface ContextResponse {
    before: ContextLine[];
    quote: ContextLine;
    after: ContextLine[];
    error?: string;
}

export interface TopSpeaker {
    name: string;
    count: number;
}

export interface LinesPerEpisode {
    episode: number;
    episodeName: string;
    characters: Record<string, number>;
}

export interface TruthPerEpisode {
    episode: number;
    red: number;
    blue: number;
}

export interface Interaction {
    nameA: string;
    nameB: string;
    count: number;
}

export interface CharacterPresence {
    name: string;
    episodes: number[];
}

export interface StatsResponse {
    topSpeakers: TopSpeaker[];
    linesPerEpisode: LinesPerEpisode[];
    truthPerEpisode: TruthPerEpisode[];
    interactions: Interaction[];
    characterPresence: CharacterPresence[];
    characterNames: Record<string, string>;
    episodeNames: Record<number, string>;
}

export interface ConfigResponse {
    hasAudio: boolean;
}

export type CharactersResponse = Record<string, string>;
