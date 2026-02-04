import { useEffect, useRef } from "react";

const BUTTERFLY_COUNT = 8;

export function Butterflies() {
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const container = containerRef.current;
        if (!container) {
            return;
        }

        for (let i = 0; i < BUTTERFLY_COUNT; i++) {
            const butterfly = document.createElement("div");
            butterfly.className = "butterfly";
            butterfly.textContent = "\uD83E\uDD8B";
            butterfly.style.setProperty("--start-x", `${Math.random() * 100}vw`);
            butterfly.style.setProperty("--duration", `${20 + Math.random() * 15}s`);
            butterfly.style.setProperty("--delay", `${Math.random() * 20}s`);
            butterfly.style.left = `${Math.random() * 100}%`;
            container.appendChild(butterfly);
        }

        return () => {
            container.innerHTML = "";
        };
    }, []);

    return <div className="butterflies" ref={containerRef} />;
}
