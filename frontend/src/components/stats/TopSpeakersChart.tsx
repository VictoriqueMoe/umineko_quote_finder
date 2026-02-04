import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { zoomConfig } from "./chartConfig";
import type { StatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

interface TopSpeakersChartProps {
    data: StatsResponse;
    onRegister: (id: string, chart: Chart) => void;
}

export function TopSpeakersChart({ data, onRegister }: TopSpeakersChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartTopSpeakers", chartRef.current);
        }
    }, [onRegister]);

    const labels = data.topSpeakers.map(s => s.name);
    const counts = data.topSpeakers.map(s => s.count);

    return (
        <Bar
            ref={chartRef}
            data={{
                labels,
                datasets: [
                    {
                        label: "Lines",
                        data: counts,
                        backgroundColor: "#d4a84b",
                        borderColor: "#a67c2e",
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
                        grid: { color: "rgba(61, 42, 92, 0.4)" },
                        ticks: { color: "#a89bb8" },
                    },
                    y: {
                        grid: { display: false },
                        ticks: { color: "#e8e0f0" },
                    },
                },
            }}
        />
    );
}
