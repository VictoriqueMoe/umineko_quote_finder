import { useCallback, useRef, useState } from "react";

interface ShareButtonProps {
    audioId: string;
    lang?: string;
}

export function ShareButton({ audioId, lang }: ShareButtonProps) {
    const [copied, setCopied] = useState(false);
    const firstId = audioId.split(", ")[0];
    const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const handleClick = useCallback(() => {
        const effectiveLang = lang || "en";
        let url = window.location.origin + "/?quote=" + firstId;
        if (effectiveLang !== "en") {
            url += "&lang=" + effectiveLang;
        }
        navigator.clipboard.writeText(url).then(() => {
            setCopied(true);
            if (timeoutRef.current) {
                clearTimeout(timeoutRef.current);
            }
            timeoutRef.current = setTimeout(() => setCopied(false), 2000);
        });
    }, [firstId, lang]);

    return (
        <button className="share-btn" onClick={handleClick}>
            {copied ? "Link Copied" : "Share this Fragment"}
        </button>
    );
}
