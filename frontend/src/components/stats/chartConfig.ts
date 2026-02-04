export const PALETTE = [
    "#d4a84b",
    "#9d7bc9",
    "#ff3333",
    "#3399ff",
    "#6b4c9a",
    "#f0d590",
    "#8b2942",
    "#a67c2e",
    "#3d2a5c",
    "#e8e0f0",
    "#c97bb4",
    "#7bc9a3",
];

export const zoomConfig = {
    zoom: {
        wheel: { enabled: true },
        pinch: { enabled: true },
        mode: "xy" as const,
    },
    pan: {
        enabled: true,
        mode: "xy" as const,
    },
};
