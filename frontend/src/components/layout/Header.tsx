import type { Language, ViewMode } from "../../types/app";
import { ThemeSelector } from "./ThemeSelector";

const HOME_VIEWS = new Set<ViewMode>(["featured", "search", "browse", "quoteLookup"]);

interface HeaderProps {
    language: Language;
    viewMode: ViewMode;
    onLanguageChange: (lang: Language) => void;
    onHomeClick: () => void;
    onStatsClick: () => void;
    onBuilderClick: () => void;
    onBookmarksClick: () => void;
    bookmarkCount: number;
}

export function Header({
    language,
    viewMode,
    onLanguageChange,
    onHomeClick,
    onStatsClick,
    onBuilderClick,
    onBookmarksClick,
    bookmarkCount,
}: HeaderProps) {
    return (
        <header className="header">
            <div className="ornament">{"\u2726 \u2726 \u2726"}</div>
            <h1 className="title" onClick={onHomeClick} style={{ cursor: "pointer" }}>
                Umineko Quotes
            </h1>
            <p className="subtitle">When the seagulls cry, none shall remain</p>
            <ThemeSelector />
            <div className="language-selector">
                <button
                    className={`lang-btn${language === "en" ? " active" : ""}`}
                    onClick={() => onLanguageChange("en")}
                >
                    English
                </button>
                <button
                    className={`lang-btn${language === "wh" ? " active" : ""}`}
                    onClick={() => onLanguageChange("wh")}
                >
                    English (WH)
                </button>
                <button
                    className={`lang-btn${language === "ja" ? " active" : ""}`}
                    onClick={() => onLanguageChange("ja")}
                >
                    {"日本語"}
                </button>
                <button
                    className={`lang-btn${language === "ru" ? " active" : ""}`}
                    onClick={() => onLanguageChange("ru")}
                >
                    {"Русский"}
                </button>
                <button
                    className={`lang-btn${language === "es" ? " active" : ""}`}
                    onClick={() => onLanguageChange("es")}
                >
                    Español
                </button>
                <button
                    className={`lang-btn${language === "pt" ? " active" : ""}`}
                    onClick={() => onLanguageChange("pt")}
                >
                    Português
                </button>
            </div>
            <nav className="header-nav">
                <button className={`header-nav-btn${HOME_VIEWS.has(viewMode) ? " active" : ""}`} onClick={onHomeClick}>
                    {"\u2302 Home"}
                </button>
                <button
                    className={`header-nav-btn${viewMode === "voiceBuilder" ? " active" : ""}`}
                    onClick={onBuilderClick}
                >
                    {"\u266B Voice Builder"}
                </button>
                <button className={`header-nav-btn${viewMode === "stats" ? " active" : ""}`} onClick={onStatsClick}>
                    {"\u2733 Statistics"}
                </button>
                <button
                    className={`header-nav-btn${viewMode === "bookmarks" ? " active" : ""}`}
                    onClick={onBookmarksClick}
                >
                    {"\u2605 Bookmarks"}
                    {bookmarkCount > 0 && <span className="bookmark-count">{bookmarkCount}</span>}
                </button>
            </nav>
        </header>
    );
}
