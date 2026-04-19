import { useCallback, useEffect, useRef, useState } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { getContext, getNearestVoiced } from "../../api/endpoints";
import type { ContextResponse } from "../../types/api";
import type { Language } from "../../types/app";

interface ContextViewerProps {
    audioId: string;
    onQuoteClick?: (audioId: string) => void;
    langOverride?: Language;
}

export function ContextViewer({ audioId, onQuoteClick, langOverride }: ContextViewerProps) {
    const { language, game } = useAppContext();
    const [data, setData] = useState<ContextResponse | null>(null);
    const [visible, setVisible] = useState(false);
    const [loading, setLoading] = useState(false);
    const [navigating, setNavigating] = useState(false);
    const [centerAudioId, setCenterAudioId] = useState<string | null>(null);
    const firstId = audioId.split(", ")[0];
    const prevFirstId = useRef(firstId);

    const isNavigated = centerAudioId !== null && centerAudioId !== firstId;
    const effectiveCenterId = centerAudioId || firstId;

    const fetchContext = useCallback(
        async (targetAudioId: string, lang?: Language) => {
            const effectiveLang = lang || langOverride || language;
            setLoading(true);
            try {
                const result = await getContext(game, targetAudioId, effectiveLang, 5);
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
        [langOverride, language, game],
    );

    useEffect(() => {
        if (firstId !== prevFirstId.current) {
            prevFirstId.current = firstId;
            setCenterAudioId(null);
            if (visible) {
                // eslint-disable-next-line react-hooks/set-state-in-effect -- intentional: refetch on prop change
                void fetchContext(firstId);
            }
        }
    }, [firstId, visible, fetchContext]);

    const handleToggle = useCallback(() => {
        if (visible) {
            setVisible(false);
        } else {
            setCenterAudioId(null);
            fetchContext(firstId);
        }
    }, [visible, fetchContext, firstId]);

    const handleNavigate = useCallback(
        async (direction: "next" | "prev") => {
            const effectiveLang = langOverride || language;
            setNavigating(true);
            try {
                const result = await getNearestVoiced(game, effectiveCenterId, effectiveLang, direction);
                if (result.audioId) {
                    setCenterAudioId(result.audioId);
                    await fetchContext(result.audioId);
                }
            } catch {
                // no more voiced quotes in that direction
            } finally {
                setNavigating(false);
            }
        },
        [effectiveCenterId, langOverride, language, fetchContext, game],
    );

    const handleReset = useCallback(() => {
        setCenterAudioId(null);
        fetchContext(firstId);
    }, [firstId, fetchContext]);

    const quoteAudioId = data?.quote?.audioId || "";

    return (
        <>
            <button className="context-btn" disabled={loading} onClick={handleToggle}>
                {loading && !navigating ? "Loading..." : visible ? "Hide Context" : "Show Context"}
            </button>
            {visible && data && (
                <div className="context-section">
                    <div className="context-nav">
                        <button
                            className="context-nav-btn"
                            disabled={navigating}
                            onClick={() => handleNavigate("prev")}
                        >
                            Prev Quote
                        </button>
                        <button
                            className="context-nav-btn context-nav-reset"
                            disabled={navigating || !isNavigated}
                            onClick={handleReset}
                        >
                            Reset
                        </button>
                        <button
                            className="context-nav-btn"
                            disabled={navigating}
                            onClick={() => handleNavigate("next")}
                        >
                            Next Quote
                        </button>
                    </div>
                    {[...data.before, data.quote, ...data.after].map((line, i) => {
                        const isHighlight = line.audioId === quoteAudioId && quoteAudioId !== "";
                        const lineFirstId = line.audioId ? line.audioId.split(", ")[0] : "";
                        const isClickable = lineFirstId && (!isHighlight || isNavigated);

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
