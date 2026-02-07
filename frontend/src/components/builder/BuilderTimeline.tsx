import {
    closestCenter,
    DndContext,
    type DragEndEvent,
    PointerSensor,
    TouchSensor,
    useSensor,
    useSensors,
} from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { restrictToParentElement, restrictToVerticalAxis } from "@dnd-kit/modifiers";
import { BuilderSegment } from "./BuilderSegment";
import type { VoiceBuilder } from "../../hooks/useVoiceBuilder";
import type { AudioPlayer } from "../../hooks/useAudioPlayer";

interface BuilderTimelineProps {
    builder: VoiceBuilder;
    audioPlayer: AudioPlayer;
}

export function BuilderTimeline({ builder, audioPlayer }: BuilderTimelineProps) {
    const sensors = useSensors(
        useSensor(PointerSensor, {
            activationConstraint: { distance: 8 },
        }),
        useSensor(TouchSensor, {
            activationConstraint: { delay: 250, tolerance: 5 },
        }),
    );

    const handleDragEnd = (event: DragEndEvent) => {
        const { active, over } = event;
        if (over && active.id !== over.id) {
            builder.reorderSegments(String(active.id), String(over.id));
        }
    };

    return (
        <div className="builder-timeline">
            <div className="builder-timeline-header">
                <h3>Your Voice Build</h3>
                <span className="builder-counter">
                    {builder.segmentCount}/{builder.maxSegments} clips
                </span>
            </div>

            {builder.segments.length === 0 ? (
                <div className="builder-timeline-empty">
                    <p>Search above and add voice clips to begin crafting your dialogue.</p>
                </div>
            ) : (
                <DndContext
                    sensors={sensors}
                    collisionDetection={closestCenter}
                    onDragEnd={handleDragEnd}
                    modifiers={[restrictToVerticalAxis, restrictToParentElement]}
                >
                    <SortableContext items={builder.segments.map(s => s.id)} strategy={verticalListSortingStrategy}>
                        <div className="builder-segment-list">
                            {builder.segments.map((segment, index) => (
                                <BuilderSegment
                                    key={segment.id}
                                    segment={segment}
                                    index={index}
                                    audioPlayer={audioPlayer}
                                    onRemove={builder.removeSegment}
                                />
                            ))}
                        </div>
                    </SortableContext>
                </DndContext>
            )}
        </div>
    );
}
