import { useCallback, useEffect, useRef, useState } from "react";
import type { FilterState, Game, Language } from "./types/app";
import { resolveLanguage } from "./types/app";
import { useAppContext } from "./hooks/useAppContext";
import { useTheme } from "./hooks/useTheme";
import { useAudioPlayer } from "./hooks/useAudioPlayer";
import { useSearch } from "./hooks/useSearch";
import { useBrowse } from "./hooks/useBrowse";
import { useStats } from "./hooks/useStats";
import { useFeaturedQuote } from "./hooks/useFeaturedQuote";
import { type ParsedRoute, useRouter } from "./hooks/useRouter";
import { enforceMutuallyExclusiveFilters, normalizeFilterCharacters } from "./utils/filters";
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
import { VoiceBuilderView } from "./components/builder/VoiceBuilderView";
import { LoadingSpinner } from "./components/common/LoadingSpinner";
import { EmptyState } from "./components/common/EmptyState";
import { BookmarksView } from "./components/quotes/BookmarksView";
import { switchBookmarkGame, useBookmarks } from "./hooks/useBookmarks";
import { normalizeCharacterKey } from "./utils/characterIds";

const DEFAULT_FILTERS: FilterState = {
    character: "",
    interactionA: "",
    interactionB: "",
    episode: "0",
    truth: "",
    arc: "",
};

