import { useCallback } from "react";
import { QuoteCard } from "./QuoteCard";
import { EmptyState } from "../common/EmptyState";
import { useBookmarks } from "../../hooks/useBookmarks";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface BookmarksViewProps {
    audioPlayer: AudioPlayer;
    onContextQuoteClick?: (audioId: string) => void;
}

export function BookmarksView({ audioPlayer, onContextQuoteClick }: BookmarksViewProps) {
    const { bookmarks, remove, clear } = useBookmarks();

    const handleClearAll = useCallback(() => {
        if (window.confirm("Remove all bookmarks?")) {
            clear();
        }
    }, [clear]);

    if (bookmarks.length === 0) {
        return <EmptyState message="No bookmarks yet. Bookmark quotes to save them here." />;
    }

    return (
        <>
            <div className="browse-header">
                <h2 className="browse-title">Bookmarks</h2>
                <p className="browse-subtitle">
                    {bookmarks.length} saved {bookmarks.length === 1 ? "fragment" : "fragments"}
                </p>
                <button className="bookmarks-clear-btn" onClick={handleClearAll}>
                    Clear All
                </button>
            </div>
            <div className="quotes-grid">
                {bookmarks.map((bookmark, index) => (
                    <div key={bookmark.quote.audioId ?? index} className="bookmark-entry">
                        <QuoteCard
                            quote={bookmark.quote}
                            index={index}
                            audioPlayer={audioPlayer}
                            onContextQuoteClick={onContextQuoteClick}
                        />
                        <button className="bookmark-remove-btn" onClick={() => remove(bookmark.quote)}>
                            Remove
                        </button>
                    </div>
                ))}
            </div>
        </>
    );
}
