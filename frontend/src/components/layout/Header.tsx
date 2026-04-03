import type { Game, Language, ViewMode } from "../../types/app";
import { ThemeSelector } from "./ThemeSelector";

const HOME_VIEWS = new Set<ViewMode>(["featured", "search", "browse", "quoteLookup"]);

const GAME_CONFIG = {
    umineko: {
        title: "Umineko Quotes",
        subtitle: "When the seagulls cry, none shall remain",
    },
    higurashi: {
        title: "Higurashi Quotes",
        subtitle: "When the cicadas cry, none shall escape",
    },
} as const;

interface HeaderProps {
    language: Language;
    game: Game;
    viewMode: ViewMode;
    onLanguageChange: (lang: Language) => void;
    onGameChange: (game: Game) => void;
    onHomeClick: () => void;
    onStatsClick: () => void;
    onBuilderClick: () => void;
    onBookmarksClick: () => void;
    bookmarkCount: number;
}

export function Header({
    language,
    game,
    viewMode,
    onLanguageChange,
    onGameChange,
    onHomeClick,
    onStatsClick,
    onBuilderClick,
    onBookmarksClick,
    bookmarkCount,
}: HeaderProps) {
    const config = GAME_CONFIG[game];

    return (
        <header className="header">
            <div className="game-toggle">
                <button
                    className={`game-toggle-btn${game === "umineko" ? " active" : ""}`}
                    onClick={() => onGameChange("umineko")}
                >
                    Umineko
                </button>
                <button
                    className={`game-toggle-btn${game === "higurashi" ? " active" : ""}`}
                    onClick={() => onGameChange("higurashi")}
                >
                    Higurashi
                </button>
            </div>
            <div className="ornament">{"\u2726 \u2726 \u2726"}</div>
            <h1 className="title" onClick={onHomeClick} style={{ cursor: "pointer" }}>
                {config.title}
            </h1>
            <p className="subtitle">{config.subtitle}</p>
            <ThemeSelector />
            <div className="language-selector">
                <button
                    className={`lang-btn${language === "auto" ? " active" : ""}`}
                    onClick={() => onLanguageChange("auto")}
                >
                    Auto
                </button>
                <button
                    className={`lang-btn${language === "en" ? " active" : ""}`}
                    onClick={() => onLanguageChange("en")}
                >
                    English
                </button>
                {game === "higurashi" && (
                    <button
                        className={`lang-btn${language === "ja" ? " active" : ""}`}
                        onClick={() => onLanguageChange("ja")}
                    >
                        {"日本語"}
                    </button>
                )}
                {game === "umineko" && (
                    <>
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
                    </>
                )}
            </div>
            <nav className="header-nav">
                <button className={`header-nav-btn${HOME_VIEWS.has(viewMode) ? " active" : ""}`} onClick={onHomeClick}>
                    {"\u2302 Home"}
                </button>
                {game === "umineko" && (
                    <button
                        className={`header-nav-btn${viewMode === "voiceBuilder" ? " active" : ""}`}
                        onClick={onBuilderClick}
                    >
                        {"\u266B Voice Builder"}
                    </button>
                )}
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
