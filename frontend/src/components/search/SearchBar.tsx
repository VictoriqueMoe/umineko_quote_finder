import type { KeyboardEvent } from "react";

interface SearchBarProps {
    value: string;
    onChange: (value: string) => void;
    onSubmit: (value: string) => void;
}

export function SearchBar({ value, onChange, onSubmit }: SearchBarProps) {
    const handleKeyPress = (e: KeyboardEvent) => {
        if (e.key === "Enter") {
            onSubmit(value);
        }
    };

    return (
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
    );
}
