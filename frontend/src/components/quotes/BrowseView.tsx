import { QuoteCard } from "./QuoteCard";
import { Pagination } from "../common/Pagination";
import { EmptyState } from "../common/EmptyState";
import type { BrowseResponse } from "../../types/api";
import type { FilterState } from "../../types/app";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface BrowseViewProps {
    data: BrowseResponse;
    offset: number;
    total: number;
    onPaginate: (newOffset: number) => void;
    audioPlayer: AudioPlayer;
    filters: FilterState;
    onContextQuoteClick?: (audioId: string) => void;
}

export function BrowseView({
    data,
    offset,
    total,
    onPaginate,
    audioPlayer,
    filters,
    onContextQuoteClick,
}: BrowseViewProps) {
    if (!data.quotes || data.quotes.length === 0) {
        return <EmptyState message="No dialogue found for this character." />;
    }

    const epLabel = filters.episode && filters.episode !== "0" ? ` \u2014 Episode ${filters.episode}` : "";
    const truthLabel =
        filters.truth === "red" ? " \u2014 Red Truth" : filters.truth === "blue" ? " \u2014 Blue Truth" : "";
    const titleName = data.character || "All Characters";

    return (
        <>
            <div className="browse-header">
                <h2 className="browse-title">
                    {titleName}
                    {epLabel}
                    {truthLabel}
                </h2>
                <p className="browse-subtitle">
                    Showing lines {data.offset + 1}-{data.offset + data.quotes.length} of {data.total} in story order
                </p>
            </div>
            <div className="quotes-grid">
                {data.quotes.map((quote, index) => {
                    const lineNum = data.offset + index + 1;
                    return (
                        <QuoteCard
                            key={`${quote.audioId || index}-${offset}`}
                            quote={quote}
                            index={index}
                            lineNumber={lineNum}
                            audioPlayer={audioPlayer}
                            onContextQuoteClick={onContextQuoteClick}
                        />
                    );
                })}
            </div>
            <Pagination total={total} offset={offset} onPaginate={onPaginate} />
        </>
    );
}
