import { useCallback, useState } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { getContext } from "../../api/endpoints";
import type { ContextResponse } from "../../types/api";
import type { Language } from "../../types/app";

interface ContextViewerProps {
    audioId: string;
    onQuoteClick?: (audioId: string) => void;
    langOverride?: Language;
}

export function ContextViewer({ audioId, onQuoteClick, langOverride }: ContextViewerProps) {
    const { language } = useAppContext();
    const [data, setData] = useState<ContextResponse | null>(null);
    const [visible, setVisible] = useState(false);
    const [loading, setLoading] = useState(false);
    const firstId = audioId.split(", ")[0];

    const fetchContext = useCallback(
        async (lang?: Language) => {
            const effectiveLang = lang || langOverride || language;
            setLoading(true);
            try {
                const result = await getContext(firstId, effectiveLang, 5);
                if (!result.error) {
                    setData(result);
                    setVisible(true);
                }
            } catch (err) {
                console.error("Failed to load context:", err);
            } finally {
                setLoading(false);
            }
        },
        [firstId, langOverride, language],
    );

    const handleToggle = useCallback(() => {
        if (visible) {
            setVisible(false);
        } else {
            fetchContext();
        }
    }, [visible, fetchContext]);

    const refreshForLang = useCallback(
        (lang: Language) => {
            if (visible) {
                fetchContext(lang);
            }
        },
        [visible, fetchContext],
    );

    // Expose refresh method via the component
    (ContextViewer as { refreshForLang?: (lang: Language) => void }).refreshForLang = refreshForLang;

    const quoteAudioId = data?.quote?.audioId || "";

    return (
        <>
            <button className="context-btn" disabled={loading} onClick={handleToggle}>
                {loading ? "Loading..." : visible ? "Hide Context" : "Show Context"}
            </button>
            {visible && data && (
                <div className="context-section">
                    {[...data.before, data.quote, ...data.after].map((line, i) => {
                        const isHighlight = line.audioId === quoteAudioId && quoteAudioId !== "";
                        const lineFirstId = line.audioId ? line.audioId.split(", ")[0] : "";
                        const isClickable = lineFirstId && !isHighlight;

                        return (
                            <div
                                key={i}
                                className={`context-line${isHighlight ? " context-highlight" : ""}${isClickable ? " context-clickable" : ""}`}
                                onClick={() => {
                                    if (isClickable && onQuoteClick) {
                                        onQuoteClick(lineFirstId);
                                    }
                                }}
                            >
                                <span className="context-character">{line.character}</span>
                                <span
                                    className="context-text"
                                    dangerouslySetInnerHTML={{ __html: line.textHtml || line.text }}
                                />
                            </div>
                        );
                    })}
                </div>
            )}
        </>
    );
}

export type ContextViewerRefresh = (lang: Language) => void;
