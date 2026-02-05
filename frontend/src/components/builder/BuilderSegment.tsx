import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { audioUrl } from "../../api/client";
import type { BuilderSegment as SegmentType } from "../../hooks/useVoiceBuilder";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface BuilderSegmentProps {
    segment: SegmentType;
    index: number;
    audioPlayer: AudioPlayer;
    onRemove: (id: string) => void;
}

export function BuilderSegment({ segment, index, audioPlayer, onRemove }: BuilderSegmentProps) {
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
        id: segment.id,
    });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
    };

    const isActive = audioPlayer.state.activeId === `builder-${segment.id}`;
    const isPlaying = isActive && audioPlayer.state.isPlaying;

    const handlePlay = () => {
        const url = audioUrl(segment.charId, segment.audioId);
        audioPlayer.play(url, `builder-${segment.id}`);
    };

    return (
        <div ref={setNodeRef} style={style} className={`builder-segment${isDragging ? " dragging" : ""}`}>
            <button className="builder-segment-handle" {...attributes} {...listeners}>
                <span className="grip-icon">{"\u2807"}</span>
            </button>
            <span className="builder-segment-index">{index + 1}</span>
            <div className="builder-segment-info">
                <span className="builder-segment-character">{segment.characterName}</span>
                <span className="builder-segment-text">
                    {"\u201C"}
                    {segment.quoteText}
                    {"\u201D"}
                </span>
                {segment.episode && <span className="builder-segment-episode">Ep {segment.episode}</span>}
            </div>
            <div className="builder-segment-actions">
                <button
                    className={`builder-segment-play${isPlaying ? " playing" : ""}`}
                    onClick={handlePlay}
                    title="Preview clip"
                >
                    {isPlaying ? "\u275A\u275A" : "\u25B6"}
                </button>
                <button className="builder-segment-remove" onClick={() => onRemove(segment.id)} title="Remove clip">
                    {"\u2715"}
                </button>
            </div>
        </div>
    );
}
