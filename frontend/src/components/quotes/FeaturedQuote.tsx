import { episodeLabel, useQuoteDisplay } from "../../hooks/useQuoteDisplay";
import { AudioPlayer } from "../audio/AudioPlayer";
import { LangToggle } from "./LangToggle";
import { ShareButton } from "./ShareButton";
import { DownloadButton } from "./DownloadButton";
import { ContextViewer } from "./ContextViewer";
import type { Quote } from "../../types/api";
import type { AudioPlayer as AudioPlayerType } from "../../hooks/useAudioPlayer";

interface FeaturedQuoteProps {
    quote: Quote;
    audioPlayer: AudioPlayerType;
    onContextQuoteClick?: (audioId: string) => void;
}

export function FeaturedQuote({ quote, audioPlayer, onContextQuoteClick }: FeaturedQuoteProps) {
    const { displayHtml, lang, hasAudio, handleTextUpdate, handleLangChange } = useQuoteDisplay(quote);

    return (
        <article className="featured-quote">
            <div className="featured-label">{"\u2726 A Fragment from the Sea \u2726"}</div>
            <p className="featured-text" dangerouslySetInnerHTML={{ __html: `&ldquo;${displayHtml}&rdquo;` }} />
            <p className="featured-character">&mdash; {quote.character}</p>
            {quote.episode ? <p className="featured-episode">{episodeLabel(quote)}</p> : null}
            {hasAudio && quote.audioId && quote.characterId && (
                <AudioPlayer audioId={quote.audioId} characterId={quote.characterId} audioPlayer={audioPlayer} />
            )}
            {quote.audioId && (
                <LangToggle audioId={quote.audioId} onTextUpdate={handleTextUpdate} onLangChange={handleLangChange} />
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
