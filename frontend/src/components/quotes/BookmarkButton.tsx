import { useCallback } from "react";
import { useBookmarks } from "../../hooks/useBookmarks";
import type { Quote } from "../../types/api";

interface BookmarkButtonProps {
    quote: Quote;
}

export function BookmarkButton({ quote }: BookmarkButtonProps) {
    const { isBookmarked, toggle } = useBookmarks();
    const saved = isBookmarked(quote);

    const handleClick = useCallback(() => {
        toggle(quote);
    }, [quote, toggle]);

    return (
        <button
            className={`bookmark-ribbon${saved ? " bookmarked" : ""}`}
            onClick={handleClick}
            title={saved ? "Remove bookmark" : "Bookmark this quote"}
            aria-label={saved ? "Remove bookmark" : "Bookmark this quote"}
        >
            <svg viewBox="0 0 24 32" width="16" height="22">
                <path d="M2 0h20v32l-10-7-10 7z" fill="currentColor" />
            </svg>
        </button>
    );
}
