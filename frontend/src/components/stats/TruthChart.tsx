import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { getGridColour, getThemeColours, zoomConfig } from "./chartConfig";
import type { StatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

interface TruthChartProps {
    data: StatsResponse;
    onRegister: (id: string, chart: Chart) => void;
}

export function TruthChart({ data, onRegister }: TruthChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartTruth", chartRef.current);
        }
    }, [onRegister]);

    const tc = getThemeColours();
    const gridColour = getGridColour();
    const labels = data.truthPerEpisode.map(t => `EP${t.episode}`);
    const redData = data.truthPerEpisode.map(t => t.red);
    const blueData = data.truthPerEpisode.map(t => t.blue);
    const goldData = data.truthPerEpisode.map(t => t.gold);
    const purpleData = data.truthPerEpisode.map(t => t.purple);

    return (
        <Bar
            ref={chartRef}
            data={{
                labels,
                datasets: [
                    {
                        label: "Red Truth",
                        data: redData,
                        backgroundColor: "#ff3333",
                        borderColor: "#cc0000",
                        borderWidth: 1,
                    },
                    {
                        label: "Blue Truth",
                        data: blueData,
                        backgroundColor: "#3399ff",
                        borderColor: "#0066cc",
                        borderWidth: 1,
                    },
                    {
                        label: "Gold Truth",
                        data: goldData,
                        backgroundColor: "#ffaa00",
                        borderColor: "#cc8800",
                        borderWidth: 1,
                    },
                    {
                        label: "Purple Statements",
                        data: purpleData,
                        backgroundColor: "#aa71ff",
                        borderColor: "#8855cc",
                        borderWidth: 1,
                    },
                ],
            }}
            options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: "bottom",
                        labels: { color: tc.textMuted },
                    },
                    zoom: zoomConfig,
                },
                scales: {
                    x: {
                        grid: { color: gridColour },
                        ticks: { color: tc.textMuted },
                    },
                    y: {
                        grid: { color: gridColour },
                        ticks: { color: tc.textMuted },
                    },
                },
            }}
        />
    );
}
