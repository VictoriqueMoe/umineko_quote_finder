import { useState } from "react";
import { audioUrl, combinedAudioUrl, resolveCharId } from "../../api/client";
import { useAppContext } from "../../hooks/useAppContext";
import { AudioControls } from "./AudioControls";
import type { AudioPlayer as AudioPlayerType } from "../../hooks/useAudioPlayer";

interface AudioPlayerProps {
    audioId: string;
    characterId: string;
    audioCharMap?: Record<string, string>;
    audioPlayer: AudioPlayerType;
}

export function AudioPlayer({ audioId, characterId, audioCharMap, audioPlayer }: AudioPlayerProps) {
    const { game } = useAppContext();
    const [showIndividual, setShowIndividual] = useState(false);
    const [delay, setDelay] = useState(false);
    const ids = audioId.split(", ");
    const hasMultiple = ids.length > 1;

    const handleClipClick = (id: string) => {
        const charId = resolveCharId(id, characterId, audioCharMap);
        const url = audioUrl(game, charId, id);
        audioPlayer.play(url, id);
    };

    const handleCombinedClick = () => {
        const segments = ids.map(id => ({
            charId: resolveCharId(id, characterId, audioCharMap),
            audioId: id,
        }));
        const url = combinedAudioUrl(game, segments, delay);
        audioPlayer.play(url, `combined-${ids.join(",")}-${delay}`);
    };

    const isActive = (id: string) => audioPlayer.state.activeId === id;
    const isCombinedActive = audioPlayer.state.activeId === `combined-${ids.join(",")}-${delay}`;
    const isAnyActive = ids.some(id => isActive(id)) || isCombinedActive;

    return (
        <div className={`audio-player${isAnyActive && audioPlayer.state.isPlaying ? " playing" : ""}`}>
            {hasMultiple ? (
                <>
                    <div className="audio-clips">
                        <button
                            className={`audio-clip-btn audio-combined-btn${isCombinedActive ? " active" : ""}`}
                            onClick={handleCombinedClick}
                        >
                            {`\u25B6 Combined (${ids.length} clips)`}
                        </button>
                        <button
                            className={`audio-clip-btn audio-delay-btn${delay ? " active" : ""}`}
                            onClick={() => setDelay(!delay)}
                            title="Add a short pause between clips"
                        >
                            {delay ? "\u23F8 Delay On" : "\u23F5 Delay Off"}
                        </button>
                        <button className="audio-expand-btn" onClick={() => setShowIndividual(!showIndividual)}>
                            {showIndividual ? "\u25B4 Individual" : "\u25BE Individual"}
                        </button>
                    </div>
                    <div className={`audio-individual-clips${showIndividual ? " visible" : ""}`}>
                        {ids.map(id => (
                            <button
                                key={id}
                                className={`audio-clip-btn${isActive(id) ? " active" : ""}`}
                                onClick={() => handleClipClick(id)}
                            >
                                {`\u25B6 ${id}.ogg`}
                            </button>
                        ))}
                    </div>
                </>
            ) : (
                <div className="audio-clips">
                    {ids.map(id => (
                        <button
                            key={id}
                            className={`audio-clip-btn${isActive(id) ? " active" : ""}`}
                            onClick={() => handleClipClick(id)}
                        >
                            {`\u25B6 ${id}.ogg`}
                        </button>
                    ))}
                </div>
            )}
            <AudioControls audioPlayer={audioPlayer} isVisible={isAnyActive} />
        </div>
    );
}
