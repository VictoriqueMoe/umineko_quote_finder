import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { getGridColour, getPalette, getThemeColours, zoomConfig } from "./chartConfig";
import type { HigurashiStatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

const ARC_ORDER = [
    "onikakushi",
    "watanagashi",
    "tatarigoroshi",
    "himatsubushi",
    "meakashi",
    "tsumihoroboshi",
    "minagoroshi",
    "matsuribayashi",
    "someutsushi",
    "kageboshi",
    "tsukiotoshi",
    "taraimawashi",
    "yoigoshi",
    "tokihogushi",
    "miotsukushi_omote",
    "kakera",
    "miotsukushi_ura",
    "kotohogushi",
    "hajisarashi",
];

interface LinesPerArcChartProps {
    data: HigurashiStatsResponse;
    onRegister: (id: string, chart: Chart) => void;
}

export function LinesPerArcChart({ data, onRegister }: LinesPerArcChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartLinesPerArc", chartRef.current);
        }
    }, [onRegister]);

    const arcKeys = ARC_ORDER.filter(a => a in data.linesPerArc);
    const arcLabels = arcKeys.map(a => a.charAt(0).toUpperCase() + a.slice(1).replace(/_/g, " "));

    const palette = getPalette();
    const tc = getThemeColours();
    const gridColour = getGridColour();

    const charSet = new Set<string>();
    for (const arc of arcKeys) {
        for (const key of Object.keys(data.linesPerArc[arc])) {
            charSet.add(key);
        }
    }

    const charIds = Array.from(charSet).filter(id => id !== "other");
    charIds.push("other");

    const datasets = charIds.map((id, ci) => ({
        label: id === "other" ? "Other" : data.characterNames[id] || id,
        data: arcKeys.map(arc => data.linesPerArc[arc][id] || 0),
        backgroundColor: palette[ci % palette.length],
    }));

    return (
        <Bar
            ref={chartRef}
            data={{ labels: arcLabels, datasets }}
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
