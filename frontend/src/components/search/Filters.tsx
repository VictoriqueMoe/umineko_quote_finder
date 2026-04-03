import { useAppContext } from "../../hooks/useAppContext";
import type { FilterState, Game, ViewMode } from "../../types/app";

interface FiltersProps {
    filters: FilterState;
    viewMode: ViewMode;
    onFilterChange: (filters: Partial<FilterState>) => void;
    onBrowseClick: () => void;
    browseDisabled: boolean;
}

const DEFAULT_FILTERS: FilterState = {
    character: "",
    interactionA: "",
    interactionB: "",
    episode: "0",
    truth: "",
    arc: "",
};

const HIGURASHI_ARCS: { value: string; label: string }[] = [
    { value: "onikakushi", label: "Onikakushi" },
    { value: "watanagashi", label: "Watanagashi" },
    { value: "tatarigoroshi", label: "Tatarigoroshi" },
    { value: "himatsubushi", label: "Himatsubushi" },
    { value: "meakashi", label: "Meakashi" },
    { value: "tsumihoroboshi", label: "Tsumihoroboshi" },
    { value: "minagoroshi", label: "Minagoroshi" },
    { value: "matsuribayashi", label: "Matsuribayashi" },
    { value: "someutsushi", label: "Someutsushi" },
    { value: "kageboshi", label: "Kageboshi" },
    { value: "tsukiotoshi", label: "Tsukiotoshi" },
    { value: "taraimawashi", label: "Taraimawashi" },
    { value: "yoigoshi", label: "Yoigoshi" },
    { value: "tokihogushi", label: "Tokihogushi" },
    { value: "miotsukushi_omote", label: "Miotsukushi Omote" },
    { value: "kakera", label: "Kakera" },
    { value: "miotsukushi_ura", label: "Miotsukushi Ura" },
    { value: "kotohogushi", label: "Kotohogushi" },
    { value: "hajisarashi", label: "Hajisarashi" },
];

function GameSpecificFilters({
    filters,
    game,
    onFilterChange,
}: {
    filters: FilterState;
    game: Game;
    onFilterChange: (f: Partial<FilterState>) => void;
}) {
    if (game === "higurashi") {
        return (
            <>
                <div className="filter-group">
                    <label className="filter-label" htmlFor="filter-arc">
                        Arc
                    </label>
                    <select
                        id="filter-arc"
                        className="truth-select"
                        value={filters.arc}
                        onChange={e => onFilterChange({ arc: e.target.value })}
                    >
                        <option value="">All Arcs</option>
                        {HIGURASHI_ARCS.map(a => (
                            <option key={a.value} value={a.value}>
                                {a.label}
                            </option>
                        ))}
                    </select>
                </div>
            </>
        );
    }

    return (
        <>
            <div className="filter-group">
                <label className="filter-label" htmlFor="filter-episode">
                    Episode
                </label>
                <select
                    id="filter-episode"
                    className="episode-select"
                    value={filters.episode}
                    onChange={e => onFilterChange({ episode: e.target.value })}
                >
                    <option value="0">All Episodes</option>
                    <option value="1">{"Episode 1 \u2014 Legend"}</option>
                    <option value="2">{"Episode 2 \u2014 Turn"}</option>
                    <option value="3">{"Episode 3 \u2014 Banquet"}</option>
                    <option value="4">{"Episode 4 \u2014 Alliance"}</option>
                    <option value="5">{"Episode 5 \u2014 End"}</option>
                    <option value="6">{"Episode 6 \u2014 Dawn"}</option>
                    <option value="7">{"Episode 7 \u2014 Requiem"}</option>
                    <option value="8">{"Episode 8 \u2014 Twilight"}</option>
                </select>
            </div>
            <div className="filter-group">
                <label className="filter-label" htmlFor="filter-truth">
                    Truth
                </label>
                <select
                    id="filter-truth"
                    className="truth-select"
                    value={filters.truth}
                    onChange={e => onFilterChange({ truth: e.target.value })}
                >
                    <option value="">All Quotes</option>
                    <option value="red">Red Truth</option>
                    <option value="blue">Blue Truth</option>
                    <option value="gold">Gold Truth</option>
                    <option value="purple">Purple Statements</option>
                </select>
            </div>
        </>
    );
}

