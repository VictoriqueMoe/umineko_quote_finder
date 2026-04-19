export interface Quote {
    text: string;
    textHtml?: string;
    textJp?: string;
    textJpHtml?: string;
    arc?: string;
    character: string;
    characterId?: string;
    episode?: number;
    contentType?: string;
    audioId?: string;
    audioCharMap?: Record<string, string>;
    audioTextMap?: Record<string, string>;
    soundEffects?: SoundEffect[];
    chapter?: string;
}

export interface SearchResult {
    quote: Quote;
}

export interface SearchResponse {
    results: SearchResult[];
    total: number;
    offset: number;
    lang?: string;
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

export interface SoundEffect {
    filename: string;
    afterClip: number;
}

export interface NearestVoicedResponse {
    audioId: string;
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
    gold: number;
    purple: number;
}

export interface Interaction {
    charA: string;
    charB: string;
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
    interactionCounts: Record<string, number>;
    characterPresence: CharacterPresence[];
    characterNames: Record<string, string>;
    episodeNames: Record<number, string>;
}

export interface HigurashiStatsResponse {
    topSpeakers: TopSpeaker[];
    linesPerArc: Record<string, Record<string, number>>;
    interactions: Interaction[];
    interactionCounts: Record<string, number>;
    characterNames: Record<string, string>;
}

export interface CiconiaStatsResponse {
    topSpeakers: TopSpeaker[];
    linesPerChapter: Record<string, Record<string, number>>;
    interactions: Interaction[];
    interactionCounts: Record<string, number>;
    characterNames: Record<string, string>;
}

export interface ConfigResponse {
    hasAudio: boolean;
}

export interface CharactersResponse {
    characters: Record<string, string>;
    additional: Record<string, string>;
}
