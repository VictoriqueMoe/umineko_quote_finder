import { useCallback, useRef, useState } from "react";
import type { FilterState, ViewMode } from "./types/app";
import { useAppContext } from "./hooks/useAppContext";
import { useAudioPlayer } from "./hooks/useAudioPlayer";
import { useSearch } from "./hooks/useSearch";
import { useBrowse } from "./hooks/useBrowse";
import { useStats } from "./hooks/useStats";
import { useFeaturedQuote } from "./hooks/useFeaturedQuote";
import { pushUrl, useUrlState } from "./hooks/useUrlState";
import { Header } from "./components/layout/Header";
import { Footer } from "./components/layout/Footer";
import { Butterflies } from "./components/layout/Butterflies";
import { SearchBar } from "./components/search/SearchBar";
import { AudioIdLookup } from "./components/search/AudioIdLookup";
import { ActionButtons } from "./components/search/ActionButtons";
import { Filters } from "./components/search/Filters";
import { QuoteList } from "./components/quotes/QuoteList";
import { FeaturedQuote } from "./components/quotes/FeaturedQuote";
import { BrowseView } from "./components/quotes/BrowseView";
import { StatsView } from "./components/stats/StatsView";
import { LoadingSpinner } from "./components/common/LoadingSpinner";
import { EmptyState } from "./components/common/EmptyState";

const DEFAULT_FILTERS: FilterState = { character: "", episode: "0", truth: "" };

