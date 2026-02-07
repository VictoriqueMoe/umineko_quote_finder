import { useCallback, useMemo, useState } from "react";
import { combinedAudioUrl, resolveCharId } from "../api/client";
import { getQuoteByAudioId } from "../api/endpoints";
import type { Quote } from "../types/api";
import type { Language } from "../types/app";
import { arrayMove } from "@dnd-kit/sortable";

export interface BuilderSegment {
    id: string;
    audioId: string;
    charId: string;
    characterName: string;
    quoteText: string;
    episode?: number;
}

const MAX_SEGMENTS = 20;
const STORAGE_KEY = "uminekoVoiceBuilder";

function loadFromStorage(): BuilderSegment[] {
    try {
        const stored = sessionStorage.getItem(STORAGE_KEY);
        if (stored) {
            return JSON.parse(stored);
        }
    } catch {
        // ignore corrupt storage
    }
    return [];
}

function saveToStorage(segments: BuilderSegment[]) {
    try {
        sessionStorage.setItem(STORAGE_KEY, JSON.stringify(segments));
    } catch {
        // ignore full storage
    }
}

export function segmentFromQuote(quote: Quote, audioId: string): BuilderSegment {
    const charId = resolveCharId(audioId, quote.characterId ?? "", quote.audioCharMap);
    const clipText = quote.audioTextMap?.[audioId] ?? quote.text;
    return {
        id: crypto.randomUUID(),
        audioId,
        charId,
        characterName: quote.character,
        quoteText: clipText,
        episode: quote.episode,
    };
}

export function useVoiceBuilder() {
    const [segments, setSegments] = useState<BuilderSegment[]>(loadFromStorage);

    const canAdd = segments.length < MAX_SEGMENTS;
    const segmentCount = segments.length;

    const updateSegments = useCallback((updater: (prev: BuilderSegment[]) => BuilderSegment[]) => {
        setSegments(prev => {
            const next = updater(prev);
            saveToStorage(next);
            return next;
        });
    }, []);

    const addSegment = useCallback(
        (segment: BuilderSegment): boolean => {
            let added = false;
            updateSegments(prev => {
                if (prev.length >= MAX_SEGMENTS) {
                    return prev;
                }
                added = true;
                return [...prev, segment];
            });
            return added;
        },
        [updateSegments],
    );

    const removeSegment = useCallback(
        (id: string) => {
            updateSegments(prev => prev.filter(s => s.id !== id));
        },
        [updateSegments],
    );

    const reorderSegments = useCallback(
        (activeId: string, overId: string) => {
            updateSegments(prev => {
                const oldIndex = prev.findIndex(s => s.id === activeId);
                const newIndex = prev.findIndex(s => s.id === overId);
                if (oldIndex === -1 || newIndex === -1) {
                    return prev;
                }
                return arrayMove(prev, oldIndex, newIndex);
            });
        },
        [updateSegments],
    );

    const clearAll = useCallback(() => {
        updateSegments(() => []);
    }, [updateSegments]);

    const combinedUrl = useMemo(() => {
        if (segments.length === 0) {
            return null;
        }
        return combinedAudioUrl(segments.map(s => ({ charId: s.charId, audioId: s.audioId })));
    }, [segments]);

    const shareUrl = useMemo(() => {
        if (segments.length === 0) {
            return "";
        }
        const param = segments.map(s => `${s.charId}:${s.audioId}`).join(",");
        return `${window.location.origin}/?builder=${param}`;
    }, [segments]);

    const loadFromUrl = useCallback(
        async (param: string, language: Language) => {
            const parts = param.split(",").filter(Boolean);
            const newSegments: BuilderSegment[] = [];

            for (const part of parts.slice(0, MAX_SEGMENTS)) {
                const colonIdx = part.indexOf(":");
                if (colonIdx === -1) {
                    continue;
                }
                const charId = part.slice(0, colonIdx);
                const audioId = part.slice(colonIdx + 1);
                if (!charId || !audioId) {
                    continue;
                }

                try {
                    const quote = await getQuoteByAudioId(audioId, language);
                    const clipText = quote.audioTextMap?.[audioId] ?? quote.text;
                    newSegments.push({
                        id: crypto.randomUUID(),
                        audioId,
                        charId,
                        characterName: quote.character,
                        quoteText: clipText,
                        episode: quote.episode,
                    });
                } catch {
                    // If we can't fetch metadata, still add with minimal info
                    newSegments.push({
                        id: crypto.randomUUID(),
                        audioId,
                        charId,
                        characterName: `Character ${charId}`,
                        quoteText: audioId,
                    });
                }
            }

            if (newSegments.length > 0) {
                updateSegments(() => newSegments);
            }
        },
        [updateSegments],
    );

    return {
        segments,
        canAdd,
        segmentCount,
        maxSegments: MAX_SEGMENTS,
        addSegment,
        removeSegment,
        reorderSegments,
        clearAll,
        combinedUrl,
        shareUrl,
        loadFromUrl,
    };
}

export type VoiceBuilder = ReturnType<typeof useVoiceBuilder>;
