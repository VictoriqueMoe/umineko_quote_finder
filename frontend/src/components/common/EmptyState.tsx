import { useAppContext } from "../../hooks/useAppContext";

interface EmptyStateProps {
    message?: string;
}

const GAME_EMPTY = {
    umineko: { icon: "\uD83E\uDD8B", title: "The Golden Land remains silent" },
    higurashi: { icon: "\uD83E\uDE77", title: "The cicadas have fallen silent" },
} as const;

export function EmptyState({ message = "No quotes found in this fragment." }: EmptyStateProps) {
    const { game } = useAppContext();
    const { icon, title } = GAME_EMPTY[game];

    return (
        <div className="empty-state">
            <div className="empty-icon">{icon}</div>
            <h3 className="empty-title">{title}</h3>
            <p className="empty-subtitle">{message}</p>
        </div>
    );
}
