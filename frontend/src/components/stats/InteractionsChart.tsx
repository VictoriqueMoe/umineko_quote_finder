import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Bar } from "react-chartjs-2";
import { getGridColour, getThemeColours, zoomConfig } from "./chartConfig";
import type { StatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

interface InteractionsChartProps {
    data: StatsResponse;
    onRegister: (id: string, chart: Chart) => void;
    onViewDialogues?: (charA: string, charB: string) => void;
}

export function InteractionsChart({ data, onRegister, onViewDialogues }: InteractionsChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);
    const [charA, setCharA] = useState("");
    const [charB, setCharB] = useState("");

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartInteractions", chartRef.current);
        }
    }, [onRegister]);

    const characters = useMemo(
        () => Object.entries(data.characterNames).sort((a, b) => a[1].localeCompare(b[1])),
        [data.characterNames],
    );

    const resetPairSelection = useCallback(() => {
        setCharA("");
        setCharB("");
    }, []);

    const swapPairSelection = useCallback(() => {
        setCharA(charB);
        setCharB(charA);
    }, [charA, charB]);

    const tc = getThemeColours();
    const gridColour = getGridColour();
    const pairKey = charA && charB && charA !== charB ? [charA, charB].sort().join("|") : "";
    const selectedPairCount = pairKey ? (data.interactionCounts[pairKey] ?? 0) : 0;

    const sortedPairs = useMemo(
        () => Object.entries(data.interactionCounts).sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0])),
        [data.interactionCounts],
    );
    const rank = pairKey ? sortedPairs.findIndex(([key]) => key === pairKey) + 1 : 0;
    const rankText = rank > 0 ? `#${rank} of ${sortedPairs.length}` : "Not ranked";

    const selectedPairKey = pairKey || null;
    const isRelatedMode = !!selectedPairKey;
    const canViewDialogues = !!onViewDialogues && !!charA && !!charB && charA !== charB;

    const baselineRows = useMemo(
        () =>
            data.interactions.map(i => ({
                key: `${i.charA}|${i.charB}`,
                charA: i.charA,
                charB: i.charB,
                nameA: i.nameA,
                nameB: i.nameB,
                count: i.count,
            })),
        [data.interactions],
    );

    const allRows = useMemo(
        () =>
            Object.entries(data.interactionCounts)
                .map(([key, count]) => {
                    const [a, b] = key.split("|");
                    return {
                        key,
                        charA: a,
                        charB: b,
                        nameA: data.characterNames[a] ?? a,
                        nameB: data.characterNames[b] ?? b,
                        count,
                    };
                })
                .sort((left, right) => right.count - left.count || left.key.localeCompare(right.key)),
        [data.interactionCounts, data.characterNames],
    );

    const chartRows = useMemo(() => {
        if (!selectedPairKey || charA === charB) {
            return baselineRows;
        }
        const [a, b] = selectedPairKey.split("|");
        const selectedPair = allRows.find(p => p.key === selectedPairKey) ?? {
            key: selectedPairKey,
            charA: a,
            charB: b,
            nameA: data.characterNames[a] ?? a,
            nameB: data.characterNames[b] ?? b,
            count: selectedPairCount,
        };

        const related = allRows
            .filter(
                p =>
                    p.key !== selectedPairKey &&
                    (p.charA === charA || p.charB === charA || p.charA === charB || p.charB === charB),
            )
            .slice(0, 11);

        return [selectedPair, ...related].sort(
            (left, right) => right.count - left.count || left.key.localeCompare(right.key),
        );
    }, [baselineRows, allRows, selectedPairKey, charA, charB, selectedPairCount, data.characterNames]);

    const isSelectedVisible = selectedPairKey ? chartRows.some(p => p.key === selectedPairKey) : false;

    const labels = chartRows.map(i => `${i.nameA} & ${i.nameB}`);
    const counts = chartRows.map(i => i.count);
    const chartSubtitle =
        selectedPairKey && charA !== charB
            ? "Selected pair + related interactions (pairs containing Character A or B), sorted by count."
            : "Default top interactions.";
    const barColors = chartRows.map(i => (selectedPairKey && i.key === selectedPairKey ? tc.gold : tc.purpleLight));
    const borderColors = chartRows.map(i => (selectedPairKey && i.key === selectedPairKey ? tc.goldDark : tc.purple));

    return (
        <div className="interactions-explorer">
            <div className="interactions-controls">
                <div className="interactions-control-group">
                    <label className="filter-label" htmlFor="stats-interaction-char-a">
                        Character A
                    </label>
                    <select
                        id="stats-interaction-char-a"
                        className="stats-character-select"
                        value={charA}
                        onChange={e => setCharA(e.target.value)}
                        title="Select Character A for pair analysis"
                    >
                        <option value="">Select character</option>
                        {characters.map(([id, name]) => (
                            <option key={id} value={id}>
                                {name}
                            </option>
                        ))}
                    </select>
                </div>
                <div className="interactions-control-group">
                    <label className="filter-label" htmlFor="stats-interaction-char-b">
                        Character B
                    </label>
                    <select
                        id="stats-interaction-char-b"
                        className="stats-character-select"
                        value={charB}
                        onChange={e => setCharB(e.target.value)}
                        title="Select Character B for pair analysis"
                    >
                        <option value="">Select character</option>
                        {characters.map(([id, name]) => (
                            <option key={id} value={id}>
                                {name}
                            </option>
                        ))}
                    </select>
                </div>
                <button
                    className="interactions-swap-btn"
                    onClick={swapPairSelection}
                    disabled={!charA && !charB}
                    title="Swap Character A and Character B"
                >
                    Swap
                </button>
                <div className="interactions-pair-stats">
                    <span className="interactions-pair-count">
                        {selectedPairKey ? selectedPairCount.toLocaleString() : "\u2014"}
                    </span>
                    <span className="interactions-pair-label">back-to-back speaker switches</span>
                    <span className="interactions-pair-rank">Pair rank: {selectedPairKey ? rankText : "\u2014"}</span>
                </div>
                <button className="interactions-reset-btn" onClick={resetPairSelection} title="Clear both selections">
                    Reset Pair
                </button>
                <button
                    className="interactions-view-btn"
                    disabled={!canViewDialogues}
                    onClick={() => onViewDialogues?.(charA, charB)}
                    title="Open Browse results for this interaction pair"
                >
                    View Dialogues
                </button>
            </div>
            <div className="interactions-mode-row">
                <span className={`interactions-mode-pill${isRelatedMode ? " is-related" : ""}`}>
                    {isRelatedMode ? "Related Mode" : "Default Mode"}
                </span>
                <p className="interactions-mode-text">
                    {isRelatedMode
                        ? "Selected pair + top related pairs (contains Character A or Character B), ordered by count."
                        : "Top global interaction pairs from the default stats list."}
                </p>
            </div>
            <p className="interactions-pair-note">{chartSubtitle}</p>
            <div className="interactions-chart-wrap">
                <Bar
                    ref={chartRef}
                    data={{
                        labels,
                        datasets: [
                            {
                                label: "Back-to-back speaker switches",
                                data: counts,
                                backgroundColor: barColors,
                                borderColor: borderColors,
                                borderWidth: 1,
                            },
                        ],
                    }}
                    options={{
                        indexAxis: "y",
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { display: false },
                            zoom: zoomConfig,
                        },
                        scales: {
                            x: {
                                grid: { color: gridColour },
                                ticks: { color: tc.textMuted },
                            },
                            y: {
                                grid: { display: false },
                                ticks: {
                                    color: tc.text,
                                    font: { size: 11 },
                                    autoSkip: false,
                                },
                            },
                        },
                    }}
                />
            </div>
            {!charA || !charB ? (
                <p className="interactions-pair-note">Select two characters to highlight their pair.</p>
            ) : null}
            {charA && charB && charA === charB && (
                <p className="interactions-pair-note">Pick two different characters to compare.</p>
            )}
            {selectedPairKey && selectedPairCount === 0 && (
                <p className="interactions-pair-note">
                    No back-to-back speaker switches found for this pair in the current scope.
                </p>
            )}
            {selectedPairKey && selectedPairCount > 0 && !isSelectedVisible && (
                <p className="interactions-pair-note">This pair is rank {rankText} globally.</p>
            )}
            {selectedPairKey && selectedPairCount > 0 && isSelectedVisible && (
                <p className="interactions-pair-note">
                    This pair appears in {selectedPairCount.toLocaleString()} back-to-back speaker transitions.
                </p>
            )}
        </div>
    );
}
