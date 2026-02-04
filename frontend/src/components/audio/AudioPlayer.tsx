import { useState } from "react";
import { audioUrl, combinedAudioUrl } from "../../api/client";
import { AudioControls } from "./AudioControls";
import type { AudioPlayer as AudioPlayerType } from "../../hooks/useAudioPlayer";

interface AudioPlayerProps {
    audioId: string;
    characterId: string;
    audioPlayer: AudioPlayerType;
}

export function AudioPlayer({ audioId, characterId, audioPlayer }: AudioPlayerProps) {
    const [showIndividual, setShowIndividual] = useState(false);
    const ids = audioId.split(", ");
    const hasMultiple = ids.length > 1;

    const handleClipClick = (id: string) => {
        const url = audioUrl(characterId, id);
        audioPlayer.play(url, id);
    };

    const handleCombinedClick = () => {
        const allIds = ids.join(",");
        const url = combinedAudioUrl(characterId, allIds);
        audioPlayer.play(url, `combined-${allIds}`);
    };

    const isActive = (id: string) => audioPlayer.state.activeId === id;
    const isCombinedActive = audioPlayer.state.activeId === `combined-${ids.join(",")}`;
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
