import { useCallback, useEffect, useState } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { getQuoteByAudioId } from "../../api/endpoints";
import type { Language } from "../../types/app";

interface LangToggleProps {
    audioId: string;
    onTextUpdate: (textHtml: string, text: string) => void;
    onLangChange?: (lang: Language) => void;
    onContextRefresh?: (lang: Language) => void;
    langOverride?: Exclude<Language, "auto">;
    textJp?: string;
    textJpHtml?: string;
    originalText?: string;
    originalTextHtml?: string;
}

export function LangToggle({
    audioId,
    onTextUpdate,
    onLangChange,
    onContextRefresh,
    langOverride,
    textJp,
    textJpHtml,
    originalText,
    originalTextHtml,
}: LangToggleProps) {
    const { language, game } = useAppContext();
    const effective = langOverride ?? (language === "auto" ? "en" : language);
    const [activeLang, setActiveLang] = useState<Language>(effective);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const next = langOverride ?? (language === "auto" ? "en" : language);
        setActiveLang(next);
    }, [language, langOverride]);
    const firstId = audioId.split(", ")[0];

    const handleToggle = useCallback(
        async (newLang: Language) => {
            if (newLang === activeLang || loading) {
                return;
            }

            if (game === "higurashi") {
                if (newLang === "ja" && textJp) {
                    setActiveLang("ja");
                    onTextUpdate(textJpHtml || textJp, textJp);
                    onLangChange?.("ja");
                } else if (newLang === "en" && originalText) {
                    setActiveLang("en");
                    onTextUpdate(originalTextHtml || originalText, originalText);
                    onLangChange?.("en");
                }
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
        [
            firstId,
            activeLang,
            loading,
            onTextUpdate,
            onLangChange,
            onContextRefresh,
            game,
            textJp,
            textJpHtml,
            originalText,
            originalTextHtml,
        ],
    );

    if (game === "higurashi") {
        if (!textJp) {
            return null;
        }
        return (
            <span className="lang-card-toggle" data-audio-id={firstId}>
                <button
                    className={`lang-card-btn${activeLang === "en" ? " active" : ""}`}
                    onClick={() => handleToggle("en")}
                >
                    EN
                </button>
                <button
                    className={`lang-card-btn${activeLang === "ja" ? " active" : ""}`}
                    onClick={() => handleToggle("ja")}
                >
                    JA
                </button>
            </span>
        );
    }

    return (
        <span className="lang-card-toggle" data-audio-id={firstId}>
            <button
                className={`lang-card-btn${activeLang === "en" ? " active" : ""}`}
                disabled={loading}
                onClick={() => handleToggle("en")}
            >
                EN
            </button>
            <button
                className={`lang-card-btn${activeLang === "wh" ? " active" : ""}`}
                disabled={loading}
                onClick={() => handleToggle("wh")}
            >
                WH
            </button>
            <button
                className={`lang-card-btn${activeLang === "ja" ? " active" : ""}`}
                disabled={loading}
                onClick={() => handleToggle("ja")}
            >
                JA
            </button>
            <button
                className={`lang-card-btn${activeLang === "ru" ? " active" : ""}`}
                disabled={loading}
                onClick={() => handleToggle("ru")}
            >
                RU
            </button>
            <button
                className={`lang-card-btn${activeLang === "es" ? " active" : ""}`}
                disabled={loading}
                onClick={() => handleToggle("es")}
            >
                ES
            </button>
            <button
                className={`lang-card-btn${activeLang === "pt" ? " active" : ""}`}
                disabled={loading}
                onClick={() => handleToggle("pt")}
            >
                PT
            </button>
        </span>
    );
}