export default function App() {
    const { language, setLanguage, game, setGame } = useAppContext();
    const { particlesEnabled, switchThemeForGame } = useTheme();
    const audioPlayer = useAudioPlayer();
    const search = useSearch();
    const browse = useBrowse();
    const stats = useStats();
    const featured = useFeaturedQuote();
    const { count: bookmarkCount } = useBookmarks();

    const [filters, setFilters] = useState<FilterState>(DEFAULT_FILTERS);
    const [searchInputValue, setSearchInputValue] = useState("");
    const [audioIdInputValue, setAudioIdInputValue] = useState("");
    const [builderInitialSegments, setBuilderInitialSegments] = useState<string | null>(null);
    const pendingDeeplinkScroll = useRef(false);
    const resultsSectionRef = useRef<HTMLElement>(null);

    const { viewMode, navigate } = useRouter({
        language,
        game,
        onRouteLoad: async (route: ParsedRoute) => {
            setLanguage(route.lang);
            if (route.game !== game) {
                setGame(route.game);
                switchBookmarkGame(route.game);
                switchThemeForGame(route.game);
            }

            const baseFilters = {
                episode: route.episode,
                truth: route.truth,
                arc: route.arc,
                interactionA: route.interactionA,
                interactionB: route.interactionB,
            };

            switch (route.viewMode) {
                case "search": {
                    const nextFilters = enforceMutuallyExclusiveFilters(
                        normalizeFilterCharacters({ character: route.character, ...baseFilters }, route.game),
                    );
                    setSearchInputValue(route.query);
                    setFilters(prev => ({ ...prev, ...nextFilters }));
                    await search.search(route.game, route.query, route.lang, route.offset, nextFilters);
                    break;
                }
                case "browse": {
                    const nextFilters = enforceMutuallyExclusiveFilters(
                        normalizeFilterCharacters({ character: route.character, ...baseFilters }, route.game),
                    );
                    setFilters(prev => ({ ...prev, ...nextFilters }));
                    await browse.browse(
                        route.game,
                        nextFilters.character,
                        resolveLanguage(route.lang),
                        route.offset,
                        nextFilters.interactionA,
                        nextFilters.interactionB,
                        nextFilters.episode,
                        nextFilters.truth,
                        nextFilters.arc,
                    );
                    break;
                }
                case "stats": {
                    setFilters(prev => ({ ...prev, ...baseFilters }));
                    await stats.loadStats(route.game, route.episode);
                    break;
                }
                case "quoteLookup": {
                    setFilters(prev => ({ ...prev, ...baseFilters }));
                    pendingDeeplinkScroll.current = true;
                    await featured.lookupByAudioId(route.game, route.audioId, resolveLanguage(route.lang));
                    break;
                }
                case "voiceBuilder": {
                    setFilters(prev => ({ ...prev, ...baseFilters }));
                    setBuilderInitialSegments(route.segments);
                    break;
                }
                case "featured": {
                    setFilters(prev => ({ ...prev, ...baseFilters }));
                    await featured.randomQuote(route.game, resolveLanguage(route.lang), {
                        character: "",
                        ...baseFilters,
                    });
                    break;
                }
            }
        },
    });

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

    const handleSearchSubmit = useCallback(
        async (query: string) => {
            setSearchInputValue(query);
            audioPlayer.stop();
            const result = await search.search(game, query, language, 0, filters);
            if (result) {
                pendingDeeplinkScroll.current = true;
                navigate("search", filters, { searchOffset: result.offset, searchQuery: query });
            }
        },
        [game, filters, language, audioPlayer, search, navigate],
    );

    const handleSearchPaginate = useCallback(
        async (newOffset: number) => {
            audioPlayer.stop();
            const result = await search.search(game, search.query, language, newOffset, filters);
            if (result) {
                navigate("search", filters, { searchOffset: result.offset, searchQuery: search.query });
            }
        },
        [game, filters, language, audioPlayer, search, navigate],
    );

    const handleRandomQuote = useCallback(async () => {
        audioPlayer.stop();
        const result = await featured.randomQuote(game, resolveLanguage(language), filters);
        if (result) {
            navigate("featured", filters, { currentAudioId: result.audioId });
        }
    }, [game, filters, language, audioPlayer, featured, navigate]);

    const handleQuoteLookup = useCallback(
        async (audioId: string) => {
            audioPlayer.stop();
            const result = await featured.lookupByAudioId(game, audioId, resolveLanguage(language));
            if (result) {
                pendingDeeplinkScroll.current = true;
                navigate("quoteLookup", filters, { currentAudioId: result.audioId });
            }
        },
        [game, filters, language, audioPlayer, featured, navigate],
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
        const hasInteractionPair = !!filters.interactionA && !!filters.interactionB;
        const hasEpisodeFilter = filters.episode !== "0";
        const hasArcFilter = !!filters.arc;
        if (!filters.character && !filters.truth && !hasInteractionPair && !hasEpisodeFilter && !hasArcFilter) {
            return;
        }
        audioPlayer.stop();
        setSearchInputValue("");
        const result = await browse.browse(
            game,
            filters.character,
            resolveLanguage(language),
            0,
            filters.interactionA,
            filters.interactionB,
            filters.episode,
            filters.truth,
            filters.arc,
        );
        if (result) {
            navigate("browse", filters, { browseOffset: result.offset });
        }
    }, [game, filters, language, audioPlayer, browse, navigate]);

    const handleBrowsePaginate = useCallback(
        async (newOffset: number) => {
            audioPlayer.stop();
            const result = await browse.browse(
                game,
                filters.character,
                resolveLanguage(language),
                newOffset,
                filters.interactionA,
                filters.interactionB,
                filters.episode,
                filters.truth,
                filters.arc,
            );
            if (result) {
                navigate("browse", filters, { browseOffset: result.offset });
            }
        },
        [game, filters, language, audioPlayer, browse, navigate],
    );

    const handleLoadStats = useCallback(async () => {
        audioPlayer.stop();
        await stats.loadStats(game, filters.episode);
        navigate("stats", filters);
    }, [game, filters, audioPlayer, stats, navigate]);

    const handleViewInteractionDialogues = useCallback(
        async (interactionA: string, interactionB: string) => {
            if (!interactionA || !interactionB || interactionA === interactionB) {
                return;
            }

            audioPlayer.stop();
            setSearchInputValue("");

            const nextFilters: FilterState = {
                ...filters,
                character: "",
                interactionA: normalizeCharacterKey(interactionA, game),
                interactionB: normalizeCharacterKey(interactionB, game),
            };
            setFilters(nextFilters);

            const result = await browse.browse(
                game,
                "",
                resolveLanguage(language),
                0,
                nextFilters.interactionA,
                nextFilters.interactionB,
                nextFilters.episode,
                nextFilters.truth,
                nextFilters.arc,
            );
            if (result) {
                navigate("browse", nextFilters, { browseOffset: result.offset });
            }
        },
        [game, audioPlayer, filters, browse, language, navigate],
    );

    const handleHomeClick = useCallback(() => {
        audioPlayer.stop();
        featured.randomQuote(game, resolveLanguage(language), filters);
        navigate("featured", filters);
    }, [game, audioPlayer, language, filters, featured, navigate]);

    const handleBuilderClick = useCallback(() => {
        audioPlayer.stop();
        navigate("voiceBuilder", filters);
    }, [audioPlayer, filters, navigate]);

    const handleBookmarksClick = useCallback(() => {
        audioPlayer.stop();
        navigate("bookmarks", filters);
    }, [audioPlayer, filters, navigate]);

    const handleBuilderClose = useCallback(() => {
        audioPlayer.stop();
        featured.randomQuote(game, resolveLanguage(language), filters);
        navigate("featured", filters);
    }, [game, audioPlayer, language, filters, featured, navigate]);

    const handleClear = useCallback(() => {
        audioPlayer.stop();
        search.clear();
        browse.clear();
        stats.clear();
        featured.clear();
        setFilters(DEFAULT_FILTERS);
        setSearchInputValue("");
        setAudioIdInputValue("");
        navigate("featured", DEFAULT_FILTERS);
    }, [audioPlayer, search, browse, stats, featured, navigate]);

    const handleGameChange = useCallback(
        (newGame: Game) => {
            if (newGame === game) {
                return;
            }
            audioPlayer.stop();
            setGame(newGame);
            switchBookmarkGame(newGame);
            switchThemeForGame(newGame);
            search.clear();
            browse.clear();
            stats.clear();
            featured.clear();
            setFilters(DEFAULT_FILTERS);
            setSearchInputValue("");
            setAudioIdInputValue("");
            if (newGame === "higurashi") {
                setLanguage("auto");
            }
            featured.randomQuote(newGame, "en", DEFAULT_FILTERS);
            navigate("featured", DEFAULT_FILTERS, { game: newGame });
        },
        [game, audioPlayer, setGame, setLanguage, switchThemeForGame, search, browse, stats, featured, navigate],
    );

    const handleFilterChange = useCallback(
        (newFilters: Partial<FilterState>) => {
            const merged = enforceMutuallyExclusiveFilters(
                normalizeFilterCharacters({ ...filters, ...newFilters }, game),
            );
            setFilters(merged);

            if (viewMode === "stats") {
                if ("episode" in newFilters) {
                    audioPlayer.stop();
                    stats.loadStats(game, merged.episode).then(() => {
                        navigate("stats", merged);
                    });
                }
            } else if (viewMode === "browse") {
                if (
                    "character" in newFilters ||
                    "episode" in newFilters ||
                    "truth" in newFilters ||
                    "arc" in newFilters ||
                    "interactionA" in newFilters ||
                    "interactionB" in newFilters
                ) {
                    audioPlayer.stop();
                    browse
                        .browse(
                            game,
                            merged.character,
                            resolveLanguage(language),
                            0,
                            merged.interactionA,
                            merged.interactionB,
                            merged.episode,
                            merged.truth,
                            merged.arc,
                        )
                        .then(result => {
                            if (result) {
                                navigate("browse", merged, { browseOffset: result.offset });
                            }
                        });
                }
            } else if (searchInputValue.trim()) {
                audioPlayer.stop();
                search.search(game, searchInputValue, language, 0, merged).then(result => {
                    if (result) {
                        navigate("search", merged, { searchOffset: result.offset, searchQuery: searchInputValue });
                    }
                });
            }
        },
        [game, filters, viewMode, searchInputValue, language, audioPlayer, search, browse, stats, navigate],
    );

    const handleLanguageChange = useCallback(
        (lang: Language) => {
            setLanguage(lang);
            if (game === "higurashi") {
                return;
            }
            const resolved = resolveLanguage(lang);
            if (viewMode === "browse") {
                browse.browse(
                    game,
                    filters.character,
                    resolved,
                    browse.offset,
                    filters.interactionA,
                    filters.interactionB,
                    filters.episode,
                    filters.truth,
                    filters.arc,
                );
            } else if (viewMode === "search" && searchInputValue.trim()) {
                search.search(game, searchInputValue, lang, search.offset, filters);
            } else if (featured.currentAudioId) {
                featured.lookupByAudioId(game, featured.currentAudioId, resolved);
            } else {
                featured.randomQuote(game, resolved, filters);
            }
        },
        [game, setLanguage, viewMode, filters, searchInputValue, browse, search, featured],
    );

    const handleContextQuoteClick = useCallback(
        (audioId: string) => {
            handleQuoteLookup(audioId);
        },
        [handleQuoteLookup],
    );

    useEffect(() => {
        if (!pendingDeeplinkScroll.current || !resultsSectionRef.current) {
            return;
        }
        const ready =
            (viewMode === "quoteLookup" && featured.quote) ||
            (viewMode === "search" && search.results.length > 0) ||
            (viewMode === "featured" && featured.quote);
        if (ready) {
            pendingDeeplinkScroll.current = false;
            resultsSectionRef.current.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    }, [viewMode, featured.quote, search.results]);

    const isStatsActive = viewMode === "stats";
    const isBuilderActive = viewMode === "voiceBuilder" && game === "umineko";

    const browseDisabled =
        !!filters.interactionA !== !!filters.interactionB ||
        (!filters.character &&
            !filters.truth &&
            !filters.arc &&
            filters.episode === "0" &&
            !(filters.interactionA && filters.interactionB));

    return (
        <>
            {particlesEnabled && <Butterflies />}
            <div className="bg-pattern" />
            <div className={`container${isStatsActive ? " stats-active" : ""}`}>
                <Header
                    language={language}
                    game={game}
                    viewMode={viewMode}
                    onLanguageChange={handleLanguageChange}
                    onGameChange={handleGameChange}
                    onHomeClick={handleHomeClick}
                    onStatsClick={handleLoadStats}
                    onBuilderClick={handleBuilderClick}
                    onBookmarksClick={handleBookmarksClick}
                    bookmarkCount={bookmarkCount}
                />

                {isBuilderActive ? (
                    <VoiceBuilderView onClose={handleBuilderClose} initialBuilder={builderInitialSegments} />
                ) : (
                    <>
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
                            viewMode={viewMode}
                            onFilterChange={handleFilterChange}
                            onBrowseClick={handleBrowseClick}
                            browseDisabled={browseDisabled}
                        />

                        <section
                            ref={resultsSectionRef}
                            className={`results-section${loading && hasViewData ? " results-loading" : ""}`}
                        >
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
                                    filters={filters}
                                    onContextQuoteClick={handleContextQuoteClick}
                                    langOverride={language === "auto" ? search.detectedLang : undefined}
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
                                <StatsView
                                    data={stats.data}
                                    game={game}
                                    episode={filters.episode}
                                    onViewInteractionDialogues={handleViewInteractionDialogues}
                                />
                            )}
                            {viewMode === "bookmarks" && (
                                <BookmarksView
                                    audioPlayer={audioPlayer}
                                    onContextQuoteClick={handleContextQuoteClick}
                                />
                            )}
                        </section>
                    </>
                )}

                <Footer />
            </div>
        </>
    );
}
