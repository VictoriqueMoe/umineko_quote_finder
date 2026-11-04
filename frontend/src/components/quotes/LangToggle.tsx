import { useCallback, useState } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { getQuoteByAudioId } from "../../api/endpoints";
import type { Language } from "../../types/app";
import { GAME_LANGUAGES } from "../../types/app";

type LangToggleProps = {
    audioId: string;
    onTextUpdate: (textHtml: string, text: string) => void;
    onLangChange?: (lang: Language) => void;
    onContextRefresh?: (lang: Language) => void;
    langOverride?: Exclude<Language, "auto">;
};

export function LangToggle({ audioId, onTextUpdate, onLangChange, onContextRefresh, langOverride }: LangToggleProps) {
    const { language, game } = useAppContext();
    const effective = langOverride ?? (language === "auto" ? "en" : language);
    const [activeLang, setActiveLang] = useState<Language>(effective);
    const [prevEffective, setPrevEffective] = useState(effective);
    const [loading, setLoading] = useState(false);

    if (effective !== prevEffective) {
        setPrevEffective(effective);
        setActiveLang(effective);
    }

    const firstId = audioId.split(", ")[0];

    const handleToggle = useCallback(
        async (newLang: Language) => {
            if (newLang === activeLang || loading) {
                return;
            }

            setLoading(true);
            try {
                const quote = await getQuoteByAudioId(game, firstId, newLang);
                if (!("error" in quote)) {
                    setActiveLang(newLang);
                    onTextUpdate(quote.textHtml || quote.text, quote.text);
                    onLangChange?.(newLang);
                    onContextRefresh?.(newLang);
                }
            } catch (err) {
                console.error("Failed to toggle language:", err);
            } finally {
                setLoading(false);
            }
        },
        [firstId, activeLang, loading, onTextUpdate, onLangChange, onContextRefresh, game],
    );

    return (
        <span className="lang-card-toggle" data-audio-id={firstId}>
            {GAME_LANGUAGES[game].map(lang => (
                <button
                    key={lang}
                    className={`lang-card-btn${activeLang === lang ? " active" : ""}`}
                    disabled={loading}
                    onClick={() => handleToggle(lang)}
                >
                    {lang.toUpperCase()}
                </button>
            ))}
        </span>
    );
}
