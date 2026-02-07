import { useCallback, useRef, useState } from "react";
import { AudioControls } from "../audio/AudioControls";
import type { VoiceBuilder } from "../../hooks/useVoiceBuilder";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface BuilderControlsProps {
    builder: VoiceBuilder;
    audioPlayer: AudioPlayer;
}

export function BuilderControls({ builder, audioPlayer }: BuilderControlsProps) {
    const [downloading, setDownloading] = useState(false);
    const [copied, setCopied] = useState(false);
    const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const isEmpty = builder.segmentCount === 0;
    const isCombinedActive = audioPlayer.state.activeId === "builder-combined";
    const isCombinedPlaying = isCombinedActive && audioPlayer.state.isPlaying;

    const handlePlayCombined = useCallback(() => {
        if (!builder.combinedUrl) {
            return;
        }
        audioPlayer.play(builder.combinedUrl, "builder-combined");
    }, [builder.combinedUrl, audioPlayer]);

    const handleDownload = useCallback(async () => {
        if (!builder.combinedUrl) {
            return;
        }
        setDownloading(true);
        try {
            const response = await fetch(builder.combinedUrl);
            if (!response.ok) {
                throw new Error(`Download failed: ${response.status}`);
            }
            const blob = await response.blob();
            const blobUrl = URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = blobUrl;
            a.download = "voice-build.ogg";
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(blobUrl);
        } catch (err) {
            console.error("Download failed:", err);
        } finally {
            setDownloading(false);
        }
    }, [builder.combinedUrl]);

    const handleShare = useCallback(() => {
        if (!builder.shareUrl) {
            return;
        }
        navigator.clipboard.writeText(builder.shareUrl).then(() => {
            setCopied(true);
            if (copyTimeoutRef.current) {
                clearTimeout(copyTimeoutRef.current);
            }
            copyTimeoutRef.current = setTimeout(() => setCopied(false), 2000);
        });
    }, [builder.shareUrl]);

    const handleClear = useCallback(() => {
        audioPlayer.stop();
        builder.clearAll();
    }, [audioPlayer, builder]);

    return (
        <div className="builder-controls">
            <AudioControls audioPlayer={audioPlayer} isVisible={isCombinedActive} />
            <div className="builder-controls-buttons">
                <button
                    className="builder-control-btn builder-play-combined"
                    disabled={isEmpty}
                    onClick={handlePlayCombined}
                >
                    {isCombinedPlaying ? "\u275A\u275A Pause" : "\u25B6 Play Combined"}
                </button>
                <button
                    className="builder-control-btn builder-download"
                    disabled={isEmpty || downloading}
                    onClick={handleDownload}
                >
                    {downloading ? "Downloading..." : "\u2913 Download Audio"}
                </button>
                <button className="builder-control-btn builder-share" disabled={isEmpty} onClick={handleShare}>
                    {copied ? "Link Copied!" : "\u2197 Share Link"}
                </button>
                <button className="builder-control-btn builder-clear" disabled={isEmpty} onClick={handleClear}>
                    {"\u2715 Clear All"}
                </button>
            </div>
        </div>
    );
}
