import type { ReactNode } from "react";

interface StatsCardProps {
    id: string;
    title: string;
    tall?: boolean;
    wide?: boolean;
    chartClassName?: string;
    children: ReactNode;
    onResetZoom: (chartId: string) => void;
}

export function StatsCard({ id, title, tall, wide, chartClassName, children, onResetZoom }: StatsCardProps) {
    return (
        <div className={`stats-card${wide ? " stats-card-wide" : ""}`}>
            <div className="stats-card-header">
                <h3 className="stats-card-title" dangerouslySetInnerHTML={{ __html: title }} />
                <button className="stats-zoom-reset" onClick={() => onResetZoom(id)}>
                    Reset Zoom
                </button>
            </div>
            <div
                className={`stats-chart-container${tall ? " stats-chart-tall" : ""}${chartClassName ? ` ${chartClassName}` : ""}`}
            >
                {children}
            </div>
            <p className="stats-zoom-hint">Scroll to zoom &middot; drag to pan</p>
        </div>
    );
}
