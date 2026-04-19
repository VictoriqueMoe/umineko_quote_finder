import { useEffect, useMemo, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { getGridColour, getPalette, getThemeColours, zoomConfig } from "./chartConfig";
import type { CiconiaStatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

interface LinesPerChapterChartProps {
    data: CiconiaStatsResponse;
    onRegister: (id: string, chart: Chart) => void;
}

function chapterSortKey(chapter: string): number {
    if (chapter === "00") {
        return 0;
    }
    if (chapter === "25b") {
        return 25.5;
    }
    if (chapter.startsWith("df")) {
        return 100 + parseInt(chapter.slice(2), 10);
    }
    return parseInt(chapter, 10);
}

function chapterLabel(chapter: string): string {
    if (chapter === "00") {
        return "Prologue";
    }
    if (chapter === "25b") {
        return "Ch 25b";
    }
    if (chapter.startsWith("df")) {
        return `DF ${parseInt(chapter.slice(2), 10)}`;
    }
    return `Ch ${parseInt(chapter, 10)}`;
}

export function LinesPerChapterChart({ data, onRegister }: LinesPerChapterChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartLinesPerChapter", chartRef.current);
        }
    }, [onRegister]);

    const chapterKeys = useMemo(() => {
        return Object.keys(data.linesPerChapter).sort((a, b) => chapterSortKey(a) - chapterSortKey(b));
    }, [data.linesPerChapter]);

    const chapterLabels = chapterKeys.map(chapterLabel);

    const palette = getPalette();
    const tc = getThemeColours();
    const gridColour = getGridColour();

    const charSet = new Set<string>();
    for (const ch of chapterKeys) {
        for (const key of Object.keys(data.linesPerChapter[ch])) {
            charSet.add(key);
        }
    }

    const charIds = Array.from(charSet).filter(id => id !== "other");
    charIds.push("other");

    const datasets = charIds.map((id, ci) => ({
        label: id === "other" ? "Other" : data.characterNames[id] || id,
        data: chapterKeys.map(ch => data.linesPerChapter[ch][id] || 0),
        backgroundColor: palette[ci % palette.length],
    }));

    return (
        <Bar
            ref={chartRef}
            data={{ labels: chapterLabels, datasets }}
            options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: "bottom",
                        labels: { color: tc.textMuted, boxWidth: 12 },
                    },
                    zoom: zoomConfig,
                },
                scales: {
                    x: {
                        stacked: true,
                        grid: { color: gridColour },
                        ticks: { color: tc.textMuted },
                    },
                    y: {
                        stacked: true,
                        grid: { color: gridColour },
                        ticks: { color: tc.textMuted },
                    },
                },
            }}
        />
    );
}
