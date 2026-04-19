import { useEffect, useRef, useState } from "react";
import { useTheme } from "../../hooks/useTheme";
import { useAppContext } from "../../hooks/useAppContext";
import type { ThemeType } from "../../types/app";
import { ToggleSwitch } from "../common/ToggleSwitch";

interface ThemeDefinition {
    id: ThemeType;
    name: string;
    description: string;
}

const UMINEKO_THEMES: ThemeDefinition[] = [
    { id: "featherine", name: "Featherine", description: "Witch of Theatergoing, Drama, and Spectating" },
    { id: "bernkastel", name: "Lady Bernkastel", description: "Witch of Miracles" },
    { id: "lambdadelta", name: "Lady Lambdadelta", description: "Witch of Certainty" },
];

const HIGURASHI_THEMES: ThemeDefinition[] = [
    { id: "rika", name: "Rika", description: "The one who knows the truth of June 1983" },
    { id: "mion", name: "Mion", description: "The eldest of the Sonozaki twins" },
    { id: "satoko", name: "Satoko", description: "The master of traps" },
];

const CICONIA_THEMES: ThemeDefinition[] = [
    { id: "miyao", name: "Miyao", description: "AOU Gauntlet Knight — sky and sun" },
    { id: "lingji", name: "Lingji", description: "COU Gauntlet Knight — red banner, gold star" },
    { id: "stanislaw", name: "Stanisław", description: "ABN Gauntlet Knight — constellation over the void" },
];

export function ThemeSelector() {
    const { theme, setTheme, particlesEnabled, setParticlesEnabled } = useTheme();
    const { game } = useAppContext();
    const [isOpen, setIsOpen] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    const themes = game === "higurashi" ? HIGURASHI_THEMES : game === "ciconia" ? CICONIA_THEMES : UMINEKO_THEMES;
    const currentTheme = themes.find(t => t.id === theme) || themes[0];

    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        }

        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    return (
        <div className="theme-selector" ref={dropdownRef}>
            <button
                className="theme-trigger"
                onClick={() => setIsOpen(!isOpen)}
                aria-expanded={isOpen}
                aria-haspopup="listbox"
            >
                <span className="theme-trigger-label">Theme</span>
                <span className="theme-trigger-sep">{"\u2726"}</span>
                <span className="theme-trigger-name">{currentTheme.name}</span>
                <span className={`theme-chevron${isOpen ? " open" : ""}`}>{"\u25BC"}</span>
            </button>

            {isOpen && (
                <div className="theme-dropdown" role="listbox">
                    {themes.map(t => (
                        <button
                            key={t.id}
                            className={`theme-option${t.id === theme ? " active" : ""}`}
                            onClick={() => {
                                setTheme(t.id);
                                setIsOpen(false);
                            }}
                            role="option"
                            aria-selected={t.id === theme}
                        >
                            <div className="theme-option-info">
                                <span className="theme-option-name">{t.name}</span>
                                <span className="theme-option-desc">{t.description}</span>
                            </div>
                            {t.id === theme && <span className="theme-check">{"\u2713"}</span>}
                        </button>
                    ))}
                    <div className="theme-dropdown-divider" />
                    <ToggleSwitch
                        enabled={particlesEnabled}
                        onChange={setParticlesEnabled}
                        label="Particles"
                        description="Floating butterflies & sparkles"
                    />
                </div>
            )}
        </div>
    );
}
