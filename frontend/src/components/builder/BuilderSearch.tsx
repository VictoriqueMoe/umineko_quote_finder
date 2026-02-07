import { type KeyboardEvent, useCallback, useState } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { useSearch } from "../../hooks/useSearch";
import { audioUrl, resolveCharId } from "../../api/client";
import { getRandomQuote, PAGE_SIZE } from "../../api/endpoints";
import { episodeLabel } from "../../hooks/useQuoteDisplay";
import { segmentFromQuote, type VoiceBuilder } from "../../hooks/useVoiceBuilder";
import type { Quote } from "../../types/api";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";
import type { FilterState } from "../../types/app";

interface BuilderSearchProps {
    builder: VoiceBuilder;
    audioPlayer: AudioPlayer;
}

export function BuilderSearch({ builder, audioPlayer }: BuilderSearchProps) {
    const { language, sortedCharacters } = useAppContext();
    const search = useSearch();
    const [query, setQuery] = useState("");
    const [collapsed, setCollapsed] = useState(false);
    const [filters, setFilters] = useState<FilterState>({ character: "", episode: "0", truth: "" });

    const handleSearch = useCallback(async () => {
        if (!query.trim()) {
            return;
        }
        await search.search(query, language, 0, filters);
    }, [query, language, filters, search]);

    const handleKeyDown = (e: KeyboardEvent) => {
        if (e.key === "Enter") {
            handleSearch();
        }
    };

    const handlePaginate = useCallback(
        async (newOffset: number) => {
            await search.search(search.query, language, newOffset, filters);
        },
        [search, language, filters],
    );

    const handleRandom = useCallback(async () => {
        try {
            const quote = await getRandomQuote(language, filters.character, filters.episode, filters.truth);
            if (quote && quote.audioId) {
                const ids = quote.audioId.split(", ");
                const firstId = ids[0];
                const segment = segmentFromQuote(quote, firstId);
                builder.addSegment(segment);
            }
        } catch {
            // ignore errors
        }
    }, [language, filters, builder]);

    const handleAddQuote = useCallback(
        (quote: Quote, audioId: string) => {
            const segment = segmentFromQuote(quote, audioId);
            builder.addSegment(segment);
        },
        [builder],
    );

    const handleAddAll = useCallback(
        (quote: Quote) => {
            if (!quote.audioId) {
                return;
            }
            const ids = quote.audioId.split(", ");
            for (const id of ids) {
                if (!builder.canAdd) {
                    break;
                }
                const segment = segmentFromQuote(quote, id);
                builder.addSegment(segment);
            }
        },
        [builder],
    );

    const handlePreviewClip = useCallback(
        (quote: Quote, audioIdSingle: string) => {
            const charId = resolveCharId(audioIdSingle, quote.characterId ?? "", quote.audioCharMap);
            const url = audioUrl(charId, audioIdSingle);
            audioPlayer.play(url, `preview-${audioIdSingle}`);
        },
        [audioPlayer],
    );

    const audioQuotes = (search.results || []).map(r => r.quote).filter(q => !!q.audioId);

    const total = search.total;
    const offset = search.offset;
    const totalPages = Math.ceil(total / PAGE_SIZE);
    const currentPage = Math.floor(offset / PAGE_SIZE) + 1;
    const hasPrev = offset > 0;
    const hasNext = offset + PAGE_SIZE < total;

    return (
        <div className={`builder-search${collapsed ? " collapsed" : ""}`}>
            <button className="builder-search-toggle" onClick={() => setCollapsed(!collapsed)}>
                {collapsed ? "\u25BE Show Search" : "\u25B4 Hide Search"}
            </button>

            {!collapsed && (
                <>
                    <div className="builder-search-bar">
                        <input
                            type="text"
                            className="builder-search-input"
                            placeholder="Search for voice clips..."
                            value={query}
                            onChange={e => setQuery(e.target.value)}
                            onKeyDown={handleKeyDown}
                        />
                        <button className="builder-search-btn" onClick={handleSearch}>
                            Search
                        </button>
                        <button className="builder-random-btn" onClick={handleRandom} disabled={!builder.canAdd}>
                            Random
                        </button>
                    </div>

                    <div className="builder-filters">
                        <select
                            className="builder-filter-select"
                            value={filters.character}
                            onChange={e => setFilters(prev => ({ ...prev, character: e.target.value }))}
                        >
                            <option value="">All Characters</option>
                            {sortedCharacters.map(([id, name]) => (
                                <option key={id} value={id}>
                                    {name}
                                </option>
                            ))}
                        </select>
                        <select
                            className="builder-filter-select"
                            value={filters.episode}
                            onChange={e => setFilters(prev => ({ ...prev, episode: e.target.value }))}
                        >
                            <option value="0">All Episodes</option>
                            <option value="1">{"Ep 1 \u2014 Legend"}</option>
                            <option value="2">{"Ep 2 \u2014 Turn"}</option>
                            <option value="3">{"Ep 3 \u2014 Banquet"}</option>
                            <option value="4">{"Ep 4 \u2014 Alliance"}</option>
                            <option value="5">{"Ep 5 \u2014 End"}</option>
                            <option value="6">{"Ep 6 \u2014 Dawn"}</option>
                            <option value="7">{"Ep 7 \u2014 Requiem"}</option>
                            <option value="8">{"Ep 8 \u2014 Twilight"}</option>
                        </select>
                    </div>

                    {search.loading && <div className="builder-search-loading">Searching...</div>}

                    {!search.loading && search.error && <div className="builder-search-error">{search.error}</div>}

                    {!search.loading && !search.error && audioQuotes.length > 0 && (
                        <>
                            <div className="builder-search-info">
                                Showing {offset + 1}&ndash;{Math.min(offset + PAGE_SIZE, total)} of {total} results
                                (audio only) &mdash;{" "}
                                <span className="builder-search-hint">click on text to add a clip</span>
                            </div>
                            <div className="builder-results">
                                {audioQuotes.map((quote, i) => {
                                    const ids = quote.audioId!.split(", ");
                                    const hasMultiple = ids.length > 1;

                                    return (
                                        <div key={`${quote.audioId}-${i}`} className="builder-result">
                                            <div className="builder-result-text">
                                                <span className="builder-result-quote">
                                                    {hasMultiple && quote.audioTextMap ? (
                                                        ids.map((id, j) => {
                                                            const fragment = quote.audioTextMap?.[id] ?? id;
                                                            const previewing =
                                                                audioPlayer.state.activeId === `preview-${id}` &&
                                                                audioPlayer.state.isPlaying;
                                                            return (
                                                                <span key={id} className="builder-clip-wrap">
                                                                    <span
                                                                        className="builder-clip-preview"
                                                                        onClick={e => {
                                                                            e.stopPropagation();
                                                                            handlePreviewClip(quote, id);
                                                                        }}
                                                                    >
                                                                        {previewing ? "\u275A\u275A" : "\u25B6"}
                                                                    </span>
                                                                    <span
                                                                        className={`builder-clip-text${!builder.canAdd ? " disabled" : ""}`}
                                                                        onClick={() => {
                                                                            if (builder.canAdd) {
                                                                                handleAddQuote(quote, id);
                                                                            }
                                                                        }}
                                                                    >
                                                                        {j === 0 ? "\u201C" : ""}
                                                                        {fragment}
                                                                        {j === ids.length - 1 ? "\u201D" : ""}
                                                                    </span>
                                                                </span>
                                                            );
                                                        })
                                                    ) : (
                                                        <span className="builder-clip-wrap">
                                                            <span
                                                                className="builder-clip-preview"
                                                                onClick={e => {
                                                                    e.stopPropagation();
                                                                    handlePreviewClip(quote, ids[0]);
                                                                }}
                                                            >
                                                                {audioPlayer.state.activeId === `preview-${ids[0]}` &&
                                                                audioPlayer.state.isPlaying
                                                                    ? "\u275A\u275A"
                                                                    : "\u25B6"}
                                                            </span>
                                                            <span
                                                                className={`builder-clip-text${!builder.canAdd ? " disabled" : ""}`}
                                                                onClick={() => {
                                                                    if (builder.canAdd) {
                                                                        handleAddQuote(quote, ids[0]);
                                                                    }
                                                                }}
                                                            >
                                                                {"\u201C"}
                                                                {quote.text}
                                                                {"\u201D"}
                                                            </span>
                                                        </span>
                                                    )}
                                                </span>
                                                <span className="builder-result-meta">
                                                    {"\u2014 "}
                                                    {quote.character}
                                                    {quote.episode ? ` \u00B7 ${episodeLabel(quote)}` : ""}
                                                </span>
                                            </div>
                                            {hasMultiple && (
                                                <div className="builder-result-actions">
                                                    <button
                                                        className="builder-result-btn builder-add-all-btn"
                                                        onClick={() => handleAddAll(quote)}
                                                        disabled={!builder.canAdd}
                                                        title={`Add all ${ids.length} clips`}
                                                    >
                                                        {`+ All (${ids.length})`}
                                                    </button>
                                                </div>
                                            )}
                                        </div>
                                    );
                                })}
                            </div>

                            {total > PAGE_SIZE && (
                                <div className="builder-pagination">
                                    <button
                                        className="pagination-btn"
                                        disabled={!hasPrev}
                                        onClick={() => handlePaginate(offset - PAGE_SIZE)}
                                    >
                                        {"\u25C0 Prev"}
                                    </button>
                                    <span className="pagination-info">
                                        {currentPage} / {totalPages}
                                    </span>
                                    <button
                                        className="pagination-btn"
                                        disabled={!hasNext}
                                        onClick={() => handlePaginate(offset + PAGE_SIZE)}
                                    >
                                        {"Next \u25B6"}
                                    </button>
                                </div>
                            )}
                        </>
                    )}

                    {!search.loading && !search.error && search.query && audioQuotes.length === 0 && (
                        <div className="builder-search-empty">No voice clips found. Try a different search.</div>
                    )}
                </>
            )}
        </div>
    );
}
