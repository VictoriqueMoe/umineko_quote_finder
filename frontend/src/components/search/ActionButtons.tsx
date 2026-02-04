interface ActionButtonsProps {
    onRandom: () => void;
    onClear: () => void;
}

export function ActionButtons({ onRandom, onClear }: ActionButtonsProps) {
    return (
        <div className="actions">
            <button className="action-btn" onClick={onRandom}>
                {"\u2726 Random Quote"}
            </button>
            <button className="action-btn" onClick={onClear}>
                Clear Results
            </button>
        </div>
    );
}
