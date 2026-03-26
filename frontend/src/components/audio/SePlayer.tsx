import { seAudioUrl } from "../../api/client";
import type { SoundEffect } from "../../types/api";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface SePlayerProps {
    soundEffects: SoundEffect[];
    audioPlayer: AudioPlayer;
}

export function SePlayer({ soundEffects, audioPlayer }: SePlayerProps) {
    return (
        <div className="se-player">
            <span className="se-label">Sound Effects</span>
            <div className="se-clips">
                {soundEffects.map(se => {
                    const id = `se-${se.filename}`;
                    const isActive = audioPlayer.state.activeId === id;
                    return (
                        <button
                            key={`${se.filename}-${se.afterClip}`}
                            className={`se-clip-btn${isActive ? " active" : ""}`}
                            onClick={() => audioPlayer.play(seAudioUrl(se.filename), id)}
                        >
                            {`\u266A ${se.filename}.ogg`}
                        </button>
                    );
                })}
            </div>
        </div>
    );
}
