import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { getGridColour, getThemeColours, zoomConfig } from "./chartConfig";
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

    const tc = getThemeColours();
    const gridColour = getGridColour();
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
                        backgroundColor: tc.gold,
                        borderColor: tc.goldDark,
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
                        ticks: { color: tc.text },
                    },
                },
            }}
        />
    );
}
