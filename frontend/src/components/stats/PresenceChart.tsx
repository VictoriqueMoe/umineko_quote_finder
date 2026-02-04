import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { PALETTE, zoomConfig } from "./chartConfig";
import type { StatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

interface PresenceChartProps {
    data: StatsResponse;
    onRegister: (id: string, chart: Chart) => void;
}

export function PresenceChart({ data, onRegister }: PresenceChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartPresence", chartRef.current);
        }
    }, [onRegister]);

    const epLabels = Array.from({ length: 8 }, (_, i) => `EP${i + 1}`);

    const datasets = data.characterPresence.map((cp, i) => ({
        label: cp.name,
        data: cp.episodes,
        backgroundColor: PALETTE[i % PALETTE.length],
        borderColor: PALETTE[i % PALETTE.length],
        borderWidth: 1,
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
                        grid: { color: "rgba(61, 42, 92, 0.4)" },
                        ticks: { color: "#a89bb8" },
                    },
                    y: {
                        grid: { color: "rgba(61, 42, 92, 0.4)" },
                        ticks: { color: "#a89bb8" },
                    },
                },
            }}
        />
    );
}
