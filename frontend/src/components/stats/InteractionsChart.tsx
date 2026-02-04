import { useEffect, useRef } from "react";
import { Bar } from "react-chartjs-2";
import { zoomConfig } from "./chartConfig";
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
                        backgroundColor: "#9d7bc9",
                        borderColor: "#6b4c9a",
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
                        ticks: { color: "#e8e0f0", font: { size: 11 } },
                    },
                },
            }}
        />
    );
}
