import { useCallback, useEffect, useRef } from "react";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface AudioControlsProps {
    audioPlayer: AudioPlayer;
    isVisible: boolean;
}

function formatTime(sec: number): string {
    if (!sec || !isFinite(sec)) {
        return "0:00";
    }
    const m = Math.floor(sec / 60);
    const s = Math.floor(sec % 60);
    return m + ":" + (s < 10 ? "0" : "") + s;
}

export function AudioControls({ audioPlayer, isVisible }: AudioControlsProps) {
    const sliderRef = useRef<HTMLInputElement>(null);
    const savedVolume = localStorage.getItem("uminekoVolume") ?? "1.0";
    const volumeValue = Math.floor(parseFloat(savedVolume) * 100);

    const progressPct = audioPlayer.state.duration
        ? (audioPlayer.state.currentTime / audioPlayer.state.duration) * 100
        : 0;

    const handleTrackClick = useCallback(
        (e: React.MouseEvent<HTMLDivElement>) => {
            const rect = e.currentTarget.getBoundingClientRect();
            const ratio = (e.clientX - rect.left) / rect.width;
            audioPlayer.seek(ratio);
        },
        [audioPlayer],
    );

    const handleVolumeChange = useCallback(
        (e: React.ChangeEvent<HTMLInputElement>) => {
            const v = parseFloat(e.target.value) / 100;
            audioPlayer.setVolume(v);
            updateSliderFill(e.target);
        },
        [audioPlayer],
    );

    useEffect(() => {
        if (sliderRef.current) {
            updateSliderFill(sliderRef.current);
        }
    }, []);

    return (
        <div className={`audio-controls${isVisible ? " visible" : ""}`}>
            <div className="audio-track" onClick={handleTrackClick}>
                <div className="audio-progress" style={{ width: `${progressPct}%` }} />
            </div>
            <div className="audio-volume">
                <span>VOL</span>
                <input
                    ref={sliderRef}
                    className="audio-volume-slider"
                    type="range"
                    min="0"
                    max="100"
                    step="1"
                    defaultValue={volumeValue}
                    onChange={handleVolumeChange}
                />
            </div>
            <span className="audio-time">
                {formatTime(audioPlayer.state.currentTime)} / {formatTime(audioPlayer.state.duration)}
            </span>
        </div>
    );
}

function updateSliderFill(slider: HTMLInputElement) {
    const pct = slider.value;
    slider.style.background = `linear-gradient(to right, #d4a84b 0%, #d4a84b ${pct}%, #3d2a5c ${pct}%, #3d2a5c 100%)`;
}