export function Filters({ filters, viewMode, onFilterChange, onBrowseClick, browseDisabled }: FiltersProps) {
    const { sortedCharacters, sortedAdditionalCharacters, game } = useAppContext();
    const characterNameById = new Map(sortedCharacters);
    const ensureValueOption = (value: string): [string, string][] => {
        if (!value) {
            return sortedCharacters;
        }
        if (characterNameById.has(value)) {
            return sortedCharacters;
        }
        const fallbackLabel = /^\d+$/.test(value) ? `Character ID ${value}` : value;
        return [[value, fallbackLabel], ...sortedCharacters];
    };
    const { character, interactionA, interactionB } = filters;
    const hasCharacterFilter = character.length > 0;
    const hasInteractionFilter = interactionA.length > 0 || interactionB.length > 0;
    const hasBothInteractionChars = interactionA.length > 0 && interactionB.length > 0;
    const interactionHalfFilled = hasInteractionFilter && !hasBothInteractionChars;
    const characterConflictsWithInteraction =
        hasBothInteractionChars && character.length > 0 && character !== interactionA && character !== interactionB;

    const handleSwapInteractions = () => {
        onFilterChange({ interactionA: interactionB, interactionB: interactionA });
    };

    const handleClearInteractions = () => {
        onFilterChange({ interactionA: "", interactionB: "" });
    };

    const handleCharacterChange = (value: string) => {
        if (value) {
            onFilterChange({ character: value, interactionA: "", interactionB: "" });
            return;
        }
        onFilterChange({ character: "" });
    };

    const handleInteractionAChange = (value: string) => {
        onFilterChange({ character: "", interactionA: value });
    };

    const handleInteractionBChange = (value: string) => {
        onFilterChange({ character: "", interactionB: value });
    };

    const handleResetFilters = () => {
        onFilterChange({ character: "", interactionA: "", interactionB: "", episode: "0", truth: "", arc: "" });
    };

    const selectedInteractionAName = interactionA ? characterNameById.get(interactionA) || interactionA : "";
    const selectedInteractionBName = interactionB ? characterNameById.get(interactionB) || interactionB : "";
    const characterOptions = ensureValueOption(character);
    const interactionAOptions = ensureValueOption(interactionA);
    const interactionBOptions = ensureValueOption(interactionB);
    const hasAnyFilter = (Object.keys(DEFAULT_FILTERS) as (keyof FilterState)[]).some(
        key => filters[key] !== DEFAULT_FILTERS[key],
    );

    return (
        <section className="filter-section">
            <p className="filter-unified-note">All filters below apply to both Search and Browse.</p>
            <div className="filter-row">
                <div className="filter-group">
                    <label className="filter-label" htmlFor="filter-character">
                        Character
                    </label>
                    <select
                        id="filter-character"
                        className="character-select"
                        value={character}
                        onChange={e => handleCharacterChange(e.target.value)}
                        disabled={hasInteractionFilter}
                    >
                        <option value="">All Characters</option>
                        {characterOptions.map(([id, name]) => (
                            <option key={id} value={id}>
                                {name}
                            </option>
                        ))}
                        {sortedAdditionalCharacters.length > 0 && (
                            <optgroup label="Additional Characters">
                                {sortedAdditionalCharacters.map(([id, name]) => (
                                    <option key={id} value={id}>
                                        {name}
                                    </option>
                                ))}
                            </optgroup>
                        )}
                    </select>
                </div>
                <GameSpecificFilters filters={filters} game={game} onFilterChange={onFilterChange} />
            </div>
            <div className="interaction-filter-panel">
                <div className="interaction-filter-head">
                    <p className="interaction-filter-title">Interaction Pair (optional)</p>
                </div>
                <div className="interaction-filter-row">
                    <div className="filter-group">
                        <label className="filter-label" htmlFor="filter-interaction-a">
                            Character A
                        </label>
                        <select
                            id="filter-interaction-a"
                            className="interaction-select"
                            value={interactionA}
                            onChange={e => handleInteractionAChange(e.target.value)}
                            disabled={hasCharacterFilter}
                        >
                            <option value="">Select character</option>
                            {interactionAOptions.map(([id, name]) => (
                                <option key={id} value={id}>
                                    {name}
                                </option>
                            ))}
                            {sortedAdditionalCharacters.length > 0 && (
                                <optgroup label="Additional Characters">
                                    {sortedAdditionalCharacters.map(([id, name]) => (
                                        <option key={id} value={id}>
                                            {name}
                                        </option>
                                    ))}
                                </optgroup>
                            )}
                        </select>
                    </div>
                    <div className="filter-group">
                        <label className="filter-label" htmlFor="filter-interaction-b">
                            Character B
                        </label>
                        <select
                            id="filter-interaction-b"
                            className="interaction-select"
                            value={interactionB}
                            onChange={e => handleInteractionBChange(e.target.value)}
                            disabled={hasCharacterFilter}
                        >
                            <option value="">Select character</option>
                            {interactionBOptions.map(([id, name]) => (
                                <option key={id} value={id}>
                                    {name}
                                </option>
                            ))}
                            {sortedAdditionalCharacters.length > 0 && (
                                <optgroup label="Additional Characters">
                                    {sortedAdditionalCharacters.map(([id, name]) => (
                                        <option key={id} value={id}>
                                            {name}
                                        </option>
                                    ))}
                                </optgroup>
                            )}
                        </select>
                    </div>
                    <div className="filter-group">
                        <label className="filter-label">&nbsp;</label>
                        <div className="interaction-filter-actions">
                            <button
                                className="interaction-action-btn"
                                disabled={!hasInteractionFilter || hasCharacterFilter}
                                onClick={handleSwapInteractions}
                            >
                                Swap
                            </button>
                            <button
                                className="interaction-action-btn"
                                disabled={!hasInteractionFilter}
                                onClick={handleClearInteractions}
                            >
                                Clear
                            </button>
                        </div>
                    </div>
                </div>
                <p className="interaction-filter-hint">
                    Activates only when both characters are selected. Uses adjacent A/B exchanges.
                </p>
                <p className="interaction-filter-state">
                    {hasCharacterFilter
                        ? "Character filter is active, so Interaction Pair is disabled."
                        : interactionHalfFilled
                          ? "Interaction filter inactive: select both Character A and Character B."
                          : hasBothInteractionChars
                            ? `Interaction filter active: ${selectedInteractionAName} x ${selectedInteractionBName}.${viewMode === "browse" ? " Browse is in interaction mode." : ""}${characterConflictsWithInteraction ? " Character filter is outside the selected pair, so Search may return zero results." : ""}`
                            : "Interaction filter inactive."}
                </p>
            </div>
            <div className="filter-actions-row">
                <button className="browse-btn" disabled={browseDisabled} onClick={onBrowseClick}>
                    Browse Dialogue
                </button>
                <button
                    className="interaction-action-btn filter-reset-btn"
                    disabled={!hasAnyFilter}
                    onClick={handleResetFilters}
                >
                    Reset Filters
                </button>
            </div>
        </section>
    );
}
