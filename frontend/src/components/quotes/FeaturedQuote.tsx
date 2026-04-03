import { episodeLabel, useQuoteDisplay } from "../../hooks/useQuoteDisplay";
import { useAppContext } from "../../hooks/useAppContext";
import { AudioPlayer } from "../audio/AudioPlayer";
import { SePlayer } from "../audio/SePlayer";
import { LangToggle } from "./LangToggle";
import { ShareButton } from "./ShareButton";
import { DownloadButton } from "./DownloadButton";
import { BookmarkButton } from "./BookmarkButton";
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
    const { game } = useAppContext();

    return (
        <article className="featured-quote">
            <BookmarkButton quote={quote} />
            <div className="featured-label">{"\u2726 A Fragment from the Sea \u2726"}</div>
            <p className="featured-text" dangerouslySetInnerHTML={{ __html: `&ldquo;${displayHtml}&rdquo;` }} />
            <p className="featured-character">&mdash; {quote.character}</p>
            {quote.episode ? <p className="featured-episode">{episodeLabel(quote)}</p> : null}
            {quote.arc && <p className="featured-episode">{quote.arc}</p>}
            {hasAudio && quote.audioId && quote.characterId && (
                <AudioPlayer
                    audioId={quote.audioId}
                    characterId={quote.characterId}
                    audioCharMap={quote.audioCharMap}
                    audioPlayer={audioPlayer}
                />
            )}
            {hasAudio && quote.soundEffects && quote.soundEffects.length > 0 && (
                <SePlayer soundEffects={quote.soundEffects} audioPlayer={audioPlayer} />
            )}
            {quote.audioId && (
                <LangToggle
                    audioId={quote.audioId}
                    onTextUpdate={handleTextUpdate}
                    onLangChange={handleLangChange}
                    textJp={quote.textJp}
                    textJpHtml={quote.textJpHtml}
                    originalText={quote.text}
                    originalTextHtml={quote.textHtml}
                />
            )}
            {quote.audioId && (
                <div className="quote-actions">
                    <ContextViewer audioId={quote.audioId} onQuoteClick={onContextQuoteClick} langOverride={lang} />
                    <ShareButton audioId={quote.audioId} lang={lang} game={game} />
                    <DownloadButton audioId={quote.audioId} lang={lang} />
                </div>
            )}
        </article>
    );
}
