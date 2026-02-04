import { QuoteCard } from "./QuoteCard";
import { Pagination } from "../common/Pagination";
import { EmptyState } from "../common/EmptyState";
import type { SearchResult } from "../../types/api";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface QuoteListProps {
    results: SearchResult[];
    query: string;
    total: number;
    offset: number;
    onPaginate: (newOffset: number) => void;
    audioPlayer: AudioPlayer;
    onContextQuoteClick?: (audioId: string) => void;
}

export function QuoteList({
    results,
    query,
    total,
    offset,
    onPaginate,
    audioPlayer,
    onContextQuoteClick,
}: QuoteListProps) {
    if (!results || results.length === 0) {
        return <EmptyState />;
    }

    const start = offset + 1;
    const end = offset + results.length;

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
                    const quote = item.quote || item;
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
