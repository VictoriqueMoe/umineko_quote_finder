import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { getGridColour, getThemeColours, zoomConfig } from "./chartConfig";
import type { StatsResponse } from "../../types/api";
import type { Chart } from "chart.js";

interface InteractionsChartProps {
    data: StatsResponse;
    onRegister: (id: string, chart: Chart) => void;
}

export function InteractionsChart({ data, onRegister }: InteractionsChartProps) {
    const chartRef = useRef<Chart<"bar"> | null>(null);

    useEffect(() => {
        if (chartRef.current) {
            onRegister("chartInteractions", chartRef.current);
        }
    }, [onRegister]);

    const tc = getThemeColours();
    const gridColour = getGridColour();
    const labels = data.interactions.map(i => `${i.nameA} & ${i.nameB}`);
    const counts = data.interactions.map(i => i.count);

    return (
        <Bar
            ref={chartRef}
            data={{
                labels,
                datasets: [
                    {
                        label: "Adjacent Lines",
                        data: counts,
                        backgroundColor: tc.purpleLight,
                        borderColor: tc.purple,
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
                        ticks: { color: tc.text, font: { size: 11 } },
                    },
                },
            }}
        />
    );
}
