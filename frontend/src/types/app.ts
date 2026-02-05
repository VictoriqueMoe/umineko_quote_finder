export type Language = "en" | "ja";

export type ThemeType = "featherine" | "bernkastel" | "lambdadelta";

export type ViewMode = "search" | "browse" | "stats" | "featured" | "quoteLookup" | "voiceBuilder";

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
