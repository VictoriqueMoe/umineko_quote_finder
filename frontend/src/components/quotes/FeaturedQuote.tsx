import { useCallback, useEffect, useState } from "react";
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

interface FeaturedQuoteProps {
    quote: Quote;
    audioPlayer: AudioPlayerType;
    onContextQuoteClick?: (audioId: string) => void;
}

export function FeaturedQuote({ quote, audioPlayer, onContextQuoteClick }: FeaturedQuoteProps) {
    const { language, hasAudio } = useAppContext();
    const [displayHtml, setDisplayHtml] = useState(quote.textHtml || quote.text);
    const [lang, setLang] = useState<Language>(language);

    useEffect(() => {
        setLang(language);
    }, [language]);

    const handleTextUpdate = useCallback((textHtml: string) => {
        setDisplayHtml(textHtml);
    }, []);

    const handleLangChange = useCallback((newLang: Language) => {
        setLang(newLang);
    }, []);

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
                </div>
            )}
        </article>
    );
}
