export type Language = "en" | "wh" | "ja" | "es" | "pt";

export type ThemeType = "featherine" | "bernkastel" | "lambdadelta";

export type ViewMode = "search" | "browse" | "stats" | "featured" | "quoteLookup" | "voiceBuilder" | "bookmarks";

export interface FilterState {
    character: string;
    interactionA: string;
    interactionB: string;
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
