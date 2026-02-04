export type Language = "en" | "ja";

export type ViewMode = "search" | "browse" | "stats" | "featured" | "quoteLookup";

export interface FilterState {
    character: string;
    episode: string;
    truth: string;
}

export interface PushUrlParams {
    viewMode: ViewMode;
    filters: FilterState;
    currentAudioId: string | null;
    browseOffset: number;
    searchOffset: number;
}
