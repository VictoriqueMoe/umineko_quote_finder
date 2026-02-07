import type { KeyboardEvent } from "react";

interface AudioIdLookupProps {
    value: string;
    onChange: (value: string) => void;
    onSubmit: (value: string) => void;
}

export function AudioIdLookup({ value, onChange, onSubmit }: AudioIdLookupProps) {
    const handleKeyPress = (e: KeyboardEvent) => {
        if (e.key === "Enter") {
            onSubmit(value);
        }
    };

    return (
        <div className="search-wrapper audio-id-wrapper">
            <span className="search-icon">{"\u266C"}</span>
            <input
                type="text"
                className="search-input"
                placeholder="Look up by audio ID, e.g. 92100169"
                autoComplete="off"
                value={value}
                onChange={e => onChange(e.target.value)}
                onKeyDown={handleKeyPress}
            />
            <button className="search-btn" onClick={() => onSubmit(value)}>
                Lookup
            </button>
        </div>
    );
}
