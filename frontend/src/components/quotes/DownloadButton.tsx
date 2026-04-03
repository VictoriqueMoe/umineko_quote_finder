import { useCallback, useState } from "react";
import { ogImageUrl } from "../../api/client";
import { useAppContext } from "../../hooks/useAppContext";

interface DownloadButtonProps {
    audioId: string;
    lang?: string;
}

export function DownloadButton({ audioId, lang }: DownloadButtonProps) {
    const { game } = useAppContext();
    const [downloading, setDownloading] = useState(false);
    const firstId = audioId.split(", ")[0];

    const handleClick = useCallback(async () => {
        setDownloading(true);
        try {
            const url = ogImageUrl(game, firstId, lang ?? "en", true);
            const response = await fetch(url, { cache: "no-cache" });
            if (!response.ok) {
                throw new Error(`Failed to fetch image: ${response.status}`);
            }
            const blob = await response.blob();
            const blobUrl = URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = blobUrl;
            a.download = `${firstId}.png`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(blobUrl);
        } catch (err) {
            console.error("Download failed:", err);
        } finally {
            setDownloading(false);
        }
    }, [game, firstId, lang]);

    return (
        <button className="share-btn" onClick={handleClick} disabled={downloading}>
            {downloading ? "Downloading..." : "Save as Image"}
        </button>
    );
}
