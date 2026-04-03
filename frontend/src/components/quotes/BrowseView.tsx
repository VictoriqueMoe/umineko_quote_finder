import { QuoteCard } from "./QuoteCard";
import { InteractionResults } from "./InteractionResults";
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
        const noResultsMessage =
            filters.interactionA && filters.interactionB
                ? "No interaction lines found for the selected filters."
                : "No dialogue found for this character.";
        return <EmptyState message={noResultsMessage} />;
    }

    const epLabel = filters.episode && filters.episode !== "0" ? ` \u2014 Episode ${filters.episode}` : "";
    const truthLabel =
        filters.truth === "red" ? " \u2014 Red Truth" : filters.truth === "blue" ? " \u2014 Blue Truth" : "";
    const arcLabel = filters.arc ? ` \u2014 ${filters.arc}` : "";
    const titleName = data.character || "All Characters";
    const isInteractionMode = !!filters.interactionA && !!filters.interactionB;

    if (isInteractionMode) {
        return (
            <>
                <InteractionResults
                    mode="browse"
                    quotes={data.quotes}
                    offset={offset}
                    total={total}
                    interactionA={filters.interactionA}
                    interactionB={filters.interactionB}
                    characterFilter={filters.character}
                    episodeFilter={filters.episode}
                    truthFilter={filters.truth}
                    onContextQuoteClick={onContextQuoteClick}
                />
                <Pagination total={total} offset={offset} onPaginate={onPaginate} />
            </>
        );
    }

    return (
        <>
            <div className="browse-header">
                <h2 className="browse-title">
                    {titleName}
                    {epLabel}
                    {truthLabel}
                    {arcLabel}
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
