import { useEffect, useRef } from "react";
import { useAppContext } from "../../hooks/useAppContext";
import { useAudioPlayer } from "../../hooks/useAudioPlayer";
import { useVoiceBuilder, type VoiceBuilder } from "../../hooks/useVoiceBuilder";
import { BuilderSearch } from "./BuilderSearch";
import { BuilderTimeline } from "./BuilderTimeline";
import { BuilderControls } from "./BuilderControls";

interface VoiceBuilderViewProps {
    onClose: () => void;
    initialBuilder?: string | null;
}

export function VoiceBuilderView({ onClose, initialBuilder }: VoiceBuilderViewProps) {
    const { language } = useAppContext();
    const audioPlayer = useAudioPlayer();
    const builder = useVoiceBuilder();
    const hydratedRef = useRef(false);

    useEffect(() => {
        if (initialBuilder && initialBuilder !== "1" && !hydratedRef.current) {
            hydratedRef.current = true;
            builder.loadFromUrl(initialBuilder, language);
        }
    }, [initialBuilder, language, builder]);

    return (
        <div className="builder-view">
            <div className="builder-header">
                <div className="builder-header-text">
                    <h2 className="builder-title">{"\u266B"} Voice Builder</h2>
                    <p className="builder-subtitle">Craft a custom voice dialogue from fragments</p>
                </div>
                <button className="builder-close-btn" onClick={onClose} title="Close builder">
                    {"\u2715"}
                </button>
            </div>

            <BuilderSearch builder={builder} audioPlayer={audioPlayer} />
            <BuilderTimeline builder={builder} audioPlayer={audioPlayer} />
            <BuilderControls builder={builder} audioPlayer={audioPlayer} />
        </div>
    );
}

export type { VoiceBuilder };
