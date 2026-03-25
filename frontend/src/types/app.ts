export type Language = "auto" | "en" | "wh" | "ja" | "ru" | "es" | "pt";

export const REAL_LANGUAGES: Exclude<Language, "auto">[] = ["en", "wh", "ja", "ru", "es", "pt"];

export function resolveLanguage(lang: Language): Exclude<Language, "auto"> {
    if (lang === "auto") {
        return "en";
    }
    return lang;
}

export type ThemeType = "featherine" | "bernkastel" | "lambdadelta";

export type ViewMode = "search" | "browse" | "stats" | "featured" | "quoteLookup" | "voiceBuilder" | "bookmarks";

export interface FilterState {
    character: string;
    interactionA: string;
    interactionB: string;
    episode: string;
    truth: string;
}

