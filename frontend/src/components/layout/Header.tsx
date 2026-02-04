import type { Language } from "../../types/app";

interface HeaderProps {
    language: Language;
    onLanguageChange: (lang: Language) => void;
    onStatsClick: () => void;
}

export function Header({ language, onLanguageChange, onStatsClick }: HeaderProps) {
    return (
        <header className="header">
            <div className="ornament">{"\u2726 \u2726 \u2726"}</div>
            <h1 className="title">Umineko Quotes</h1>
            <p className="subtitle">When the seagulls cry, none shall remain</p>
            <div className="language-selector">
                <button
                    className={`lang-btn${language === "en" ? " active" : ""}`}
                    onClick={() => onLanguageChange("en")}
                >
                    English
                </button>
                <button
                    className={`lang-btn${language === "ja" ? " active" : ""}`}
                    onClick={() => onLanguageChange("ja")}
                >
                    {"日本語"}
                </button>
            </div>
            <nav className="header-nav">
                <button className="header-nav-btn" onClick={onStatsClick}>
                    {"\u2733 Statistics"}
                </button>
            </nav>
        </header>
    );
}
