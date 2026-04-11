import type { KeyboardEvent } from "react";

interface SearchBarProps {
    value: string;
    onChange: (value: string) => void;
    onSubmit: (value: string) => void;
    exact: boolean;
    onExactChange: (exact: boolean) => void;
}

export function SearchBar({ value, onChange, onSubmit, exact, onExactChange }: SearchBarProps) {
    const handleKeyPress = (e: KeyboardEvent) => {
        if (e.key === "Enter") {
            onSubmit(value);
        }
    };

    return (
        <>
            <div className="search-wrapper">
                <span className="search-icon">{"\uD83E\uDD8B"}</span>
                <input
                    type="text"
                    className="search-input"
                    placeholder="Search for truth within the fragments..."
                    autoComplete="off"
                    value={value}
                    onChange={e => onChange(e.target.value)}
                    onKeyDown={handleKeyPress}
                />
                <button className="search-btn" onClick={() => onSubmit(value)}>
                    Search
                </button>
            </div>
            <button
                type="button"
                className={`search-exact-toggle${exact ? " is-active" : ""}`}
                role="switch"
                aria-checked={exact}
                onClick={() => onExactChange(!exact)}
                title="Match the whole word only — searching &ldquo;broth&rdquo; won&rsquo;t match &ldquo;Brother&rdquo;"
            >
                <span className="search-exact-toggle__track">
                    <span className="search-exact-toggle__thumb" />
                </span>
                <span className="search-exact-toggle__label">Whole word only</span>
            </button>
        </>
    );
}
