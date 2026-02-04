import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { PALETTE, zoomConfig } from "./chartConfig";
import type { StatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

interface LinesPerEpisodeChartProps {
    data: StatsResponse;
    onRegister: (id: string, chart: Chart) => void;
}

export function LinesPerEpisodeChart({ data, onRegister }: LinesPerEpisodeChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartLinesPerEpisode", chartRef.current);
        }
    }, [onRegister]);

    const epLabels = data.linesPerEpisode.map(ep => `EP${ep.episode} ${ep.episodeName}`);

    const charSet = new Set<string>();
    for (const ep of data.linesPerEpisode) {
        for (const key of Object.keys(ep.characters)) {
            charSet.add(key);
        }
    }

    const charIds = Array.from(charSet).filter(id => id !== "other");
    charIds.push("other");

    const datasets = charIds.map((id, ci) => ({
        label: id === "other" ? "Other" : data.characterNames[id] || id,
        data: data.linesPerEpisode.map(ep => ep.characters[id] || 0),
        backgroundColor: PALETTE[ci % PALETTE.length],
    }));

    return (
        <Bar
            ref={chartRef}
            data={{ labels: epLabels, datasets }}
            options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: "bottom",
                        labels: { color: "#a89bb8", boxWidth: 12 },
                    },
                    zoom: zoomConfig,
                },
                scales: {
                    x: {
                        stacked: true,
                        grid: { color: "rgba(61, 42, 92, 0.4)" },
                        ticks: { color: "#a89bb8" },
                    },
                    y: {
                        stacked: true,
                        grid: { color: "rgba(61, 42, 92, 0.4)" },
                        ticks: { color: "#a89bb8" },
                    },
                },
            }}
        />
    );
}
