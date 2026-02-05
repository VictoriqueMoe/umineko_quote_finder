import { useCallback, useRef } from "react";
import { episodeLabel, useQuoteDisplay } from "../../hooks/useQuoteDisplay";
import { AudioPlayer } from "../audio/AudioPlayer";
import { LangToggle } from "./LangToggle";
import { ShareButton } from "./ShareButton";
import { DownloadButton } from "./DownloadButton";
import { ContextViewer } from "./ContextViewer";
import type { Quote } from "../../types/api";
import type { AudioPlayer as AudioPlayerType } from "../../hooks/useAudioPlayer";
import type { Language } from "../../types/app";

interface QuoteCardProps {
    quote: Quote;
    index: number;
    lineNumber?: number;
    audioPlayer: AudioPlayerType;
    onContextQuoteClick?: (audioId: string) => void;
}

export function QuoteCard({ quote, index, lineNumber, audioPlayer, onContextQuoteClick }: QuoteCardProps) {
    const { displayHtml, lang, hasAudio, handleTextUpdate, handleLangChange } = useQuoteDisplay(quote);
    const contextRefreshRef = useRef<((lang: Language) => void) | null>(null);

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
                <AudioPlayer
                    audioId={quote.audioId}
                    characterId={quote.characterId}
                    audioCharMap={quote.audioCharMap}
                    audioPlayer={audioPlayer}
                />
            )}
            {quote.audioId && (
                <div className="quote-actions">
                    <ContextViewer audioId={quote.audioId} onQuoteClick={onContextQuoteClick} langOverride={lang} />
                    <ShareButton audioId={quote.audioId} lang={lang} />
                    <DownloadButton audioId={quote.audioId} lang={lang} />
                </div>
            )}
        </article>
    );
}
