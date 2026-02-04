import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { zoomConfig } from "./chartConfig";
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

    const labels = data.truthPerEpisode.map(t => `EP${t.episode}`);
    const redData = data.truthPerEpisode.map(t => t.red);
    const blueData = data.truthPerEpisode.map(t => t.blue);

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
                ],
            }}
            options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: "bottom",
                        labels: { color: "#a89bb8" },
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
