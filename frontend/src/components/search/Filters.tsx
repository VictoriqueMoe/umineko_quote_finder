import { useAppContext } from "../../hooks/useAppContext";
import type { FilterState } from "../../types/app";

interface FiltersProps {
    filters: FilterState;
    onFilterChange: (filters: Partial<FilterState>) => void;
    onBrowseClick: () => void;
    browseDisabled: boolean;
}

export function Filters({ filters, onFilterChange, onBrowseClick, browseDisabled }: FiltersProps) {
    const { sortedCharacters } = useAppContext();

    return (
        <section className="filter-section">
            <div className="filter-row">
                <div className="filter-group">
                    <label className="filter-label">Character</label>
                    <select
                        className="character-select"
                        value={filters.character}
                        onChange={e => onFilterChange({ character: e.target.value })}
                    >
                        <option value="">All Characters</option>
                        {sortedCharacters.map(([id, name]) => (
                            <option key={id} value={id}>
                                {name}
                            </option>
                        ))}
                    </select>
                </div>
                <div className="filter-group">
                    <label className="filter-label">Episode</label>
                    <select
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
                    <label className="filter-label">Truth</label>
                    <select
                        className="truth-select"
                        value={filters.truth}
                        onChange={e => onFilterChange({ truth: e.target.value })}
                    >
                        <option value="">All Quotes</option>
                        <option value="red">Red Truth</option>
                        <option value="blue">Blue Truth</option>
                    </select>
                </div>
                <div className="filter-group">
                    <label className="filter-label">&nbsp;</label>
                    <button className="browse-btn" disabled={browseDisabled} onClick={onBrowseClick}>
                        Browse Dialogue
                    </button>
                </div>
            </div>
        </section>
    );
}
