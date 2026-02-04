import { useCallback, useState } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { getQuoteByAudioId } from "../../api/endpoints";
import type { Language } from "../../types/app";

interface LangToggleProps {
    audioId: string;
    onTextUpdate: (textHtml: string, text: string) => void;
    onLangChange?: (lang: Language) => void;
    onContextRefresh?: (lang: Language) => void;
}

export function LangToggle({ audioId, onTextUpdate, onLangChange, onContextRefresh }: LangToggleProps) {
    const { language } = useAppContext();
    const [activeLang, setActiveLang] = useState<Language>(language);
    const [loading, setLoading] = useState(false);
    const firstId = audioId.split(", ")[0];

    const handleToggle = useCallback(
        async (newLang: Language) => {
            if (newLang === activeLang || loading) {
                return;
            }
            setLoading(true);
            try {
                const quote = await getQuoteByAudioId(firstId, newLang);
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
        [firstId, activeLang, loading, onTextUpdate, onLangChange, onContextRefresh],
    );

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
                className={`lang-card-btn${activeLang === "ja" ? " active" : ""}`}
                disabled={loading}
                onClick={() => handleToggle("ja")}
            >
                JA
            </button>
        </span>
    );
}
