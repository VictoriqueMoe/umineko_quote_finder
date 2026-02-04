interface EmptyStateProps {
    message?: string;
}

export function EmptyState({ message = "No quotes found in this fragment." }: EmptyStateProps) {
    return (
        <div className="empty-state">
            <div className="empty-icon">{"\uD83E\uDD8B"}</div>
            <h3 className="empty-title">The Golden Land remains silent</h3>
            <p className="empty-subtitle">{message}</p>
        </div>
    );
}
