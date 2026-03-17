import { QuoteCard } from "./QuoteCard";
import { InteractionResults } from "./InteractionResults";
import { Pagination } from "../common/Pagination";
import { EmptyState } from "../common/EmptyState";
import type { SearchResult } from "../../types/api";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";
import type { FilterState } from "../../types/app";

interface QuoteListProps {
    results: SearchResult[];
    query: string;
    total: number;
    offset: number;
    onPaginate: (newOffset: number) => void;
    audioPlayer: AudioPlayer;
    filters: FilterState;
    onContextQuoteClick?: (audioId: string) => void;
}

export function QuoteList({
    results,
    query,
    total,
    offset,
    onPaginate,
    audioPlayer,
    filters,
    onContextQuoteClick,
}: QuoteListProps) {
    const isInteractionMode = !!filters.interactionA && !!filters.interactionB;

    if (!results || results.length === 0) {
        return isInteractionMode ? (
            <EmptyState message="No interaction matches found for this query and filter set." />
        ) : (
            <EmptyState />
        );
    }

    const start = offset + 1;
    const end = offset + results.length;

    if (isInteractionMode) {
        return (
            <>
                <InteractionResults
                    mode="search"
                    quotes={results.map(item => item.quote)}
                    offset={offset}
                    total={total}
                    interactionA={filters.interactionA}
                    interactionB={filters.interactionB}
                    characterFilter={filters.character}
                    episodeFilter={filters.episode}
                    truthFilter={filters.truth}
                    query={query}
                    onContextQuoteClick={onContextQuoteClick}
                />
                <Pagination total={total} offset={offset} onPaginate={onPaginate} />
            </>
        );
    }

    return (
        <>
            {query && (
                <div className="results-header">
                    <span className="results-count">
                        Showing{" "}
                        <span>
                            {start}-{end}
                        </span>{" "}
                        of <span>{total}</span> fragments for &ldquo;{query}&rdquo;
                    </span>
                </div>
            )}
            <div className="quotes-grid">
                {results.map((item, index) => {
                    const quote = item.quote;
                    return (
                        <QuoteCard
                            key={`${quote.audioId || index}-${offset}`}
                            quote={quote}
                            index={index}
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