export default function App() {
    const { language, setLanguage } = useAppContext();
    const audioPlayer = useAudioPlayer();
    const search = useSearch();
    const browse = useBrowse();
    const stats = useStats();
    const featured = useFeaturedQuote();

    const [viewMode, setViewMode] = useState<ViewMode>("featured");
    const [filters, setFilters] = useState<FilterState>(DEFAULT_FILTERS);
    const [searchInputValue, setSearchInputValue] = useState("");
    const [audioIdInputValue, setAudioIdInputValue] = useState("");
    const urlInitialised = useRef(false);

    const loading = search.loading || browse.loading || stats.loading || featured.loading;
    const error =
        (viewMode === "search" && search.error) ||
        (viewMode === "browse" && browse.error) ||
        (viewMode === "stats" && stats.error) ||
        ((viewMode === "featured" || viewMode === "quoteLookup") && featured.error) ||
        null;
    const hasViewData =
        (viewMode === "search" && !!search.query) ||
        (viewMode === "browse" && !!browse.data) ||
        ((viewMode === "featured" || viewMode === "quoteLookup") && !!featured.quote) ||
        (viewMode === "stats" && !!stats.data);

    const doPushUrl = useCallback(
        (
            vm: ViewMode,
            f: FilterState,
            opts?: {
                searchOffset?: number;
                browseOffset?: number;
                currentAudioId?: string | null;
                searchQuery?: string;
            },
        ) => {
            if (!urlInitialised.current) {
                return;
            }
            pushUrl(
                {
                    viewMode: vm,
                    filters: f,
                    currentAudioId: opts?.currentAudioId ?? null,
                    searchOffset: opts?.searchOffset ?? 0,
                    browseOffset: opts?.browseOffset ?? 0,
                },
                language,
                opts?.searchQuery ?? "",
            );
        },
        [language],
    );

    const handleSearchSubmit = useCallback(
        async (query: string) => {
            setSearchInputValue(query);
            audioPlayer.stop();
            const result = await search.search(query, language, 0, filters);
            if (result) {
                setViewMode("search");
                doPushUrl("search", filters, { searchOffset: result.offset, searchQuery: query });
            }
        },
        [filters, language, audioPlayer, search, doPushUrl],
    );

    const handleSearchPaginate = useCallback(
        async (newOffset: number) => {
            audioPlayer.stop();
            const result = await search.search(search.query, language, newOffset, filters);
            if (result) {
                doPushUrl("search", filters, { searchOffset: result.offset, searchQuery: search.query });
            }
        },
        [filters, language, audioPlayer, search, doPushUrl],
    );

    const handleRandomQuote = useCallback(async () => {
        audioPlayer.stop();
        const result = await featured.randomQuote(language, filters);
        if (result) {
            setViewMode("featured");
            doPushUrl("featured", filters, { currentAudioId: result.audioId });
        }
    }, [filters, language, audioPlayer, featured, doPushUrl]);

    const handleQuoteLookup = useCallback(
        async (audioId: string) => {
            audioPlayer.stop();
            const result = await featured.lookupByAudioId(audioId, language);
            if (result) {
                setViewMode("quoteLookup");
                doPushUrl("quoteLookup", filters, { currentAudioId: result.audioId });
            }
        },
        [filters, language, audioPlayer, featured, doPushUrl],
    );

    const handleAudioIdSubmit = useCallback(
        (audioId: string) => {
            if (audioId.trim()) {
                handleQuoteLookup(audioId.trim());
            }
        },
        [handleQuoteLookup],
    );

    const handleBrowseClick = useCallback(async () => {
        if (!filters.character && !filters.truth) {
            return;
        }
        audioPlayer.stop();
        setSearchInputValue("");
        const result = await browse.browse(filters.character, language, 0, filters.episode, filters.truth);
        if (result) {
            setViewMode("browse");
            doPushUrl("browse", filters, { browseOffset: result.offset });
        }
    }, [filters, language, audioPlayer, browse, doPushUrl]);

    const handleBrowsePaginate = useCallback(
        async (newOffset: number) => {
            audioPlayer.stop();
            const result = await browse.browse(filters.character, language, newOffset, filters.episode, filters.truth);
            if (result) {
                doPushUrl("browse", filters, { browseOffset: result.offset });
            }
        },
        [filters, language, audioPlayer, browse, doPushUrl],
    );

    const handleLoadStats = useCallback(async () => {
        audioPlayer.stop();
        await stats.loadStats(filters.episode);
        setViewMode("stats");
        doPushUrl("stats", filters);
    }, [filters, audioPlayer, stats, doPushUrl]);

    const handleClear = useCallback(() => {
        audioPlayer.stop();
        search.clear();
        browse.clear();
        stats.clear();
        featured.clear();
        setViewMode("featured");
        setFilters(DEFAULT_FILTERS);
        setSearchInputValue("");
        setAudioIdInputValue("");
        pushUrl(
            { viewMode: "featured", filters: DEFAULT_FILTERS, currentAudioId: null, searchOffset: 0, browseOffset: 0 },
            "en",
            "",
        );
    }, [audioPlayer, search, browse, stats, featured]);

    const handleFilterChange = useCallback(
        (newFilters: Partial<FilterState>) => {
            const merged = { ...filters, ...newFilters };
            setFilters(merged);

            if (viewMode === "stats") {
                if ("episode" in newFilters) {
                    audioPlayer.stop();
                    stats.loadStats(merged.episode).then(() => {
                        doPushUrl("stats", merged);
                    });
                }
            } else if (viewMode === "browse") {
                if ("episode" in newFilters || "truth" in newFilters) {
                    audioPlayer.stop();
                    browse.browse(filters.character, language, 0, merged.episode, merged.truth).then(result => {
                        if (result) {
                            doPushUrl("browse", merged, { browseOffset: result.offset });
                        }
                    });
                }
            } else if (searchInputValue.trim()) {
                audioPlayer.stop();
                search.search(searchInputValue, language, 0, merged).then(result => {
                    if (result) {
                        doPushUrl("search", merged, { searchOffset: result.offset, searchQuery: searchInputValue });
                    }
                });
            }
        },
        [filters, viewMode, searchInputValue, language, audioPlayer, search, browse, stats, doPushUrl],
    );

    const handleLanguageChange = useCallback(
        (lang: "en" | "ja") => {
            setLanguage(lang);
            if (viewMode === "browse") {
                browse.browse(filters.character, lang, browse.offset, filters.episode, filters.truth);
            } else if (viewMode === "search" && searchInputValue.trim()) {
                search.search(searchInputValue, lang, search.offset, filters);
            } else if (featured.currentAudioId) {
                featured.lookupByAudioId(featured.currentAudioId, lang);
            } else {
                featured.randomQuote(lang, filters);
            }
        },
        [setLanguage, viewMode, filters, searchInputValue, browse, search, featured],
    );

    // URL state initialisation
    useUrlState({
        onSearch: (query, offset, character, lang) => {
            setSearchInputValue(query);
            setFilters(prev => ({ ...prev, character }));
            search.search(query, lang, offset, { ...filters, character }).then(() => {
                setViewMode("search");
                urlInitialised.current = true;
            });
        },
        onBrowse: (character, offset, episode, lang) => {
            setFilters(prev => ({ ...prev, character, episode }));
            browse.browse(character, lang, offset, episode, filters.truth).then(() => {
                setViewMode("browse");
                urlInitialised.current = true;
            });
        },
        onStats: _lang => {
            stats.loadStats(filters.episode).then(() => {
                setViewMode("stats");
                urlInitialised.current = true;
            });
        },
        onQuoteLookup: (audioId, lang) => {
            featured.lookupByAudioId(audioId, lang).then(() => {
                setViewMode("quoteLookup");
                urlInitialised.current = true;
            });
        },
        onDefault: lang => {
            featured.randomQuote(lang, filters).then(() => {
                setViewMode("featured");
                urlInitialised.current = true;
            });
        },
        setLanguage,
        setFilters: f => {
            setFilters(prev => ({ ...prev, ...f }));
        },
    });

    const handleContextQuoteClick = useCallback(
        (audioId: string) => {
            handleQuoteLookup(audioId);
        },
        [handleQuoteLookup],
    );

    const isStatsActive = viewMode === "stats";

    return (
        <>
            <Butterflies />
            <div className="bg-pattern" />
            <div className={`container${isStatsActive ? " stats-active" : ""}`}>
                <Header language={language} onLanguageChange={handleLanguageChange} onStatsClick={handleLoadStats} />

                <section className="search-section">
                    <div className="search-container">
                        <SearchBar
                            value={searchInputValue}
                            onChange={setSearchInputValue}
                            onSubmit={handleSearchSubmit}
                        />
                        <AudioIdLookup
                            value={audioIdInputValue}
                            onChange={setAudioIdInputValue}
                            onSubmit={handleAudioIdSubmit}
                        />
                        <ActionButtons onRandom={handleRandomQuote} onClear={handleClear} />
                    </div>
                </section>

                <Filters
                    filters={filters}
                    onFilterChange={handleFilterChange}
                    onBrowseClick={handleBrowseClick}
                    browseDisabled={!filters.character && !filters.truth}
                />

                <section className={`results-section${loading && hasViewData ? " results-loading" : ""}`}>
                    {loading && !hasViewData && <LoadingSpinner />}
                    {!loading && error && <EmptyState message={error} />}
                    {!error && viewMode === "search" && !!search.query && (
                        <QuoteList
                            results={search.results}
                            query={search.query}
                            total={search.total}
                            offset={search.offset}
                            onPaginate={handleSearchPaginate}
                            audioPlayer={audioPlayer}
                            onContextQuoteClick={handleContextQuoteClick}
                        />
                    )}
                    {!error && (viewMode === "featured" || viewMode === "quoteLookup") && featured.quote && (
                        <FeaturedQuote
                            quote={featured.quote}
                            audioPlayer={audioPlayer}
                            onContextQuoteClick={handleContextQuoteClick}
                        />
                    )}
                    {!error && viewMode === "browse" && browse.data && (
                        <BrowseView
                            data={browse.data}
                            offset={browse.offset}
                            total={browse.total}
                            onPaginate={handleBrowsePaginate}
                            audioPlayer={audioPlayer}
                            filters={filters}
                            onContextQuoteClick={handleContextQuoteClick}
                        />
                    )}
                    {!error && viewMode === "stats" && stats.data && (
                        <StatsView data={stats.data} episode={filters.episode} />
                    )}
                </section>

                <Footer />
            </div>
        </>
    );
}
