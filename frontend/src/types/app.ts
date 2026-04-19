export type Language = "auto" | "en" | "wh" | "ja" | "ru" | "es" | "pt";

export const REAL_LANGUAGES: Exclude<Language, "auto">[] = ["en", "wh", "ja", "ru", "es", "pt"];

export function resolveLanguage(lang: Language): Exclude<Language, "auto"> {
    if (lang === "auto") {
        return "en";
    }
    return lang;
}

export type Game = "umineko" | "higurashi" | "ciconia";

export type UminekoTheme = "featherine" | "bernkastel" | "lambdadelta";
export type HigurashiTheme = "rika" | "mion" | "satoko";
export type CiconiaTheme = "miyao" | "lingji" | "stanislaw";
export type ThemeType = UminekoTheme | HigurashiTheme | CiconiaTheme;

export type ViewMode = "search" | "browse" | "stats" | "featured" | "quoteLookup" | "voiceBuilder" | "bookmarks";

export interface FilterState {
    character: string;
    interactionA: string;
    interactionB: string;
    episode: string;
    truth: string;
    arc: string;
    chapter: string;
}
