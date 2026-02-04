import { useCallback, useRef } from "react";
import type { StatsResponse } from "../../types/api";
import { StatsCard } from "./StatsCard";
import { TopSpeakersChart } from "./TopSpeakersChart";
import { LinesPerEpisodeChart } from "./LinesPerEpisodeChart";
import { TruthChart } from "./TruthChart";
import { InteractionsChart } from "./InteractionsChart";
import { PresenceChart } from "./PresenceChart";
import type { Chart } from "chart.js";

interface StatsViewProps {
    data: StatsResponse;
    episode: string;
}

export function StatsView({ data, episode }: StatsViewProps) {
    const chartsRef = useRef<Map<string, Chart>>(new Map());

    const ep = parseInt(episode) || 0;
    const epLabel =
        ep > 0
            ? data.episodeNames[ep]
                ? `Episode ${ep} \u2014 ${data.episodeNames[ep]}`
                : `Episode ${ep}`
            : "All Episodes";

    const hasAllEpisodes = data.linesPerEpisode && data.linesPerEpisode.length > 0;

    const registerChart = useCallback((id: string, chart: Chart) => {
        chartsRef.current.set(id, chart);
    }, []);

    const handleResetZoom = useCallback((chartId: string) => {
        const chart = chartsRef.current.get(chartId);
        if (chart) {
            (chart as Chart & { resetZoom: () => void }).resetZoom();
        }
    }, []);

    return (
        <>
            <div className="stats-header">
                <h2 className="stats-title">Umineko Statistics</h2>
                <p className="stats-subtitle">{epLabel} &mdash; English script lines</p>
            </div>
            <div className="stats-grid">
                <StatsCard id="chartTopSpeakers" title="Top Speakers" tall wide onResetZoom={handleResetZoom}>
                    <TopSpeakersChart data={data} onRegister={registerChart} />
                </StatsCard>
                {hasAllEpisodes && (
                    <StatsCard id="chartLinesPerEpisode" title="Lines per Episode" onResetZoom={handleResetZoom}>
                        <LinesPerEpisodeChart data={data} onRegister={registerChart} />
                    </StatsCard>
                )}
                {hasAllEpisodes && (
                    <StatsCard id="chartTruth" title="Red Truth &amp; Blue Truth" onResetZoom={handleResetZoom}>
                        <TruthChart data={data} onRegister={registerChart} />
                    </StatsCard>
                )}
                <StatsCard
                    id="chartInteractions"
                    title="Character Interactions"
                    tall
                    wide
                    onResetZoom={handleResetZoom}
                >
                    <InteractionsChart data={data} onRegister={registerChart} />
                </StatsCard>
                {hasAllEpisodes && (
                    <StatsCard
                        id="chartPresence"
                        title="Character Presence by Episode"
                        wide
                        onResetZoom={handleResetZoom}
                    >
                        <PresenceChart data={data} onRegister={registerChart} />
                    </StatsCard>
                )}
            </div>
        </>
    );
}
