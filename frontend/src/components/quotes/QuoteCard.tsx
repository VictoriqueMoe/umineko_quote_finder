import { useCallback, useEffect, useRef, useState } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { AudioPlayer } from "../audio/AudioPlayer";
import { LangToggle } from "./LangToggle";
import { ShareButton } from "./ShareButton";
import { ContextViewer } from "./ContextViewer";
import type { Quote } from "../../types/api";
import type { AudioPlayer as AudioPlayerType } from "../../hooks/useAudioPlayer";
import type { Language } from "../../types/app";

const CONTENT_TYPE_LABELS: Record<string, string> = {
    tea: "Tea Party",
    ura: "????",
    omake: "Omake",
};

function episodeLabel(quote: Quote): string {
    if (!quote.episode) {
        return "";
    }
    let label = `Episode ${quote.episode}`;
    if (quote.contentType && CONTENT_TYPE_LABELS[quote.contentType]) {
        label += ` \u2014 ${CONTENT_TYPE_LABELS[quote.contentType]}`;
    }
    return label;
}

interface QuoteCardProps {
    quote: Quote;
    index: number;
    lineNumber?: number;
    audioPlayer: AudioPlayerType;
    onContextQuoteClick?: (audioId: string) => void;
}

export function QuoteCard({ quote, index, lineNumber, audioPlayer, onContextQuoteClick }: QuoteCardProps) {
    const { language, hasAudio } = useAppContext();
    const [displayHtml, setDisplayHtml] = useState(quote.textHtml || quote.text);
    const [lang, setLang] = useState<Language>(language);

    useEffect(() => {
        setLang(language);
    }, [language]);
    const contextRefreshRef = useRef<((lang: Language) => void) | null>(null);

    const handleTextUpdate = useCallback((textHtml: string) => {
        setDisplayHtml(textHtml);
    }, []);

    const handleLangChange = useCallback((newLang: Language) => {
        setLang(newLang);
    }, []);

    const handleContextRefresh = useCallback((lang: Language) => {
        contextRefreshRef.current?.(lang);
    }, []);

    return (
        <article className="quote-card" style={{ "--index": index } as React.CSSProperties}>
            {lineNumber !== undefined && <span className="quote-number">#{lineNumber}</span>}
            <span className="quote-mark">&ldquo;</span>
            <p className="quote-text" dangerouslySetInnerHTML={{ __html: displayHtml }} />
            <div className="quote-meta">
                <span className="quote-character">&mdash; {quote.character}</span>
                <div className="quote-details">
                    {quote.episode ? <span className="quote-episode">{episodeLabel(quote)}</span> : null}
                    {quote.audioId && (
                        <LangToggle
                            audioId={quote.audioId}
                            onTextUpdate={handleTextUpdate}
                            onLangChange={handleLangChange}
                            onContextRefresh={handleContextRefresh}
                        />
                    )}
                </div>
            </div>
            {hasAudio && quote.audioId && quote.characterId && (
                <AudioPlayer audioId={quote.audioId} characterId={quote.characterId} audioPlayer={audioPlayer} />
            )}
            {quote.audioId && (
                <div className="quote-actions">
                    <ContextViewer audioId={quote.audioId} onQuoteClick={onContextQuoteClick} langOverride={lang} />
                    <ShareButton audioId={quote.audioId} lang={lang} />
                </div>
            )}
        </article>
    );
}
