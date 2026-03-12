import { useCallback, useSyncExternalStore } from "react";
import type { Quote } from "../types/api";

export interface Bookmark {
    quote: Quote;
    savedAt: number;
}

const STORAGE_KEY = "umineko-bookmarks";
const EMPTY: Bookmark[] = [];

let listeners: (() => void)[] = [];
let cached: Bookmark[] | null = null;

function subscribe(listener: () => void) {
    listeners = [...listeners, listener];
    return () => {
        listeners = listeners.filter(l => l !== listener);
    };
}

function getSnapshot(): Bookmark[] {
    if (cached !== null) {
        return cached;
    }
    try {
        const raw = localStorage.getItem(STORAGE_KEY);
        if (!raw) {
            cached = EMPTY;
            return cached;
        }
        cached = JSON.parse(raw);
        return cached!;
    } catch {
        cached = EMPTY;
        return cached;
    }
}

function saveBookmarks(bookmarks: Bookmark[]) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(bookmarks));
    cached = bookmarks;
    listeners.forEach(l => l());
}

function bookmarkKey(quote: Quote): string {
    return quote.audioId?.split(", ")[0] ?? quote.text;
}

export function useBookmarks() {
    const bookmarks = useSyncExternalStore(subscribe, getSnapshot);

    const isBookmarked = useCallback(
        (quote: Quote) => {
            const key = bookmarkKey(quote);
            return bookmarks.some(b => bookmarkKey(b.quote) === key);
        },
        [bookmarks],
    );

    const toggle = useCallback(
        (quote: Quote) => {
            const key = bookmarkKey(quote);
            const existing = bookmarks.findIndex(b => bookmarkKey(b.quote) === key);
            if (existing >= 0) {
                const updated = [...bookmarks];
                updated.splice(existing, 1);
                saveBookmarks(updated);
            } else {
                saveBookmarks([{ quote, savedAt: Date.now() }, ...bookmarks]);
            }
        },
        [bookmarks],
    );

    const remove = useCallback(
        (quote: Quote) => {
            const key = bookmarkKey(quote);
            saveBookmarks(bookmarks.filter(b => bookmarkKey(b.quote) !== key));
        },
        [bookmarks],
    );

    const clear = useCallback(() => {
        saveBookmarks([]);
    }, []);

    return { bookmarks, isBookmarked, toggle, remove, clear, count: bookmarks.length };
}
